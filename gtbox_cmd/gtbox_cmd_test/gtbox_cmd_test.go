package gtbox_cmd_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_cmd"
	"testing"
)

func TestGTCmd_ExecuteCommands(t *testing.T) {
	cmdMap := map[string]string{
		"ifconfig": "ifconfig",
		"iostat":   "iostat",
	}
	cmdRes := gtbox_cmd.RunWith(cmdMap)

	if cmdRes != nil {
		value, err := cmdRes.Load("ifconfig")
		fmt.Printf("%v%v", value, err)
	}
}
