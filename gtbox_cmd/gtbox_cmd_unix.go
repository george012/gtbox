//go:build !windows
// +build !windows

/*
Package gtbox_encoding 编码转换"目前仅支持gbk 转utf-8"
*/
package gtbox_cmd

// GetGitBashPath attempts to retrieve the installation path of Git Bash from the Windows registry.
func getWindowsGitBashPath() string {
	return ""
}
