/*
Package gtbox_log Log工具
*/
package gtbox_log

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_color"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ansiEscapeRe 匹配 ANSI 色码,落文件前 strip 用。
var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSIWriter 包一个 io.Writer,写之前去掉 ANSI 色码;Both 模式 file sink 用,避免 \x1b[36m 这种入文件。
type stripANSIWriter struct{ inner io.Writer }

func (w *stripANSIWriter) Write(p []byte) (int, error) {
	_, err := w.inner.Write(ansiEscapeRe.ReplaceAll(p, nil))
	return len(p), err
}

var (
	currentLogConfig *GTLogConf
	logConfigOnce    sync.Once
	setupComplete    bool
	mainLog          *GTLog
	mainLogOnce      sync.Once
)

// GTLogStyle 日志样式
type GTLogStyle string

const (
	GTLogStyleDebug   GTLogStyle = "debug"   // Debug
	GTLogStyleError   GTLogStyle = "error"   // Error
	GTLogStyleWarning GTLogStyle = "warning" // Warning
	GTLogStyleInfo    GTLogStyle = "info"    // Info
	GTLogStyleTrace   GTLogStyle = "trace"   // Trace
	GTLogStyleFatal   GTLogStyle = "fatal"   // Fatal
	GTLogStylePanic   GTLogStyle = "panic"   // Panic
)

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

type GTLogConf struct {
	productName       string
	productLogDir     string
	enableSaveLogFile bool
	keepStdout        bool // enableSaveLogFile=true 时是否同时保留 stdout(双输出);=false 时无效
	logLeve           GTLogStyle
	logMaxSaveDays    int64
	logSaveType       GTLogSaveType
}

func instanceConfig() *GTLogConf {
	logConfigOnce.Do(func() {
		currentLogConfig = &GTLogConf{}
	})
	return currentLogConfig
}

// setupDefaultLog 懒加载默认 mainLog。所有包级 LogXxxf 入口都经它取 logger，多 goroutine
// 首次并发调日志时会同时跑到这里——故用 mainLogOnce 收口「check-then-act」，保证 mainLog 只
// 被创建一次、并发安全（对齐文件内已有的 logConfigOnce 风格）。原裸 if 判断（无锁读写包级
// mainLog/setupComplete）在并发下是 data race（go test -race 可复现）。
// setupComplete 语义不变：它只在 SetupLogTools 被置 false（且全局从不置 true），故等价于
// 「mainLog==nil 时懒建」，Once 内保留该判断不改变既有行为。
func setupDefaultLog() *GTLog {
	mainLogOnce.Do(func() {
		if setupComplete == false && mainLog == nil {
			mainLog = NewGTLog(
				strings.ToLower(instanceConfig().productName),
			)
		}
	})
	return mainLog
}

type GTLog struct {
	sync.RWMutex
	saveFileEnabled bool
	keepStdout      bool           // saveFileEnabled=true 时是否同时保留 stdout(双输出)
	logger          *logrus.Logger // 添加这一行
	modelName       string
	logDir          string
	logDirWithDate  string
	entryTime       time.Time // 日志初始化时间,留作后续比对使用
	lastCheckTime   time.Time // 记录最后一次检查时间,用作日志轮转
}

func GetProjectName() string {
	return instanceConfig().productName
}

func GetLogLevel() GTLogStyle {
	return instanceConfig().logLeve
}

func GetProductMainLogDir() string {
	return instanceConfig().productLogDir
}

// logF 快捷日志Function，含模块字段封装
// Params [style] log类型  fatal、trace、info、warning、error、debug
// Params [format] 模块名称：自定义字符串
// Params [args...] 模块名称：自定义字符串
func (aLog *GTLog) logF(style GTLogStyle, format string, args ...interface{}) {
	aLog.Lock()
	defer aLog.Unlock()

	colorFormat := format
	// stdout 在场(只 stdout 或 stdout+file 双输出)就上色;只 file 时不上色保文件干净。
	// Both 模式 stdout 拿带 ANSI 的 string;file sink 在 SetOutput 时已 wrap stripANSIWriter,落盘前 strip 掉。
	if aLog.saveFileEnabled == false || aLog.keepStdout {
		// 对每个占位符、非占位符片段和'['、']'进行迭代，为它们添加相应的颜色
		re := regexp.MustCompile(`(%[vTsdfqTbcdoxXUeEgGp]+)|(\[|\])|([^%\[\]]+)`)
		colorFormat = re.ReplaceAllStringFunc(format, func(s string) string {
			switch {
			case strings.HasPrefix(s, "%"):
				return fmt.Sprintf("%s%s%s", gtbox_color.ANSIColorForegroundBrightYellow, s, gtbox_color.ANSIColorReset)
			case s == "[" || s == "]":
				return s // 保持 `[` 和 `]` 的原始颜色
			default:
				if style == GTLogStyleError {
					return fmt.Sprintf("%s%s%s", gtbox_color.ANSIColorForegroundBrightRed, s, gtbox_color.ANSIColorReset)
				} else if style == GTLogStyleInfo {
					return fmt.Sprintf("%s%s%s", gtbox_color.ANSIColorForegroundBrightGreen, s, gtbox_color.ANSIColorReset)
				} else {
					return fmt.Sprintf("%s%s%s", gtbox_color.ANSIColorForegroundBrightCyan, s, gtbox_color.ANSIColorReset)
				}
			}
		})
	}

	if style != GTLogStyleInfo {
		pc, _, _, _ := runtime.Caller(2)
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
		aLog.logger.Fatalf(colorFormat, args...)
	case GTLogStyleTrace:
		aLog.logger.Tracef(colorFormat, args...)
	case GTLogStyleInfo:
		aLog.logger.Infof(colorFormat, args...)
	case GTLogStyleWarning:
		aLog.logger.Warnf(colorFormat, args...)
	case GTLogStyleError:
		aLog.logger.Errorf(colorFormat, args...)
	case GTLogStyleDebug:
		aLog.logger.Debugf(colorFormat, args...)
	case GTLogStylePanic:
		aLog.logger.Panicf(colorFormat, args...)
	}
}

// LogInfof format格式化log--info信息
func (aLog *GTLog) LogInfof(format string, args ...interface{}) {
	aLog.logF(GTLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func (aLog *GTLog) LogErrorf(format string, args ...interface{}) {
	aLog.logF(GTLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func (aLog *GTLog) LogDebugf(format string, args ...interface{}) {
	aLog.logF(GTLogStyleDebug, format, args...)
}

// LogTracef format格式化log--Trace信息
func (aLog *GTLog) LogTracef(format string, args ...interface{}) {
	aLog.logF(GTLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息 !!!慎用，使用后程序会退出!!!
func (aLog *GTLog) LogFatalf(format string, args ...interface{}) {
	aLog.logF(GTLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func (aLog *GTLog) LogWarnf(format string, args ...interface{}) {
	aLog.logF(GTLogStyleWarning, format, args...)
}

// determineRotationTime 辅助函数：根据日志保存类型决定轮转时间
func determineRotationTime(logSaveType GTLogSaveType) time.Duration {
	switch logSaveType {
	case GTLogSaveHours:
		return time.Hour
	case GTLogSaveTypeDays:
		return time.Hour * 24
	default:
		return time.Hour * 24 // 默认按天轮转
	}
}

func newLogSaveHandler(gtLog *GTLog) (rotateLogger *rotatelogs.RotateLogs) {
	// 确保日志目录存在
	err := os.MkdirAll(gtLog.logDir, 0755)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	logFilePath := fmt.Sprintf("%s/run", gtLog.logDirWithDate)
	linkLogFilePath := fmt.Sprintf("%s/run.log", gtLog.logDir)

	/* 日志轮转相关函数
	   `WithLinkName` 为最新的日志建立软连接
	   `WithRotationTime` 设置日志分割的时间，隔多久分割一次
	   WithMaxAge 和 WithRotationCount二者只能设置一个
	    `WithMaxAge` 设置文件清理前的最长保存时间
	    `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	writer, err := rotatelogs.New(
		logFilePath+".%Y-%m-%d_%H.log",
		rotatelogs.WithLinkName(linkLogFilePath),
		rotatelogs.WithMaxAge(time.Duration(instanceConfig().logMaxSaveDays)*24*time.Hour),
		rotatelogs.WithRotationTime(determineRotationTime(instanceConfig().logSaveType)),
	)
	if err != nil {
		// 处理错误
		fmt.Println("Error setting up log writer:", err)
		return nil
	}
	return writer
}

// applyOutput 按三态装配 logger 输出:
//
//	saveFileEnabled=false                 → 只 stdout
//	saveFileEnabled=true,keepStdout=false → 只 file
//	saveFileEnabled=true,keepStdout=true  → stdout + file(file 端 strip ANSI)
//
// rotate 句柄创建失败(如日志目录不可写:容器只读 rootfs、CI runner 无 /var/log 权限)时
// 降级为只 stdout 并置 saveFileEnabled=false——log 库绝不能因落盘不可用把宿主进程打崩;
// 修复前 nil 句柄被直接包进 MultiWriter,后台维护 goroutine 首次写日志即 SIGSEGV。
// 并发约束:仅允许 NewGTLog(goroutine 启动前)或持有 aLog 锁的调用方使用。
func (aLog *GTLog) applyOutput() {
	if !aLog.saveFileEnabled {
		aLog.logger.SetOutput(os.Stdout)
		return
	}
	rLog := newLogSaveHandler(aLog)
	if rLog == nil {
		fmt.Printf("[gtbox_log] log file sink unavailable, fallback to stdout only [dir=%s]\n", aLog.logDir)
		aLog.saveFileEnabled = false
		aLog.logger.SetOutput(os.Stdout)
		return
	}
	if aLog.keepStdout {
		aLog.logger.SetOutput(io.MultiWriter(os.Stdout, &stripANSIWriter{inner: rLog}))
	} else {
		aLog.logger.SetOutput(rLog)
	}
}

// NewGTLog 添加GTLog模块
func NewGTLog(modelName string) *GTLog {
	currentTime := time.Now().UTC()

	gtLog := &GTLog{
		modelName:      modelName,
		logDir:         fmt.Sprintf("%s/%s", instanceConfig().productLogDir, modelName),
		logDirWithDate: fmt.Sprintf("%s/%s/%s", instanceConfig().productLogDir, modelName, currentTime.Format("2006-01-02")),
		logger:         logrus.New(),
		entryTime:      currentTime,
		lastCheckTime:  currentTime,
	}

	// 初始化日志设置（代码简化，具体初始化逻辑可以根据需要调整）
	gtLog.logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05Z07:00",
	})

	gtLog.logger.SetLevel(logrus.TraceLevel)

	// 根据LogLevel设置logrus的日志级别
	switch currentLogConfig.logLeve {
	case GTLogStyleFatal:
		gtLog.logger.SetLevel(logrus.FatalLevel)
	case GTLogStyleTrace:
		gtLog.logger.SetLevel(logrus.TraceLevel)
	case GTLogStyleInfo:
		gtLog.logger.SetLevel(logrus.InfoLevel)
	case GTLogStyleWarning:
		gtLog.logger.SetLevel(logrus.WarnLevel)
	case GTLogStyleError:
		gtLog.logger.SetLevel(logrus.ErrorLevel)
	case GTLogStyleDebug:
		gtLog.logger.SetLevel(logrus.DebugLevel)
	default:
		gtLog.logger.SetLevel(logrus.InfoLevel)
	}

	gtLog.saveFileEnabled = instanceConfig().enableSaveLogFile
	gtLog.keepStdout = instanceConfig().keepStdout

	gtLog.applyOutput()

	// 启动日志维护 Goroutine，首次执行完成后继续初始化操作
	gtLog.startLogMaintenance(func(done chan struct{}) {
		// 在这里执行首次操作，比如检查日志目录和初始化逻辑
		gtLog.LogInfof("[%s] log policy started\n", gtLog.modelName)
		// 通知首次任务完成
		close(done)
	})

	return gtLog
}

// LogInfof format格式化log--info信息
func LogInfof(format string, args ...interface{}) {
	setupDefaultLog().logF(GTLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func LogErrorf(format string, args ...interface{}) {
	setupDefaultLog().logF(GTLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func LogDebugf(format string, args ...interface{}) {
	setupDefaultLog().logF(GTLogStyleDebug, format, args...)
}

// LogTracef format格式化log--Trace信息
func LogTracef(format string, args ...interface{}) {
	setupDefaultLog().logF(GTLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息 !!!慎用，使用后程序会退出!!!
func LogFatalf(format string, args ...interface{}) {
	setupDefaultLog().logF(GTLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func LogWarnf(format string, args ...interface{}) {

	setupDefaultLog().logF(GTLogStyleWarning, format, args...)
}

// SetupLogTools 初始化日志。keepStdout 仅在 enableSaveLogFile=true 时有效:true=stdout+file 双输出,false=只 file。
func SetupLogTools(productName string, enableSaveLogFile bool, keepStdout bool, logLeve GTLogStyle, logMaxSaveDays int64, logSaveType GTLogSaveType, productLogDir string) {
	setupComplete = false

	instanceConfig().productName = productName
	instanceConfig().enableSaveLogFile = enableSaveLogFile
	instanceConfig().keepStdout = keepStdout
	instanceConfig().logLeve = logLeve
	instanceConfig().logMaxSaveDays = logMaxSaveDays
	instanceConfig().logSaveType = logSaveType
	instanceConfig().productLogDir = productLogDir

	if productLogDir == "" {
		if runtime.GOOS == "linux" {
			instanceConfig().productLogDir = fmt.Sprintf("%s/%s", "/var/log", strings.ToLower(instanceConfig().productName))
		} else {
			instanceConfig().productLogDir = "./logs"
		}
	}

	if mainLog == nil {
		setupDefaultLog()
	}
}
