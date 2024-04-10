package ltc_scrypt_test

import (
	"encoding/hex"
	"fmt"
	"github.com/george012/gtbox/ltc_scrypt"
	"testing"
)

func TestCalcScryptHash(t *testing.T) {
	// "0000000110c8357966576df46f3b802ca897deb7ad18b12f1c24ecff6386ebd9"
	b1, _ := hex.DecodeString("01000000f615f7ce3b4fc6b8f61e8f89aedb1d0852507650533a9e3b10b9bbcc30639f279fcaa86746e1ef52d3edb3c4ad8259920d509bd073605c9bf1d59983752a6b06b817bb4ea78e011d012d59d4")
	b1h := ltc_scrypt.CalcScryptHash(b1)
	b1r := make([]byte, len(b1h))
	for i := 0; i < len(b1h); i++ {
		b1r[i] = b1h[len(b1h)-i-1]
	}
	fmt.Println("b1r: ", hex.EncodeToString(b1r))
}

func TestCalcScryptHashHex(t *testing.T) {
	// "0000000110c8357966576df46f3b802ca897deb7ad18b12f1c24ecff6386ebd9"
	b1, _ := hex.DecodeString("01000000f615f7ce3b4fc6b8f61e8f89aedb1d0852507650533a9e3b10b9bbcc30639f279fcaa86746e1ef52d3edb3c4ad8259920d509bd073605c9bf1d59983752a6b06b817bb4ea78e011d012d59d4")
	b1hHex := ltc_scrypt.CalcScryptHashHex(b1)
	b1h, _ := hex.DecodeString(b1hHex)
	b1r := make([]byte, len(b1h))
	for i := 0; i < len(b1h); i++ {
		b1r[i] = b1h[len(b1h)-i-1]
	}
	fmt.Println("b1r: ", hex.EncodeToString(b1r))
}
