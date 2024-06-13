package gtbox_sys_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_sys"
	"testing"
)

func TestNameGetHardInfo(t *testing.T) {
	hardInfo := gtbox_sys.GTGetHardInfo()
	fmt.Printf("aHardInfo:%s", hardInfo)
}
