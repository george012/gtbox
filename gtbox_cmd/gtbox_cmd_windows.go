//go:build windows
// +build windows

/*
Package gtbox_cmd 编码转换"目前仅支持gbk 转utf-8"
*/
package gtbox_cmd

import (
	"golang.org/x/sys/windows/registry"
)

// GetWindowsGitBashPath 获取windows环境下的git-bash
func GetWindowsGitBashPath() string {
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
