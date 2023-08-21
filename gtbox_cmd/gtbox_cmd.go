package gtbox_cmd

import (
	"bytes"
	"fmt"
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
		cmd = exec.Command("cmd", "/C", command)
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

	if gcmd.results == nil {
		gcmd.results = &sync.Map{}
	}

	gcmd.results.Store(key, result)
}
