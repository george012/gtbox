/*
Package gtbox 工具库主入口
*/
package gtbox

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_http"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type GTAppSignalInfo struct {
	SigCode string
	Msg     string
}

// GTSysUseSignalWaitAppExit 处理程序信号,并且做一些操作，比如：保存状态、保存配置文件
func GTSysUseSignalWaitAppExit(exitHandleFunc func(sigInfo *GTAppSignalInfo)) {
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	aSignal := &GTAppSignalInfo{
		SigCode: fmt.Sprintf("%s", <-chSig),
		Msg:     "程序即将退出",
	}
	exitHandleFunc(aSignal)
}

// SetupGTBox
// 必须--YES
// 必须使用此方法初始化工具库,未使用此方法初始化，无法使用完整功能，亦存在兼容性问题
// debugToCut App如果是Debug模式默认不开启日志切片，方便IDE调试
// httpRequestTimeOut 网络请求超时时间
// projectName--项目名称，
// logLevel--日志等级，
// logMaxSaveTime--默认365天,
// logSaveType--日志分片格式，默认按天分片，可选按小时分片
func SetupGTBox(projectName string, debugMode bool, logLevel logrus.Level, logMaxSaveDays int64, logSaveType gtbox_log.GTLogSaveType, httpRequestTimeOut int) {
	debugToCut := false
	if debugMode == false {
		debugToCut = true
	}
	gtbox_log.SetupLogTools(projectName, debugToCut, logLevel, logMaxSaveDays, logSaveType)
	gtbox_http.DefaultTimeout = httpRequestTimeOut
	fmt.Printf("gtbox Tools Setup End\nProjcetName=[%s]\nlogLeve=[%s]\nlogpath=[%s]\nlogCutType=[%s]\nlogSaveDays=[%d]\nhttpRequestTimeout=[%d Second]\n",
		gtbox_log.ProjectName,
		gtbox_log.LogLevel.String(),
		gtbox_log.LogPath,
		logSaveType.String(),
		logMaxSaveDays,
		gtbox_http.DefaultTimeout,
	)
}
