package gtbox_log

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// skipIfPermUnenforceable 无法构造"不可写目录"的环境跳过(windows 不吃 POSIX 权限,root 无视权限)
func skipIfPermUnenforceable(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX 目录权限在 windows 不生效,无法构造不可写目录")
	}
	if os.Geteuid() == 0 {
		t.Skip("root 无视目录权限,无法构造不可写目录")
	}
}

// withLogConfig 临时改写包级 config 单例,测试结束还原
func withLogConfig(t *testing.T, logDir string, enableSave bool, keepStdout bool) {
	t.Helper()
	cfg := instanceConfig()
	oldDir, oldSave, oldKeep := cfg.productLogDir, cfg.enableSaveLogFile, cfg.keepStdout
	cfg.productLogDir = logDir
	cfg.enableSaveLogFile = enableSave
	cfg.keepStdout = keepStdout
	t.Cleanup(func() {
		cfg.productLogDir, cfg.enableSaveLogFile, cfg.keepStdout = oldDir, oldSave, oldKeep
	})
}

// TestNewGTLogDegradeAndSelfHeal 双路模式日志目录不可写:降级只 stdout 不 panic,
// 目录恢复可写后维护逻辑自愈重建 file sink。
//
// 回归背景:CI(linux runner)上 productLogDir 落到无权限目录 → newLogSaveHandler
// MkdirAll 失败返回 nil → 修复前 nil 句柄被包进 MultiWriter,后台维护 goroutine
// 首次写日志即 nil RotateLogs.Write → SIGSEGV 全进程崩。
func TestNewGTLogDegradeAndSelfHeal(t *testing.T) {
	skipIfPermUnenforceable(t)

	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	withLogConfig(t, filepath.Join(parent, "denied"), true, true)

	aLog := NewGTLog("degrade_test")
	if aLog == nil {
		t.Fatal("NewGTLog should not return nil")
	}
	if !aLog.saveFileEnabled {
		t.Fatal("降级不应永久关落盘:saveFileEnabled 保持 true 供维护 tick 重试自愈")
	}
	if aLog.rotateHandle != nil {
		t.Fatal("目录不可写时 rotateHandle 应为 nil")
	}
	// 写日志不 panic 即为降级成功(修复前此处 SIGSEGV)
	aLog.LogInfof("degrade path alive [%s]", "ok")

	// 自愈:恢复目录可写,维护逻辑(生产中每分钟 tick)重建句柄
	if err := os.Chmod(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	aLog.checkAndUpdateLogDir()
	if aLog.rotateHandle == nil {
		t.Fatal("目录恢复可写后应自愈重建 file sink")
	}
	if aLog.sinkFailedDir != "" {
		t.Fatal("自愈后失败标记应清空")
	}
	aLog.LogInfof("recovered [%s]", "ok")
}
