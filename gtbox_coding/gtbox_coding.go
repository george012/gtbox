package gtbox_coding

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_cmd"
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
	// 路径转换适配 Windows
	if runtime.GOOS == "windows" {
		// 只转换驱动器字母为小写，并替换反斜杠为正斜杠
		currentDir = "/" + strings.ToLower(string(currentDir[0])) + currentDir[2:]
		currentDir = strings.Replace(currentDir, "\\", "/", -1)
	}
	// 在找到的项目根目录上运行命令统计代码行数
	cmdStr := fmt.Sprintf("find %s -name \"*.go\" | xargs wc -l | tail -n 1 | awk '{print $1}'", currentDir)

	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		win_git_bash_path := gtbox_cmd.GetWindowsGitBashPath()
		if win_git_bash_path == "" {
		} else {
			cmd = exec.Command(win_git_bash_path, "-c", cmdStr)
		}
	case "darwin":
		cmd = exec.Command("/bin/zsh", "-c", cmdStr)
	default:
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

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
