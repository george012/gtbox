package gtbox_common

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetBinaryRunAsDir 获取 二进制程序 执行目录
func GetBinaryRunAsDir() string {
	// 获取当前执行程序的路径
	exePath, _ := os.Executable()
	// 获取当前执行程序所在的文件夹
	runAsDir := filepath.Dir(exePath)
	return runAsDir
}

// GetFileRunAsDir 获取 源代码文件 所在目录
func GetFileRunAsDir() string {
	// 获取当前执行程序的路径
	_, filePath, _, _ := runtime.Caller(2)
	// 获取当前执行程序所在的文件夹
	runAsDir := filepath.Dir(filePath)
	return runAsDir
}

// IsRunningFromTemp 判断程序是否在临时目录下运行
// 判断程序是否在临时目录下运行
func IsRunningFromTemp() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	// 获取可执行文件的父目录
	exeDir := filepath.Dir(exePath)

	// 获取系统的临时目录
	tempDir := os.TempDir()

	// 使用 filepath.Rel 来检查 exeDir 是否是 tempDir 的子路径
	relPath, err := filepath.Rel(tempDir, exeDir)
	if err != nil {
		return false
	}

	// 如果 relPath 以 ".." 开头，说明 exeDir 不在 tempDir 内
	isAt := strings.HasSuffix(relPath, "..")
	return !isAt
}
