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

// TestNewGTLogNoPanicOnUnwritableDir 日志目录建不出来:stderr 报因、句柄置 nil、
// 日志走仍可用的通道,进程不 fatal 不崩溃。
//
// 回归背景:CI(linux runner)上 productLogDir 落到无权限目录 → newLogSaveHandler
// MkdirAll 失败返回 nil → 修复前 nil 句柄被包进 MultiWriter,后台维护 goroutine
// 首次写日志即 nil RotateLogs.Write → SIGSEGV 全进程崩。
//
// 覆盖边界如实声明:本防护只覆盖"logDir 建不出来"这一种创建期失败;
// 目录存在但不可写、运行中权限变化、盘满等写入期失败不在本库探测范围。
func TestNewGTLogNoPanicOnUnwritableDir(t *testing.T) {
	skipIfPermUnenforceable(t)

	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	withLogConfig(t, filepath.Join(parent, "denied"), true, true)

	aLog := NewGTLog("nopanic_test")
	if aLog == nil {
		t.Fatal("NewGTLog should not return nil")
	}
	if aLog.rotateHandle != nil {
		t.Fatal("目录不可写时 rotateHandle 应为 nil")
	}
	if !aLog.saveFileEnabled {
		t.Fatal("创建失败不应改写配置语义(saveFileEnabled 保持 true)")
	}
	// 写日志不 panic 即为防护生效(修复前此处 SIGSEGV)
	aLog.LogInfof("still alive [%s]", "ok")
}

// TestRolloverRebuildsHandle 日切触发句柄重建:换日期目录、按三态重新接线、关闭旧句柄。
// 回归背景:修复前日切直接 SetOutput(rLog),丢 keepStdout 双输出与 stripANSI,旧句柄泄漏 fd。
func TestRolloverRebuildsHandle(t *testing.T) {
	// 句柄随进程存活(守护进程设计,无 Close API),windows 下 t.TempDir 清理会因文件占用失败
	logDir, err := os.MkdirTemp("", "gtbox_log_rollover")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(logDir) })
	withLogConfig(t, logDir, true, false)

	aLog := NewGTLog("rollover_test")
	oldHandle := aLog.rotateHandle
	if oldHandle == nil {
		t.Fatal("可写目录下启动应建出句柄")
	}

	// 伪造"昨天",促使日切判定生效
	fakeYesterday := aLog.logDir + "/1970-01-01"
	aLog.Lock()
	aLog.logDirWithDate = fakeYesterday
	aLog.rotateHandleDir = fakeYesterday
	aLog.Unlock()

	aLog.checkAndUpdateLogDir()

	if aLog.rotateHandle == nil || aLog.rotateHandle == oldHandle {
		t.Fatal("日切后应重建新句柄")
	}
	if aLog.rotateHandleDir == fakeYesterday {
		t.Fatal("rotateHandleDir 应更新为新日期目录")
	}
	aLog.LogInfof("after rollover [%s]", "ok")
}
