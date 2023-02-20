package gtbox

import (
	"fmt"
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
// debugToCut debug模式下是否开启日志分割,默认false,方便IDE调试
// projectName--项目名称，
// logLevel--日志等级，
// logMaxSaveTime--默认365天,
// logSaveType--日志分片格式，默认按天分片，可选按小时分片
func SetupGTBox(projectName string, debugToCut bool, logLevel logrus.Level, logMaxSaveDays int64, logSaveType gtbox_log.GTLogSaveType) {
	gtbox_log.SetupLogTools(projectName, debugToCut, logLevel, logMaxSaveDays, logSaveType)
	fmt.Printf("gtbox Tools Setup End\nProjcetName=[%s], logLeve=[%s], logpath=[%s] logCutType=[%s] logSaveDays=[%d]\n", gtbox_log.ProjectName, gtbox_log.LogLevel.String(), gtbox_log.LogPath, logSaveType.String(), logMaxSaveDays)
}
