/*
Package gtbox_cmd 本地命令行工具
*/
package gtbox_cmd

import (
	"bytes"
	"fmt"
	"github.com/george012/gtbox/gtbox_encoding"
	"github.com/george012/gtbox/gtbox_string"
	"os/exec"
	"runtime"
	"sync"
)

var rc_wg sync.WaitGroup

type gtCmd struct {
	results *sync.Map
}

func RunWith(CommandMap map[string]string) *sync.Map {
	gcmd := &gtCmd{}

	for key, command := range CommandMap {
		rc_wg.Add(1)
		go gcmd.execute(key, command)
	}
	rc_wg.Wait()
	return gcmd.results
}

func (gcmd *gtCmd) execute(key string, command string) {
	defer rc_wg.Done()
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		bit_bahs_path := getWindowsGitBashPath()
		if bit_bahs_path == "" {
		} else {
			cmd = exec.Command(bit_bahs_path, "-c", command)
		}
	case "darwin":
		cmd = exec.Command("/bin/zsh", "-c", command)
	default:
		cmd = exec.Command("/bin/sh", "-c", command)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}

	result := out.String()
	if runtime.GOOS == "windows" {
		result, _ = gtbox_encoding.ConvertToUTF8UsedLocalENV(result)
	}

	if gcmd.results == nil {
		gcmd.results = &sync.Map{}
	}
	gtbox_string.DelStringEndNewlines(&result)

	gcmd.results.Store(key, result)
}
