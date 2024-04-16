package gtbox_string

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"regexp"
	"strings"
	"unsafe"
)

var simpleBytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

// RandomString 获取随机字符串
// n 多少位 比如获取10位随机数即传10
func RandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// GTCheckMobile 判断是否为手机号
func GTCheckMobile(phoneNumber string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phoneNumber)
}

// GTStringSliceContain 判断是否包含元素
func GTStringSliceContain(slice []string, s string) bool {
	for _, s2 := range slice {
		if s == s2 {
			return true
		}
	}
	return false
}

// GTValidHostnamePort 检测IP+端口字符串是否正确
func GTValidHostnamePort(s string) bool {
	sp := strings.Split(s, ":")
	if len(sp) != 2 {
		return false
	}
	if sp[0] == "" || sp[1] == "" {
		return false
	}
	return true
}

func GTRecodingString(src string, srcCode string, toCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(toCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// GTUTF8String2GBKString UTF-8转GBK
func GTUTF8String2GBKString(src string) string {
	return GTRecodingString(src, "GBK", "UTF-8")
}

// GTGBKString2UTF8String	GBK转UTF-8
func GTGBKString2UTF8String(src string) string {
	return GTRecodingString(src, "GBK", "UTF-8")
}

// GTBytes2String byte转string
func GTBytes2String(BytesData []byte) string {
	return *(*string)(unsafe.Pointer(&BytesData))
}

// GTString2Bytes string转bytes
func GTString2Bytes(strData string) []byte {
	return *(*[]byte)(unsafe.Pointer(&strData))
}

// GTStruct2JsonString struct转json
func GTStruct2JsonString(value interface{}) (jsonString string) {
	var cuValue, _ = json.Marshal(value)
	jsonString = string(cuValue)
	return jsonString
}

// DelStringEndNewlines 删除字符串结尾的 \n or \r\n
func DelStringEndNewlines(s *string) {
	b := []byte(*s)
	b = bytes.TrimSuffix(b, []byte("\r\n"))
	b = bytes.TrimSuffix(b, []byte("\n"))
	*s = string(b)
}

// StringCoverBool str value is true or false  covert to golang bool type
func StringCoverBool(str string) bool {

	return str == "true"
}
