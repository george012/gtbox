package gtbox_unit_bit

import (
	"math/big"
)

// UnitBits Bit 比特 显示单位 Model
type UnitBits int

const (
	UnitBitsBit      UnitBits = iota // bit
	UnitBitsKiloBit                  // Kbit 千比特
	UnitBitsMegaBit                  // Mbit 兆比特
	UnitBitsGigaBit                  // Gbit 吉比特
	UnitBitsTeraBit                  // Tbit 太比特
	UnitBitsPetaBit                  // Pbit 拍比特
	UnitBitsExaBit                   // Ebit 艾比特
	UnitBitsZettaBit                 // Zbit 泽比特
	UnitBitsYottaBit                 // Ybit 尧比特

)

func (aBits UnitBits) String() string {
	switch aBits {
	case UnitBitsKiloBit:
		return "Kbit"
	case UnitBitsMegaBit:
		return "Mbit"
	case UnitBitsGigaBit:
		return "Gbit"
	case UnitBitsTeraBit:
		return "Tbit"
	case UnitBitsPetaBit:
		return "Pbit"
	case UnitBitsExaBit:
		return "Ebit"
	case UnitBitsZettaBit:
		return "Zbit"
	case UnitBitsYottaBit:
		return "Ybit"
	case UnitBitsBit:
		return "bit"
	default:
		return "bit"
	}
}

type GTUnitBit struct {
	BitValue *big.Float // 比特值
	Unit     UnitBits   // 单位
}
