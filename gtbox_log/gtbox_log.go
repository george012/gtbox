/*
Package gtbox_log Log工具
*/
package gtbox_log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
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
	ALog           *GTLog
	GTLogOnce      sync.Once
	ProjectName    = "test"
	LogLevel       = logrus.DebugLevel
	LogSaveMaxDays int64
	LogSaveFlag    = GTLogSaveTypeDays
	LogPath        = "./logs/run"
	LogDebugToCut  = false //debug模式下是否开启日志分割	默认false方便IDE调试
)

// GTGetLogsDir 获取Log目录
func GTGetLogsDir() string {
	return LogPath
}

func (alog *GTLog) Setup() {

	//	设置Log
	if LogDebugToCut == true {

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

func SetupLogTools(productName string, debugToCut bool, settingLogLeve logrus.Level, logMaxSaveDays int64, logSaveType GTLogSaveType) {
	ProjectName = productName
	LogLevel = settingLogLeve
	LogSaveMaxDays = logMaxSaveDays
	LogSaveFlag = logSaveType
	LogDebugToCut = debugToCut
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

func (aLog *GTLog) Ininfof(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Infof(format, args...)
}

func (aLog *GTLog) Warnf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Warnf(format, args...)
}

func (aLog *GTLog) Errorf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Errorf(format, args...)
}
func (aLog *GTLog) Debugf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Debugf(format, args...)
}
func (aLog *GTLog) Tracef(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Tracef(format, args...)
}
func (aLog *GTLog) Fatalf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Fatalf(format, args...)
}

// LogInfof format格式化log--info信息
func LogInfof(format string, args ...interface{}) {
	Instance().Ininfof(format, args...)
}

// LogErrorf format格式化log--error信息
func LogErrorf(format string, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
		return
	}
	callerClass := fullName[:lastDot]
	method := fullName[lastDot+1:]

	endForMat := fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	Instance().Errorf(endForMat, args...)
}

// LogDebugf format格式化log--debug信息
func LogDebugf(format string, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
		return
	}
	callerClass := fullName[:lastDot]
	method := fullName[lastDot+1:]

	endForMat := fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	Instance().Debugf(endForMat, args...)
}

// LogTracef format格式化log--Trace信息
func LogTracef(format string, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
		return
	}
	callerClass := fullName[:lastDot]
	method := fullName[lastDot+1:]

	endForMat := fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	Instance().Tracef(endForMat, args...)
}

// LogFatalf format格式化log--Fatal信息
func LogFatalf(format string, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
		return
	}
	callerClass := fullName[:lastDot]
	method := fullName[lastDot+1:]

	endForMat := fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	Instance().Fatalf(endForMat, args...)
}

// LogWarnf format格式化log--Warning信息
func LogWarnf(format string, args ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
		return
	}
	callerClass := fullName[:lastDot]
	method := fullName[lastDot+1:]

	endForMat := fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	Instance().Warnf(endForMat, args...)
}

// LogF 快捷日志Function，含模块字段封装
// Params [style] log类型  fatal、trace、info、warning、error、debug
// Params [format] 模块名称：自定义字符串
// Params [args...] 模块名称：自定义字符串
func LogF(style GTLogStyle, format string, args ...interface{}) {
	green := "\033[32m"
	red := "\033[31m"
	reset := "\033[0m"

	endForMat := green + format + reset // 绿色

	// 使用红色为每个参数上色
	coloredArgs := make([]interface{}, len(args))
	for i, arg := range args {
		coloredArgs[i] = red + fmt.Sprintf("%v", arg) + reset
	}

	if style != GTLogStyleInfo {
		pc, _, _, _ := runtime.Caller(1)
		fullName := runtime.FuncForPC(pc).Name()

		lastDot := strings.LastIndex(fullName, ".")
		if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
			return
		}
		callerClass := fullName[:lastDot]
		method := fullName[lastDot+1:]

		endForMat = fmt.Sprintf("[pkg--%s--][method--%s--] [%v]", callerClass, method, format)
	}

	switch style {
	case GTLogStyleFatal:
		Instance().Fatalf(endForMat, coloredArgs...)
	case GTLogStyleTrace:
		Instance().Tracef(endForMat, coloredArgs...)
	case GTLogStyleInfo:
		Instance().Ininfof(endForMat, coloredArgs...)
	case GTLogStyleWarning:
		Instance().Warnf(endForMat, coloredArgs...)
	case GTLogStyleError:
		Instance().Errorf(endForMat, coloredArgs...)
	case GTLogStyleDebug:
		Instance().Debugf(endForMat, coloredArgs...)
	}
}
