#ifndef GTGO_H // 包含保护开始
#define GTGO_H

#ifdef __cplusplus // 如果是 C++ 编译环境
extern "C"
{ // 开始 extern "C" 块
#endif

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
    const char* GT_getVerison();
    char* GT_MD5(char *inputStr);
    int GT_encryptionStr(const char *inputStr, char *outputStr, const char *keyString);
    int GT_decryptionStr(const char *inputStr, char *outputStr, const char *keyString);
    int gt_enc(const char *inputStr, char **outputStr, char *keyString);
    int gt_dec(const char *inputStr, char **outputStr, char *keyString);
#ifdef __cplusplus // 如果是 C++ 编译环境
} // 结束 extern "C" 块
#endif

#endif // GTGO_H  // 包含保护结束
