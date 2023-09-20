package gtbox_coding

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// GetProjectCodeLines 获取当前项目的有效代码行数
func GetProjectCodeLines() int64 {
	// 获取调用这个函数的文件位置
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return 0
	}

	// 从调用者的位置开始，不断向上查找，直到找到包含 go.mod 的目录
	currentDir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// 找到 go.mod，这应该是项目的根目录
			break
		}

		// 如果已经到达文件系统的根目录，则停止查找
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return 0
		}
		currentDir = parentDir
	}

	// 在找到的项目根目录上运行命令统计代码行数
	cmdStr := fmt.Sprintf("find %s -name \"*.go\" | xargs wc -l | tail -n 1 | awk '{print $1}'", currentDir)
	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	linesStr := strings.TrimSpace(string(out))
	lines, err := strconv.ParseInt(linesStr, 10, 64)
	if err != nil {
		return 0
	}
	return lines
}
