package gtbox

import (
	"github.com/george012/gtbox/gtbox_log"
	"testing"
)

func TestName(t *testing.T) {
	//	TODO 初始化gtbox及log分片
	SetupGTBox("test_gtbox", RunModeRelease, "", 3, gtbox_log.GTLogSaveHours, 60)
	gtbox_log.LogDebugf("level  %s", "aleve")
	gtbox_log.LogInfof("aaaaa%s", "qqq")
	gtbox_log.LogErrorf("bbb%s", "ccc")

}
