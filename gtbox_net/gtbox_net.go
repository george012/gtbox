package gtbox_net

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_string"
	"os/exec"
	"strings"
)

// GTGetPublicIPV4 获取公网IP
func GTGetPublicIPV4() *string {
	curl := exec.Command("curl", "https://ipinfo.io/ip")
	out, err := curl.Output()
	if err != nil {
		fmt.Println("erorr", err)
		return nil
	}
	aStr := gtbox_string.GTBytes2String(out)
	return &aStr
}

// GTGetRandomTag 获取基于公网IP的随机字符
func GTGetRandomTag() string {
	if sArr := strings.Split(*GTGetPublicIPV4(), "."); len(sArr) == 4 {
		ak := sArr[len(sArr)-1] + "x" + sArr[0]
		return ak
	}
	return ""
}
