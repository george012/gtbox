package gtbox

import (
	"fmt"
	"github.com/george012/gtbox/config"
	"github.com/george012/gtbox/gtbox_coding"
	"github.com/george012/gtbox/gtbox_http"
	"github.com/george012/gtbox/gtbox_log"
)

type RunMode int

const (
	RunModeUnknown RunMode = iota
	RunModeDebug
	RunModeRelease
	RunModeTest
)

func (rm RunMode) String() string {
	switch rm {
	case RunModeDebug:
		return "Debug"
	case RunModeRelease:
		return "Release"
	case RunModeTest:
		return "Test"
	case RunModeUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

type GTAppSignalInfo struct {
	SigCode string
	Msg     string
}

var (
	currentRunMode RunMode
)

func GetCurrentRunMode() RunMode {
	return currentRunMode
}

// SetupGTBox
// 必须--YES
// 必须使用此方法初始化工具库,未使用此方法初始化，无法使用完整功能，亦存在兼容性问题
// logMaxSaveDays Log是否开启文件存储模式
// log_dir 自定义日志目录,默认为:/usr/logs/${projectName},如果传"" 即使用默认值
// httpRequestTimeOut 网络请求超时时间
// projectName--项目名称，
// run_mode 运行模式 debug
// logLevel--日志等级，
// logMaxSaveTime--默认365天,
// logSaveType--日志分片格式，默认按天分片，可选按小时分片
func SetupGTBox(projectName string, run_mode RunMode, productLogDir string, logMaxSaveDays int64, logSaveType gtbox_log.GTLogSaveType, httpRequestTimeOut int) {
	enableSaveLogFile := false
	logLevel := gtbox_log.GTLogStyleDebug
	currentRunMode = run_mode
	switch run_mode {
	case RunModeDebug:
		enableSaveLogFile = false
		logLevel = gtbox_log.GTLogStyleDebug
	case RunModeTest:
		enableSaveLogFile = true
		logLevel = gtbox_log.GTLogStyleDebug
	case RunModeRelease:
		enableSaveLogFile = true
		logLevel = gtbox_log.GTLogStyleInfo
	}

	gtbox_log.SetupLogTools(projectName, enableSaveLogFile, logLevel, logMaxSaveDays, logSaveType, productLogDir)

	gtbox_http.DefaultTimeout = httpRequestTimeOut
	config.IsSetup = true
	fmt.Printf("gtbox Tools Setup End\nProjcetName=[%s]\nrunMode=[%s]\nlogLeve=[%s]\nproduct main logdir=[%s]\nlogCutType=[%s]\nlogSaveDays=[%d]\nhttpRequestTimeout=[%d Second]\ngtbox Effective lines of code=[%d]\n",
		gtbox_log.GetProjectName(),
		run_mode.String(),
		gtbox_log.GetLogLevel().String(),
		gtbox_log.GetProductMainLogDir(),
		logSaveType.String(),
		logMaxSaveDays,
		gtbox_http.DefaultTimeout,
		gtbox_coding.GetProjectCodeLines(),
	)
}
