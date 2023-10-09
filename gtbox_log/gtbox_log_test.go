package gtbox_log

import "testing"

func TestCustomLog(t *testing.T) {
	LogDebugf("aadebug  %s", "debug")
	LogInfof("info %s", "info")
	LogErrorf("error %s", "error")
}
