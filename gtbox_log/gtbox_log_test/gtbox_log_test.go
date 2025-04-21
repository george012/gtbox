package gtbox_log_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_log"
	"testing"
)

func TestCustomLog(t *testing.T) {
	gtbox_log.SetupLogTools("testP", true, gtbox_log.GTLogStyleTrace, 3, gtbox_log.GTLogSaveHours, "")

	for q := 0; q < 1; q++ {
		gtbox_log.LogDebugf("main_test %d debug", q)
		gtbox_log.LogInfof("main_test %d info", q)
		gtbox_log.LogWarnf("main_test %d warning", q)
		gtbox_log.LogTracef("main_test %d trace", q)
		gtbox_log.LogErrorf("main_test %d error", q)
	}

	for i := 0; i < 2; i++ {
		a_loger := gtbox_log.NewGTLog(fmt.Sprintf("at_%d_all", i))
		for j := 0; j < 1; j++ {
			a_loger.LogDebugf("at_ %d debug", j)
			a_loger.LogInfof("at_ %d info", j)
			a_loger.LogWarnf("at_ %d warning", j)
			a_loger.LogTracef("at_ %d trace", j)
			a_loger.LogErrorf("at_ %d error", j)

		}
	}
}
