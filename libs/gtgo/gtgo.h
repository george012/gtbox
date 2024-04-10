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

    void scrypt_1024_1_1_256(const char* input, char* output);
    void scrypt_1024_1_1_256_sp(const char* input, char* output, char* scratchpad);
    const int scrypt_scratchpad_size = 131583;
#ifdef __cplusplus // 如果是 C++ 编译环境
} // 结束 extern "C" 块
#endif

#endif // GTGO_H  // 包含保护结束
