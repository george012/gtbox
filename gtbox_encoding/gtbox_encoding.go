/*
Package gtbox_encoding 编码转换"目前仅支持gbk 转utf-8"
*/
package gtbox_encoding

import (
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io"
	"strings"
	"syscall"
	"unsafe"
)

const (
	LOCALE_IDEFAULTANSICODEPAGE = 0x00001004
)

// ConvertToUTF8UsedLocalENV 将本地获取到的字符串强制转为UTF-8。
func ConvertToUTF8UsedLocalENV(str string) (string, error) {
	codePage := getDefaultAnsiCodePage()
	enc := getEncodingByCodePage(codePage)
	if enc == nil {
		return "", fmt.Errorf("unsupported code page: %s", codePage)
	}

	reader := transform.NewReader(strings.NewReader(str), enc.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(utf8Data), nil
}

func getDefaultAnsiCodePage() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getUserDefaultLCID := kernel32.NewProc("GetUserDefaultLCID")
	getLocaleInfo := kernel32.NewProc("GetLocaleInfoW")

	lcid, _, _ := getUserDefaultLCID.Call()

	buf := make([]uint16, 6)
	getLocaleInfo.Call(lcid, LOCALE_IDEFAULTANSICODEPAGE, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func getEncodingByCodePage(codePage string) encoding.Encoding {
	switch codePage {
	case "936":
		return simplifiedchinese.GBK
	case "1252":
		return charmap.Windows1252
	case "950":
		return traditionalchinese.Big5

	// 添加其他对应的编码...
	default:
		return nil
	}
}
