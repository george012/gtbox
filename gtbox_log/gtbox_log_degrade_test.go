package gtbox_log

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestNewGTLogDegradeOnUnwritableDir 日志目录不可写时必须降级为只 stdout,绝不 panic。
//
// 回归背景:CI(linux runner)上 productLogDir 落到无权限目录 → newLogSaveHandler
// MkdirAll 失败返回 nil → 修复前 nil 句柄被包进 MultiWriter,后台维护 goroutine
// 首次写日志即 nil RotateLogs.Write → SIGSEGV 全进程崩。
func TestNewGTLogDegradeOnUnwritableDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX 目录权限在 windows 不生效,无法构造不可写目录")
	}
	if os.Geteuid() == 0 {
		t.Skip("root 无视目录权限,无法构造不可写目录")
	}

	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })

	cfg := instanceConfig()
	oldDir, oldSave, oldKeep := cfg.productLogDir, cfg.enableSaveLogFile, cfg.keepStdout
	cfg.productLogDir = filepath.Join(parent, "denied")
	cfg.enableSaveLogFile = true
	cfg.keepStdout = true
	t.Cleanup(func() {
		cfg.productLogDir, cfg.enableSaveLogFile, cfg.keepStdout = oldDir, oldSave, oldKeep
	})

	aLog := NewGTLog("degrade_test")
	if aLog == nil {
		t.Fatal("NewGTLog should not return nil")
	}
	if aLog.saveFileEnabled {
		t.Fatal("unwritable log dir should degrade to stdout-only (saveFileEnabled=false)")
	}
	// 写日志不 panic 即为降级成功(修复前此处 SIGSEGV)
	aLog.LogInfof("degrade path alive [%s]", "ok")
}
