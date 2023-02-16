package gtbox_sys

import (
	"github.com/tidwall/gjson"
	"os/exec"
	"runtime"
	"strings"
)

type GTHardInfo struct {
	CPUNumber       string `json:"cpu_number"`
	DiskNumber      string `json:"disk_number"`
	BaseBoardNumber string `json:"base_board_number"`
	BiosNumber      string `json:"bios_number"`
}

// GTGetHardInfo 获取自定义唯一硬件号
func GTGetHardInfo() *GTHardInfo {
	aHardInfo := &GTHardInfo{
		CPUNumber:       GTGetDiskNo(),
		DiskNumber:      GTGetCPUNo(),
		BaseBoardNumber: GTGetBaseBoardNo(),
		BiosNumber:      GTGetBiosNo(),
	}
	return aHardInfo
}

// GTGetHardUUIDFromMacOS 获取Mac 系统的 UUID
func GTGetHardUUIDFromMacOS() string {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPHardwareDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		as = gjson.ParseBytes(nob).Get("SPHardwareDataType.0.platform_UUID").Str
	}
	return as
}

// GTGetProvisioningHardUUIDFromMacOS 获取Mac系统的ProvisionUUID
func GTGetProvisioningHardUUIDFromMacOS() string {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPHardwareDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		as = gjson.ParseBytes(nob).Get("SPHardwareDataType.0.provisioning_UDID").Str
	}
	return as
}

// GTGetDeviceSNFromMacOS 获取Mac的SN号
func GTGetDeviceSNFromMacOS() string {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPHardwareDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		as = gjson.ParseBytes(nob).Get("SPHardwareDataType.0.serial_number").Str
	}
	return as
}

// GTGetDiskNoWithMac 获取Mac 硬盘序列号
func GTGetDiskNoWithMac() string {
	as := ""
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("system_profiler", "SPNVMeDataType", "-json", "-detailLevel", "basic")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		as = gjson.ParseBytes(nob).Get("SPNVMeDataType.0._items.0.device_serial").Str
	}
	return as
}

func GTGetDiskNo() (disk_seria_number string) {
	diskNo := ""
	if runtime.GOOS == "darwin" {
		diskNo = GTGetDiskNoWithMac()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "diskdrive", "get", "serialnumber")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return
		}
		diskNo = string(nob)
		diskNo = strings.ReplaceAll(diskNo, "\n", "")
		diskNo = strings.ReplaceAll(diskNo, "\r", "")
		diskNo = strings.ReplaceAll(diskNo, " ", "")
		diskNo = strings.ReplaceAll(diskNo, "SerialNumber", "")
		diskNo = strings.ReplaceAll(diskNo, ".", "")
	}

	return diskNo
}

func GTGetCPUNo() (cpu_id string) {
	var cpuid = ""

	if runtime.GOOS == "darwin" {
		cpuid = GTGetDeviceSNFromMacOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "cpu", "get", "processorid")

		b, err := cmd.CombinedOutput()

		if err != nil {
			return
		}
		cpuid = string(b)
		cpuid = cpuid[12 : len(cpuid)-2]
		cpuid = strings.ReplaceAll(cpuid, " ", "")
		cpuid = strings.ReplaceAll(cpuid, "\n", "")
		cpuid = strings.ReplaceAll(cpuid, "\r", "")
	}

	return cpuid
}

func GTGetBaseBoardNo() (baseboard_id string) {
	var boardNo = ""

	if runtime.GOOS == "darwin" {
		boardNo = GTGetProvisioningHardUUIDFromMacOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "baseboard", "get", "serialnumber")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return
		}
		boardNo = string(nob)
		boardNo = strings.ReplaceAll(boardNo, "\n", "")
		boardNo = strings.ReplaceAll(boardNo, "\r", "")
		boardNo = strings.ReplaceAll(boardNo, " ", "")
		boardNo = strings.ReplaceAll(boardNo, "SerialNumber", "")
	}
	return boardNo
}

func GTGetBiosNo() (bios_no string) {

	var biosNo = ""
	if runtime.GOOS == "darwin" {
		biosNo = GTGetHardUUIDFromMacOS()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "bios", "get", "serialnumber")
		nob, err := cmd.CombinedOutput()
		if err != nil {
			return
		}
		biosNo = string(nob)
		biosNo = strings.ReplaceAll(biosNo, "\n", "")
		biosNo = strings.ReplaceAll(biosNo, "\r", "")
		biosNo = strings.ReplaceAll(biosNo, " ", "")
		biosNo = strings.ReplaceAll(biosNo, "SerialNumber", "")
	}
	return biosNo
}
