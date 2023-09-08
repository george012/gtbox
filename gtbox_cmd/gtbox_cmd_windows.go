//go:build windows
// +build windows

/*
Package gtbox_encoding 编码转换"目前仅支持gbk 转utf-8"
*/
package gtbox_cmd

import (
	"golang.org/x/sys/windows/registry"
)

// GetGitBashPath attempts to retrieve the installation path of Git Bash from the Windows registry.
func getWindowsGitBashPath() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\GitForWindows`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	s, _, err := k.GetStringValue("InstallPath")
	if err != nil {
		return ""
	}
	return s + "\\bin\\bash.exe"
}
