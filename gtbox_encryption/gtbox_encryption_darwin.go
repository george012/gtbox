//go:build darwin
// +build darwin

package gtbox_encryption

/*
#cgo CFLAGS: -I../libs/gtgo
#cgo LDFLAGS: -L../libs/gtgo -lgtgo.dylib
#include <stdlib.h>
#include <stdio.h>
#include <string.h> // 添加这一行
extern int GT_encryptionStr(const char *inputStr,char *outputStr,const char *keyString);
extern int GT_decryptionStr(const char *inputStr, char *outputStr, const char *keyString);
extern int gt_dec(const char *inputStr, char **outputStr, const char *keyString);
extern int gt_enc(const char *inputStr, char **outputStr, const char *keyString);
extern const char* GT_getVerison();
extern int fatal_enc(const char *inputStr, char **outputStr, char *keyString);
extern int fatal_dec(const char *inputStr, char **outputStr, char *keyString);
extern char* GT_MD5(char *inputStr);
*/
import "C"
