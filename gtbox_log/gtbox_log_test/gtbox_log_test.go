package gtbox_log_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_log"
	"testing"
)

func TestCustomLog(t *testing.T) {
	gtbox_log.SetupLogTools("testP", true, gtbox_log.GTLogStyleTrace, 3, gtbox_log.GTLogSaveHours, "")

	for q := 0; q < 10; q++ {
		gtbox_log.LogDebugf("main_test %d", q)
		gtbox_log.LogInfof("main_test %d", q)
		gtbox_log.LogWarnf("main_test %d", q)
		gtbox_log.LogTracef("main_test %d", q)
		gtbox_log.LogErrorf("main_test %d", q)
	}

	for i := 0; i < 4; i++ {
		a_loger := gtbox_log.NewGTLog(fmt.Sprintf("at_%d", i))
		for j := 0; j < 5; j++ {
			a_loger.LogDebugf("at_ %d", j)
			a_loger.LogInfof("at_ %d", j)
			a_loger.LogWarnf("at_ %d", j)
			a_loger.LogTracef("at_ %d", j)
			a_loger.LogErrorf("at_ %d", j)

		}
	}
}
