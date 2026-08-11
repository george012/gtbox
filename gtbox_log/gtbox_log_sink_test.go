package gtbox_log

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestFileOnlyUnwritableDirRefusesToStart 只 file 模式(Release/Test:file 是唯一输出通道)
// 日志目录不可写 → 显式拒绝启动(exit 1),绝不盲跑。os.Exit 需子进程验证(Go 标准崩溃测试模式)。
func TestFileOnlyUnwritableDirRefusesToStart(t *testing.T) {
	skipIfPermUnenforceable(t)

	if dir := os.Getenv("GTBOX_LOG_FATAL_TEST_DIR"); dir != "" {
		// 子进程分支:构造只 file 模式 + 不可写目录,NewGTLog 应 os.Exit(1) 不返回
		cfg := instanceConfig()
		cfg.productLogDir = dir
		cfg.enableSaveLogFile = true
		cfg.keepStdout = false
		NewGTLog("fatal_test")
		os.Exit(0) // 不应到达
	}

	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })

	cmd := exec.Command(os.Args[0], "-test.run=TestFileOnlyUnwritableDirRefusesToStart")
	cmd.Env = append(os.Environ(), "GTBOX_LOG_FATAL_TEST_DIR="+filepath.Join(parent, "denied"))
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("file-only 模式日志目录不可写应 exit 1 拒绝启动, got %v", err)
	}
}

// TestLogFileDirAndNameSameClock 文件名日期与日期目录必须同钟(UTC)。
// 回归背景:rotatelogs 默认 Local 时钟,目录命名用 UTC——UTC+8 机器每天本地 00:00-08:00,
// 新一天(本地)的文件落在旧 UTC 日期目录,即"第二天的日志在上一天的目录里"。
func TestLogFileDirAndNameSameClock(t *testing.T) {
	// 句柄随进程存活(守护进程设计,无 Close API),windows 下 t.TempDir 清理会因文件占用失败
	logDir, err := os.MkdirTemp("", "gtbox_log_clock")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(logDir) })
	withLogConfig(t, logDir, true, false)

	aLog := NewGTLog("clock_test")
	aLog.LogInfof("clock probe [%s]", "ok")
	time.Sleep(50 * time.Millisecond) // 等 maintenance 首轮写完

	matches, err := filepath.Glob(filepath.Join(logDir, "clock_test", "*", "run.*.log"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("expect log file created, matches=%v err=%v", matches, err)
	}
	for _, m := range matches {
		dirDate := filepath.Base(filepath.Dir(m))          // 目录名 = YYYY-MM-DD
		fileName := filepath.Base(m)                       // run.YYYY-MM-DD_HH.log
		fileDate := fileName[len("run.") : len("run.")+10] // 文件名里的日期段
		if dirDate != fileDate {
			t.Fatalf("目录日期 %s 与文件名日期 %s 不同钟(时区裂脑回归)", dirDate, fileDate)
		}
	}
}
