package gtbox_log

import "testing"

func TestCustomLog(t *testing.T) {
	LogF(GTLogStyleInfo, "info[%v][%s]", "sdfdsfdsf", "dfsdfdfd")
	LogF(GTLogStyleDebug, "deb[---%s---]--[---%s---]", "deb", "bbbb")
	LogF(GTLogStyleError, "[rere%s]--[%s", "deb", "bbbb")

}
