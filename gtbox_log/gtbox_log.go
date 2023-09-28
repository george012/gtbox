/*
Package gtbox_log Log工具
*/
package gtbox_log

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_color"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// GTLogStyle 日志样式
type GTLogStyle int

const (
	GTLogStyleDebug   GTLogStyle = iota // Debug
	GTLogStyleError                     // Error
	GTLogStyleWarning                   // Warning
	GTLogStyleInfo                      // Info
	GTLogStyleTrace                     // Trace
	GTLogStyleFatal                     // Fatal
)

func (aStyle GTLogStyle) String() string {
	switch aStyle {
	case GTLogStyleFatal:
		return "fatal"
	case GTLogStyleTrace:
		return "trace"
	case GTLogStyleInfo:
		return "info"
	case GTLogStyleWarning:
		return "warning"
	case GTLogStyleError:
		return "error"
	case GTLogStyleDebug:
		return "debug"
	default:
		return "debug"
	}
}

// GTLogSaveType 日志分片类型
type GTLogSaveType int

const (
	GTLogSaveTypeDays GTLogSaveType = iota //按日分片
	GTLogSaveHours                         //按小时分片
)

func (aFlag GTLogSaveType) String() string {
	switch aFlag {
	case GTLogSaveTypeDays:
		return "Days"
	case GTLogSaveHours:
		return "Hours"
	default:
		return "Unknown"
	}
}

type GTLog struct {
	mux sync.RWMutex
}

var (
	ALog              *GTLog
	GTLogOnce         sync.Once
	ProjectName       = "test"
	LogLevel          = logrus.DebugLevel
	LogSaveMaxDays    int64
	LogSaveFlag       = GTLogSaveTypeDays
	LogPath           = "./logs/run"
	EnableSaveLogFile = false // EnableSaveLogFile	开启日志文件存储
)

// GTGetLogsDir 获取Log目录
func GTGetLogsDir() string {
	return LogPath
}

func (alog *GTLog) Setup() {

	//	设置Log
	if EnableSaveLogFile == true {

		if runtime.GOOS == "linux" {
			LogPath = "/var/log/" + strings.ToLower(ProjectName) + "/run" + "_" + ProjectName
		} else {
			LogPath = "./logs/run" + "_" + ProjectName
		}
		/* 日志轮转相关函数
		   `WithLinkName` 为最新的日志建立软连接
		   `WithRotationTime` 设置日志分割的时间，隔多久分割一次
		   WithMaxAge 和 WithRotationCount二者只能设置一个
		    `WithMaxAge` 设置文件清理前的最长保存时间
		    `WithRotationCount` 设置文件清理前最多保存的个数
		*/
		// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
		logRotaionFlag := time.Hour * 24

		if LogSaveFlag == GTLogSaveHours {
			logRotaionFlag = time.Hour
		}
		writer, _ := rotatelogs.New(
			LogPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(LogPath),
			rotatelogs.WithMaxAge(time.Duration(LogSaveMaxDays)*24*time.Hour),
			rotatelogs.WithRotationTime(logRotaionFlag),
		)
		logrus.SetOutput(writer)
	}
}

func SetupLogTools(productName string, enableSaveLogFile bool, settingLogLeve logrus.Level, logMaxSaveDays int64, logSaveType GTLogSaveType) {
	ProjectName = productName
	LogLevel = settingLogLeve
	LogSaveMaxDays = logMaxSaveDays
	LogSaveFlag = logSaveType
	EnableSaveLogFile = enableSaveLogFile
	Instance().Setup()
}

func Instance() *GTLog {
	GTLogOnce.Do(func() {
		ALog = &GTLog{}
		logrus.SetLevel(LogLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
		// 设置默认日志输出为控制台
		logrus.SetOutput(os.Stdout)
	})
	return ALog
}

func (aLog *GTLog) infof(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Infof(format, args...)
}

func (aLog *GTLog) warnf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Warnf(format, args...)
}

func (aLog *GTLog) errorf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Errorf(format, args...)
}
func (aLog *GTLog) debugf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Debugf(format, args...)
}
func (aLog *GTLog) tracef(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Tracef(format, args...)
}
func (aLog *GTLog) fatalf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Fatalf(format, args...)
}

// LogInfof format格式化log--info信息
func LogInfof(format string, args ...interface{}) {
	LogF(GTLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func LogErrorf(format string, args ...interface{}) {
	LogF(GTLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func LogDebugf(format string, args ...interface{}) {
	LogF(GTLogStyleDebug, format, args...)
}

// LogTracef format格式化log--Trace信息
func LogTracef(format string, args ...interface{}) {
	LogF(GTLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息
func LogFatalf(format string, args ...interface{}) {
	LogF(GTLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func LogWarnf(format string, args ...interface{}) {
	LogF(GTLogStyleWarning, format, args...)
}

// LogF 快捷日志Function，含模块字段封装
// Params [style] log类型  fatal、trace、info、warning、error、debug
// Params [format] 模块名称：自定义字符串
// Params [args...] 模块名称：自定义字符串
func LogF(style GTLogStyle, format string, args ...interface{}) {

	// 对每个占位符、非占位符片段和'['、']'进行迭代，为它们添加相应的颜色
	re := regexp.MustCompile(`(%[vTsdfqTbcdoxXUeEgGp]+)|(\[|\])|([^%\[\]]+)`)
	colorFormat := re.ReplaceAllStringFunc(format, func(s string) string {
		switch {
		case strings.HasPrefix(s, "%"):
			return gtbox_color.ANSIColorForegroundBrightYellow + s + gtbox_color.ANSIColorReset
		case s == "[" || s == "]":
			return s // 保持 `[` 和 `]` 的原始颜色
		default:
			if style == GTLogStyleError {
				return gtbox_color.ANSIColorForegroundBrightRed + s + gtbox_color.ANSIColorReset
			} else if style == GTLogStyleInfo {
				return gtbox_color.ANSIColorForegroundBrightGreen + s + gtbox_color.ANSIColorReset
			} else {
				return gtbox_color.ANSIColorForegroundBrightCyan + s + gtbox_color.ANSIColorReset
			}
		}
	})

	if style != GTLogStyleInfo {
		pc, _, _, _ := runtime.Caller(1)
		fullName := runtime.FuncForPC(pc).Name()

		lastDot := strings.LastIndex(fullName, ".")
		if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
			return
		}
		callerClass := fullName[:lastDot]
		method := fullName[lastDot+1:]

		prefixFormat := fmt.Sprintf("[pkg--%s--][method--%s--] ", callerClass, method)
		colorFormat = prefixFormat + colorFormat
	}

	switch style {
	case GTLogStyleFatal:
		Instance().fatalf(colorFormat, args...)
	case GTLogStyleTrace:
		Instance().tracef(colorFormat, args...)
	case GTLogStyleInfo:
		Instance().infof(colorFormat, args...)
	case GTLogStyleWarning:
		Instance().warnf(colorFormat, args...)
	case GTLogStyleError:
		Instance().errorf(colorFormat, args...)
	case GTLogStyleDebug:
		Instance().debugf(colorFormat, args...)
	}
}
