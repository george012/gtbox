package gtbox_ssh_simple

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_ssh"
)

// SimpleAddressOfOne 单地址示例
func SimpleAddressOfOne() {
	var commands []string
	commands = append(commands, "pwd")
	for i := 0; i < 100; i++ {
		commands = append(commands, "ls")
	}

	connectInfo := &gtbox_ssh.CLIConnectInfo{
		Address:   "192.168.1.101",
		Port:      22,
		User:      "root",
		Pwd:       "root",
		LoginType: gtbox_ssh.SSHLoginTypePWD,
	}

	gtbox_ssh.GTSSHClientRun(connectInfo, commands, func(result *gtbox_ssh.SSHResultInfo, err error) {
		if err != nil {
			fmt.Printf("\nresult Error = %s", err)
		}
	})

}

// SimpleAddressOfDual 多地址示例
func SimpleAddressOfDual() {

	var commands []string
	for i := 0; i < 100; i++ {
		commands = append(commands, "ls")
	}

	var address []*gtbox_ssh.CLIConnectInfo

	address = append(address, &gtbox_ssh.CLIConnectInfo{
		Address: "192.168.1.101",
		Port:    22,
		User:    "root",
		Pwd:     "root",
	})

	address = append(address, &gtbox_ssh.CLIConnectInfo{
		Address: "192.168.1.102",
		Port:    22,
		User:    "root",
		Pwd:     "root",
	})

	address = append(address, &gtbox_ssh.CLIConnectInfo{
		Address: "192.168.1.103",
		Port:    22,
		User:    "root",
		Pwd:     "root",
	})

	gtbox_ssh.GTSSHClientRunDualAddress(address, commands, func(results []*gtbox_ssh.SSHResultInfo, err error) {

		if err != nil {
			fmt.Printf("\nresult Error = %s", err)
		}
	})
}
