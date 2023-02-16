#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int GT_encryptionStr(const char *inputStr,char *outputStr,const char *keyString);
int GT_decryptionStr(const char *inputStr,char *outputStr,const char *keyString);


#ifdef __cplusplus
}
#endif