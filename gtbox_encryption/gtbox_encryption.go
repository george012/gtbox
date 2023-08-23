/*
Package gtbox_encryption 加密库
*/
package gtbox_encryption

/*
#cgo CFLAGS: -I../libs/gtgo
#cgo LDFLAGS: -L../libs/gtgo -lgtgo
#include <stdlib.h>
#include <stdio.h>
#include <string.h> // 添加这一行
extern int GT_encryptionStr(const char *inputStr,char *outputStr,const char *keyString);
extern int GT_decryptionStr(const char *inputStr, char *outputStr, const char *keyString);
*/
import "C"

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

// GTMD5Encoding MD5加密
func GTMD5Encoding(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GTEncryptionGo 自定义加密
func GTEncryptionGo(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString) // 创建C字符串

	defer C.free(unsafe.Pointer(srcCasting))
	defer C.free(unsafe.Pointer(cKeyString)) // 释放C字符串

	resultMemSize := C.strlen(srcCasting) * 2

	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize + 8 + 2)))
	C.memset(unsafe.Pointer(CStr), 0, resultMemSize+8+2)
	C.GT_encryptionStr(srcCasting, CStr, cKeyString) // 使用C字符串

	retStr := C.GoString(CStr)
	defer C.free(unsafe.Pointer(CStr))

	return retStr
}

// GTDecryptionGo 自定义解密
func GTDecryptionGo(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString)       // 创建C字符串
	defer C.free(unsafe.Pointer(srcCasting)) // Free srcCasting memory when function returns
	defer C.free(unsafe.Pointer(cKeyString)) // Free cKeyString memory when function returns

	resultMemSize := C.strlen(srcCasting) * 2

	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize)))
	C.memset(unsafe.Pointer(CStr), 0, resultMemSize)
	C.GT_decryptionStr(srcCasting, CStr, cKeyString)

	retStr := C.GoString(CStr)
	defer C.free(unsafe.Pointer(CStr))

	return retStr
}

// GTEncryptionGoReturnStringLength 自定义加密,并返回长度
func GTEncryptionGoReturnStringLength(srcString string, keyString string) (resultStr string, stringLength int) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString)       // 创建C字符串
	defer C.free(unsafe.Pointer(srcCasting)) // Free srcCasting memory when function returns
	defer C.free(unsafe.Pointer(cKeyString)) // Free cKeyString memory when function returns

	resultMemSize := C.strlen(srcCasting) * 2

	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize + 8 + 2)))
	C.memset(unsafe.Pointer(CStr), 0, resultMemSize+8+2)
	CStrLength := int(C.GT_encryptionStr(srcCasting, CStr, cKeyString))

	retStr := C.GoString(CStr)
	defer C.free(unsafe.Pointer(CStr))

	return retStr, CStrLength
}

// GTDecryptionGoWithLength 自定义解密,传入长度
func GTDecryptionGoWithLength(srcString string, keyString string, stringLength int) (resultStr string) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString)       // 创建C字符串
	defer C.free(unsafe.Pointer(srcCasting)) // Free srcCasting memory when function returns
	defer C.free(unsafe.Pointer(cKeyString)) // Free cKeyString memory when function returns

	CStr := (*C.char)(unsafe.Pointer(C.malloc(C.size_t(stringLength))))
	C.memset(unsafe.Pointer(CStr), 0, C.size_t(stringLength))
	C.GT_decryptionStr(srcCasting, CStr, cKeyString)

	retStr := C.GoString(CStr)
	defer C.free(unsafe.Pointer(CStr))

	return retStr
}
