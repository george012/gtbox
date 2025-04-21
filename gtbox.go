package gtbox

import (
	"fmt"
	"github.com/george012/gtbox/config"
	"github.com/george012/gtbox/gtbox_coding"
	"github.com/george012/gtbox/gtbox_encryption"
	"github.com/george012/gtbox/gtbox_http"
	"github.com/george012/gtbox/gtbox_log"
	"time"
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

// SetupGTBox [☑]Required
/*
	en: Setup gtbox;
	zh-CN: 初始化 gtbox 必须使用此方法初始化工具库,未使用此方法初始化，无法使用完整功能，亦存在兼容性问题;
	@params [☑]projectName en:  ;zh-CN: 项目名称;
	@params [☑]run_mode en:  ;zh-CN: 运行模式 debug|test|release;
	@params [☑]logMaxSaveDays en:  ;zh-CN: 日志存储最大天数;
	@params [☐]productLogDir en:  ;zh-CN: 自定义日志目录,默认为:/usr/logs/${projectName},如果传"" 即使用默认值;
	@params [☑]logSaveType en:  ;zh-CN: 日志存储类型：按天切片|按小时切片 GTLogSaveTypeDays | GTLogSaveHours;
	@params [☑]httpRequestTimeOut en:  ;zh-CN: 网络请求超时时间;
*/
func SetupGTBox(projectName string, runMode RunMode, productLogDir string, logMaxSaveDays int64, logSaveType gtbox_log.GTLogSaveType, httpRequestTimeOut time.Duration) {
	enableSaveLogFile := false
	logLevel := gtbox_log.GTLogStyleDebug
	currentRunMode = runMode
	switch runMode {
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
	fmt.Printf("gtbox Tools Setup End\nProjcetName=[%s]\nrunMode=[%s]\nlogLeve=[%s]\nproduct main logdir=[%s]\nlogCutType=[%s]\nlogSaveDays=[%d]\nhttpRequestTimeout=[%.2f Second]\ngtbox Effective lines of code=[%d]\nencryption_version=[%s]\n",
		gtbox_log.GetProjectName(),
		runMode.String(),
		gtbox_log.GetLogLevel(),
		gtbox_log.GetProductMainLogDir(),
		logSaveType.String(),
		logMaxSaveDays,
		gtbox_http.DefaultTimeout.Seconds(),
		gtbox_coding.GetProjectCodeLines(),
		gtbox_encryption.GetEncryptionLibVersion(),
	)
}
