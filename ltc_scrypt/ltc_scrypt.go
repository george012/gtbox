package ltc_scrypt

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h> // 添加这一行
extern void scrypt_1024_1_1_256(const char* input, char* output);
extern void scrypt_1024_1_1_256_sp(const char* input, char* output, char* scratchpad);
extern const int scrypt_scratchpad_size;
*/
import "C"

import (
	"encoding/hex"
	"unsafe"
)

func CalcScryptHash(bytesInput []byte) []byte {
	// avoid of the situation of empty bytes array
	bytesInput = append(bytesInput, byte(0x0))
	input := (*C.char)(unsafe.Pointer(&bytesInput[0]))

	bytesOutput := make([]byte, 32)
	output := (*C.char)(unsafe.Pointer(&bytesOutput[0]))

	C.scrypt_1024_1_1_256(input, output)

	return bytesOutput
}

func CalcScryptHashHex(bytesInput []byte) string {
	bytesRet := CalcScryptHash(bytesInput)
	return hex.EncodeToString(bytesRet)
}
