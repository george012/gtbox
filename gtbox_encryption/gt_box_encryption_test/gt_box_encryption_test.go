package gt_box_encryption_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_encryption"
	"testing"
)

func TestGTEncryptionFunctions(t *testing.T) {
	pre_en_str := "abcdd"
	en_key := "test"
	en_Str := gtbox_encryption.GTEncryptionGo(pre_en_str, en_key)

	fmt.Printf("en_str[%s]", en_Str)

	de_str := gtbox_encryption.GTDecryptionGo(en_Str, en_key)

	fmt.Printf("en_str[%s]", de_str)

}
