package gt_box_encryption_test

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_encryption"
	"testing"
)

func TestGTEncryptionFunctions(t *testing.T) {
	pre_en_str := "fgfgfgh"
	en_key := "test"
	alen := len(pre_en_str)
	en_Str := gtbox_encryption.GTEncryptionGo(pre_en_str, en_key)

	fmt.Printf("en_str[%s]%d", en_Str, alen)

	de_str := gtbox_encryption.GTDecryptionGo(en_Str, en_key)
	b_len := len(de_str)
	fmt.Printf("en_str[%s]%d", de_str, b_len)

}

func TestGTEnc(t *testing.T) {
	for i := 0; i < 100; i++ {
		pre_en_str := fmt.Sprintf("%d%s%d", i, "dsfsdf", i+1)
		en_key := "test"
		alen := len(pre_en_str)
		en_Str := gtbox_encryption.GTEnc(pre_en_str, en_key)

		fmt.Printf("en_str[%s]%d\n", en_Str, alen)

		de_str := gtbox_encryption.GTDec(en_Str, en_key)
		b_len := len(de_str)

		fmt.Printf("de_str[%s]%d\n", de_str, b_len)

		cd_str := gtbox_encryption.GTDecryptionGo(en_Str, en_key)
		cd_len := len(cd_str)
		fmt.Printf("en_str[%s]%d\n%s%d", de_str, b_len, cd_str, cd_len)
	}
}

func TestFatalFunc(t *testing.T) {
	for i := 0; i < 100; i++ {
		pre_en_str := fmt.Sprintf("%d%s%d", i, "dsfsdf", i+1)
		en_key := "test"
		alen := len(pre_en_str)
		en_Str := gtbox_encryption.FatalEnc(pre_en_str, en_key)

		fmt.Printf("en_str[%s]%d\n", en_Str, alen)

		de_str := gtbox_encryption.FatalDec(en_Str, en_key)
		b_len := len(de_str)

		fmt.Printf("de_str[%s]%d\n", de_str, b_len)

		cd_str := gtbox_encryption.GTDecryptionGo(en_Str, en_key)
		cd_len := len(cd_str)
		fmt.Printf("en_str[%s]%d\n%s%d", de_str, b_len, cd_str, cd_len)

		cf_str := gtbox_encryption.GTDec(en_Str, en_key)
		cf_len := len(cd_str)
		fmt.Printf("en_str[%s]%d\n%s%d\n%s%d", de_str, b_len, cd_str, cd_len, cf_str, cf_len)
	}
}

func TestGetEnLibVersion(t *testing.T) {
	avstion := gtbox_encryption.GetEncryptionLibVersion()
	fmt.Printf("%s", avstion)
}
