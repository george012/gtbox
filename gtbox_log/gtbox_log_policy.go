package gtbox_log

import (
	"fmt"
	"os"
	"time"
)

func (aLog *GTLog) startLogMaintenance(firstRunFunc func(done chan struct{})) {
	done := make(chan struct{}) // 创建无缓冲通道用于同步

	go func() {
		firstRunDone := false                 // 用于跟踪首次循环是否完成
		ticker := time.NewTicker(time.Minute) // 每分钟检查一次
		defer ticker.Stop()

		for {
			if !firstRunDone { // 如果还没有执行第一次
				// 首次执行传入的检查和清理函数
				aLog.checkAndUpdateLogDir()
				aLog.cleanOldLogs()
				firstRunDone = true // 标记第一次运行已完成
				firstRunFunc(done)  // 调用 firstRunFunc，并在其内部关闭通道
			} else {
				select {
				case <-ticker.C:
					// 定时任务：每分钟检查一次
					aLog.checkAndUpdateLogDir()
					aLog.cleanOldLogs()
				}
			}
		}
	}()

	<-done // 等待通道关闭，表示首次执行已完成
}

// checkAndUpdateLogDir 日切维护:UTC 日期变了就重建 rotate 句柄并切换日期目录
// (触发条件与原设计一致:目录日期变化,每日一次)。重建失败 stderr 报因,
// 旧句柄(若在)继续写,不中断运行中的服务,下次日切再试。
// 修复前此处直接 SetOutput(rLog):日切后丢 keepStdout 双输出与 stripANSI,
// rLog 为 nil 时(目录不可写)后续写日志直接 SIGSEGV,且旧句柄从不关闭(每日泄漏一个 fd)。
func (aLog *GTLog) checkAndUpdateLogDir() {
	aLog.Lock()
	defer aLog.Unlock()

	if !aLog.saveFileEnabled {
		return
	}
	now := time.Now().UTC()
	wantDir := fmt.Sprintf("%s/%s", aLog.logDir, now.Format("2006-01-02"))
	if aLog.logDirWithDate == wantDir {
		return
	}

	aLog.logDirWithDate = wantDir
	rLog := newLogSaveHandler(aLog)
	if rLog == nil {
		fmt.Fprintf(os.Stderr, "[gtbox_log] log file sink rebuild failed, keep last sink [dir=%s]\n", wantDir)
		return
	}

	oldHandle := aLog.rotateHandle
	aLog.rotateHandle = rLog
	aLog.rotateHandleDir = wantDir
	aLog.wireOutput()
	if oldHandle != nil {
		_ = oldHandle.Close()
	}
	aLog.lastCheckTime = now
}

func (aLog *GTLog) cleanOldLogs() {
	if aLog.saveFileEnabled == false {
		return
	}

	dirs, err := os.ReadDir(aLog.logDir)
	if err != nil {
		aLog.logF(GTLogStyleError, "Error reading log directory: %s\n", err)
		return
	}

	// UTC 与目录命名同钟(修复:原 time.Now() 本地时间对比 UTC-naive 目录日期,清理期限偏差 ±时区)
	now := time.Now().UTC()
	maxAge := time.Duration(instanceConfig().logMaxSaveDays) * 24 * time.Hour

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirPath := fmt.Sprintf("%s/%s", aLog.logDir, dir.Name())

		// 跳过当前目标目录与在用句柄的归属目录(日切重建失败时二者可能不同)
		if dirPath == aLog.logDirWithDate || dirPath == aLog.rotateHandleDir {
			continue
		}

		// 假设目录名格式为 YYYY-MM-DD
		dirDate, err := time.Parse("2006-01-02", dir.Name())
		if err != nil {
			// 如果解析失败，跳过此目录
			continue
		}

		// 判断目录是否超出保存期限
		if now.Sub(dirDate) > maxAge {
			err := os.RemoveAll(dirPath)
			if err != nil {
				aLog.logF(GTLogStyleError, "Error removing directory: %s, error: %v\n", dirPath, err)
			} else {
				aLog.logF(GTLogStyleInfo, "Deleted old log directory: %s\n", dirPath)
			}
		}
	}
}
