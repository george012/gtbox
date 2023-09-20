package gtbox_coding

import (
	"github.com/george012/gtbox/gtbox_cmd"
)

// GTGetCurrentProjectCodeLineNumber 获取当前项目的有效代码行数
func GTGetCurrentProjectCodeLineNumber() int64 {
	cmds := map[string]string{"code_line_total": "find . -name \"*.go\" | xargs wc -l | tail -n 1 | awk '{print $1}'"}
	cmdRes := gtbox_cmd.RunWith(cmds)
	if cmdRes != nil {
		for cmd_key := range cmds {
			cmd_res, _ := cmdRes.Load(cmd_key)
			if convertedValue, ok := cmd_res.(int64); ok {
				return convertedValue
			}
		}
	}
	return 0
}
