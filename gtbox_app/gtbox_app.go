package gtbox_app

import (
	"github.com/george012/gtbox"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultNetListenAddress    = "0.0.0.0"
	defaultNetListenTCPVersion = "tcp4"
	defaultHTTPRequestTimeOut  = 60 * time.Second
)

type App struct {
	sync.RWMutex            // 嵌入读写锁以支持并发读写
	AppName                 string
	BundleID                string
	Version                 string
	Description             string
	HTTPRequestTimeOut      time.Duration
	NetListenDefaultAddress string
	NetListenTCPVersion     string
	CurrentRunMode          gtbox.RunMode
	GitCommitHash           string
	GitCommitTime           string
	GoVersion               string
	PackageOS               string
	PackageTime             string
	AppRunAsDir             string
	AppLocalDBPath          string
	AppConfigFilePath       string
	AppLogPath              string
	StartTime               time.Time
}

func (app *App) getAppRunAsDir() {
	// 获取当前执行程序的路径
	exePath, _ := os.Executable()
	// 获取当前执行程序所在的文件夹
	app.AppRunAsDir = filepath.Dir(exePath)
	if app.CurrentRunMode == gtbox.RunModeDebug {
		app.AppRunAsDir = "."
	}
}

func (app *App) initializePaths() {
	app.getAppRunAsDir()
	app.AppLocalDBPath = app.AppRunAsDir + "/dts/dts"
	app.AppConfigFilePath = app.AppRunAsDir + "/conf/config.json"
	app.AppLogPath = app.AppRunAsDir + "/logs"
}

func NewApp(appName, version, bundleID, description string, runMode gtbox.RunMode) *App {
	app := &App{
		AppName:                 appName,
		BundleID:                bundleID,
		Version:                 version,
		Description:             description,
		CurrentRunMode:          runMode,
		HTTPRequestTimeOut:      defaultHTTPRequestTimeOut,
		NetListenTCPVersion:     defaultNetListenTCPVersion,
		NetListenDefaultAddress: defaultNetListenAddress,
		StartTime:               time.Now(),
	}
	app.initializePaths()
	return app
}
