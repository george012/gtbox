package gtgo_sys_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_sys"
)

func GTGoTestSysGetHardInfo() {
	hardInfo := gtbox_sys.GTGetHardInfo()
	fmt.Printf("aHardInfo:%s", hardInfo)
}
