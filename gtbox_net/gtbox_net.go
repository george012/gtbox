/*
Package gtbox_net 网络工具
*/
package gtbox_net

import (
	"github.com/george012/gtbox/gtbox_string"
	"net"
	"os/exec"
	"strings"
)

// GTGetLocalIPV4WithCurrentActive 获取当前活动网卡的IPV4地址
func GTGetLocalIPV4WithCurrentActive() string {
	// 尝试连接到公共地址但不发送数据
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// GTGetPublicIPV4 获取公网IP
func GTGetPublicIPV4() string {
	curl := exec.Command("curl", "https://ipinfo.io/ip")
	out, err := curl.Output()
	if err != nil {
		return ""
	}
	aStr := gtbox_string.GTBytes2String(out)
	return aStr
}

// GTGetRandomTag 获取基于公网IP的随机字符
func GTGetRandomTag() string {
	if sArr := strings.Split(GTGetPublicIPV4(), "."); len(sArr) == 4 {
		ak := sArr[len(sArr)-1] + "x" + sArr[0]
		return ak
	}
	return ""
}
