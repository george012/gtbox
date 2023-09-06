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
extern int gt_dec(const char *inputStr, char **outputStr, const char *keyString);
extern int gt_enc(const char *inputStr, char **outputStr, const char *keyString);
extern const char* GT_getVerison();
*/
import "C"

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

// GetEncryptionLibVersion 获取加密库版本
func GetEncryptionLibVersion() (version string) {
	return C.GoString(C.GT_getVerison())
}

// GTMD5EncryptionGo MD5加密
func GTMD5EncryptionGo(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GTEncryptionGo 自定义加密
func GTEncryptionGo(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	defer C.free(unsafe.Pointer(srcCasting))

	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(cKeyString))

	resultMemSize := C.strlen(srcCasting) * 2
	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize + 8 + 2)))
	if CStr == nil {
		return "" // 或者返回一个错误
	}
	defer C.free(unsafe.Pointer(CStr))

	C.memset(unsafe.Pointer(CStr), 0, resultMemSize+8+2)
	C.GT_encryptionStr(srcCasting, CStr, cKeyString)

	return C.GoString(CStr)
}

// GTDecryptionGo 自定义解密
func GTDecryptionGo(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	defer C.free(unsafe.Pointer(srcCasting))

	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(cKeyString))

	resultMemSize := C.strlen(srcCasting) * 2
	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize)))
	if CStr == nil {
		return "" // 或者返回一个错误
	}
	defer C.free(unsafe.Pointer(CStr))

	C.memset(unsafe.Pointer(CStr), 0, resultMemSize)
	C.GT_decryptionStr(srcCasting, CStr, cKeyString)

	return C.GoString(CStr)
}

// GTEncryptionGoReturnStringLength 自定义加密,并返回长度
func GTEncryptionGoReturnStringLength(srcString string, keyString string) (resultStr string, stringLength int) {
	srcCasting := C.CString(srcString)
	defer C.free(unsafe.Pointer(srcCasting))

	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(cKeyString))

	resultMemSize := C.strlen(srcCasting) * 2
	CStr := (*C.char)(unsafe.Pointer(C.malloc(resultMemSize + 8 + 2)))
	if CStr == nil {
		return "", 0 // 或者返回一个错误
	}
	defer C.free(unsafe.Pointer(CStr))

	C.memset(unsafe.Pointer(CStr), 0, resultMemSize+8+2)
	CStrLength := int(C.GT_encryptionStr(srcCasting, CStr, cKeyString))

	return C.GoString(CStr), CStrLength
}

// GTDecryptionGoWithLength 自定义解密,传入长度
func GTDecryptionGoWithLength(srcString string, keyString string, stringLength int) (resultStr string) {
	srcCasting := C.CString(srcString)
	defer C.free(unsafe.Pointer(srcCasting))

	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(cKeyString))

	CStr := (*C.char)(unsafe.Pointer(C.malloc(C.size_t(stringLength))))
	if CStr == nil {
		return "" // 或者返回一个错误
	}
	defer C.free(unsafe.Pointer(CStr))

	C.memset(unsafe.Pointer(CStr), 0, C.size_t(stringLength))
	C.GT_decryptionStr(srcCasting, CStr, cKeyString)

	return C.GoString(CStr)
}

// GTEnc 加密
func GTEnc(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(srcCasting))
	defer C.free(unsafe.Pointer(cKeyString))

	var output *C.char
	length := int(C.gt_enc(srcCasting, &output, cKeyString))

	if length <= 0 {
		return "" // 或者处理错误
	}

	result := C.GoString(output)
	C.free(unsafe.Pointer(output))

	return result
}

// GTDec 解密
func GTDec(srcString string, keyString string) (resultStr string) {
	srcCasting := C.CString(srcString)
	cKeyString := C.CString(keyString)
	defer C.free(unsafe.Pointer(srcCasting))
	defer C.free(unsafe.Pointer(cKeyString))

	var output *C.char
	length := int(C.gt_dec(srcCasting, &output, cKeyString))

	if length <= 0 {
		return "" // 或者处理错误
	}

	result := C.GoString(output)
	C.free(unsafe.Pointer(output))

	return result
}
