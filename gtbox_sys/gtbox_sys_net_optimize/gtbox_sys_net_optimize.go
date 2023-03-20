/*
Package gtbox_sys_net_optimize 主要提供网络并发优化功能
*/
package gtbox_sys_net_optimize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Current map[string]string `json:"current"`
	New     map[string]string `json:"new"`
}

func getSystemInfo() (int, uint64) {
	cpuCount := runtime.NumCPU()

	vmStat, err := exec.Command("vmstat", "-s").Output()
	if err != nil {
		log.Fatal(err)
	}

	var memorySize uint64
	lines := strings.Split(string(vmStat), "\n")
	for _, line := range lines {
		if strings.Contains(line, "total memory") {
			fields := strings.Fields(line)
			if len(fields) == 4 {
				memorySize, _ = strconv.ParseUint(fields[0], 10, 64)
				break
			}
		}
	}

	return cpuCount, memorySize
}

func readConfig(file string, delimiter string) (map[string]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	config := make(map[string]string)
	for _, line := range lines {
		fields := strings.Split(line, delimiter)
		if len(fields) == 2 {
			key, value := strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
			config[key] = value
		}
	}

	return config, nil
}

func applyConfig(file string, delimiter string, newConfig map[string]string) error {
	config, err := readConfig(file, delimiter)
	if err != nil {
		return err
	}

	for key, value := range newConfig {
		config[key] = value
	}

	content := ""
	for key, value := range config {
		content += fmt.Sprintf("%s %s %s\n", key, delimiter, value)
	}

	err = ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func calculateConfig(cpuCount int, memorySize uint64) map[string]string {
	concurrency := cpuCount * 12800
	newConfig := map[string]string{
		"fs.file-max":                   strconv.Itoa(concurrency),
		"net.core.somaxconn":            strconv.Itoa(concurrency),
		"net.ipv4.tcp_max_syn_backlog":  strconv.Itoa(concurrency),
		"net.ipv4.ip_local_port_range":  "1024 " + strconv.Itoa(concurrency),
		"net.ipv4.tcp_tw_reuse":         "1",
		"net.ipv4.tcp_sack":             "1",
		"net.ipv4.tcp_timestamps":       "1",
		"net.ipv4.tcp_keepalive_time":   "300",
		"net.ipv4.tcp_keepalive_probes": "5",
		"net.ipv4.tcp_keepalive_intvl":  "15",
		"* soft nofile":                 strconv.Itoa(concurrency),
		"* hard nofile":                 strconv.Itoa(concurrency),
		"* soft nproc":                  strconv.Itoa(concurrency),
		"* hard nproc":                  strconv.Itoa(concurrency),
	}

	return newConfig
}

func autoExecute() {
	cpuCount, memorySize := getSystemInfo()

	limitsConfig, err := readConfig("/etc/security/limits.conf", " ")
	if err != nil {
		log.Fatal(err)
	}

	sysctlConfig, err := readConfig("/etc/sysctl.conf", "=")
	if err != nil {
		log.Fatal(err)
	}

	newConfig := calculateConfig(cpuCount, memorySize)

	// 打印当前配置
	config := Config{
		Current: make(map[string]string),
		New:     make(map[string]string),
	}
	for key := range newConfig {
		if value, ok := limitsConfig[key]; ok {
			config.Current[key] = value
		} else if value, ok := sysctlConfig[key]; ok {
			config.Current[key] = value
		}
		config.New[key] = newConfig[key]
	}
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println("Current and new configurations:")
	fmt.Println(string(jsonConfig))

	// 应用新配置
	err = applyConfig("/etc/security/limits.conf", " ", newConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = applyConfig("/etc/sysctl.conf", "=", newConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 重新加载 sysctl 配置
	_, err = exec.Command("sysctl", "-p").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Configuration updated successfully.")
}
