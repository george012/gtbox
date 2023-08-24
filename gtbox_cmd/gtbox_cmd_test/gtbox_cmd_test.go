package gtbox_cmd_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_cmd"
	"testing"
)

func TestGTCmd_ExecuteCommands(t *testing.T) {
	cmdMap := map[string]string{
		"ifconfig_a": "ifconfig",
		"iostat_b":   "iostat",
	}
	cmdRes := gtbox_cmd.RunWith(cmdMap)

	if cmdRes != nil {
		for cmd_key := range cmdMap {
			cmd_res, _ := cmdRes.Load(cmd_key)
			fmt.Printf("[%s]---[%v]\n\n", cmd_key, cmd_res)
		}

	}
}
