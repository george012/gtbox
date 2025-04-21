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

// checkAndUpdateLogDir 检查并更新日志目录
func (aLog *GTLog) checkAndUpdateLogDir() {
	// TODO 每分钟检查一次是否需要更新日志文件路径
	now := time.Now().UTC()
	newLogDirWithDate := fmt.Sprintf("%s/%s", aLog.logDir, now.Format("2006-01-02"))

	if aLog.logDirWithDate != newLogDirWithDate {
		aLog.logDirWithDate = newLogDirWithDate
		rLog := newLogSaveHandler(aLog)
		aLog.logger.SetOutput(rLog)
		aLog.lastCheckTime = now
	}
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

	now := time.Now()
	maxAge := time.Duration(instanceConfig().logMaxSaveDays) * 24 * time.Hour

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirPath := fmt.Sprintf("%s/%s", aLog.logDir, dir.Name())

		// 跳过当前正在使用的日志目录
		if dirPath == aLog.logDirWithDate {
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
