package gtbox_cmd_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_cmd"
	"testing"
)

func TestGTCmd_ExecuteCommands(t *testing.T) {
	cmdMap := map[string]string{
		"git_commit_hash": "git show -s --format=%H",
		"git_commit_time": "git show -s --format=\"%ci\" | cut -d ' ' -f 1,2",
		"build_os":        "go env GOOS",
		"go_version":      "go version | awk '{print $3}'",
	}

	cmdRes := gtbox_cmd.RunWith(cmdMap)

	if cmdRes != nil {
		for cmd_key := range cmdMap {

			cmd_res, _ := cmdRes.Load(cmd_key)
			fmt.Printf("[%s]---[%v] \n\n", cmd_key, cmd_res)
		}

	}
}
