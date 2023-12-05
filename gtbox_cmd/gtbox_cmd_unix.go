//go:build !windows
// +build !windows

/*
Package gtbox_cmd 编码转换"目前仅支持gbk 转utf-8"
*/
package gtbox_cmd

// GetWindowsGitBashPath 获取windows环境下的git-bash
func GetWindowsGitBashPath() string {
	return ""
}
