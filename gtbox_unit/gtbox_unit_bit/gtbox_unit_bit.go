package gtbox_unit_bit

import (
	"github.com/george012/gtbox/gtbox_unit/gtbox_unit_bytes"
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

// Covert2Bytes 以2进制计算方式 将 Bit(比特) 换算成 Byte (字节)
func (aBit *GTUnitBit) Covert2Bytes() *gtbox_unit_bytes.GTUnitBytes {
	// 初始比特值
	bits := aBit.BitValue
	returnBytes := &gtbox_unit_bytes.GTUnitBytes{}

	// 根据比特的单位进行相应的换算
	switch aBit.Unit {
	case UnitBitsBit:
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesByte
	case UnitBitsKiloBit: // Kilobit
		bits.Mul(bits, big.NewFloat(1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesKiloBytes
	case UnitBitsMegaBit: // Megabit
		bits.Mul(bits, big.NewFloat(1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesMegaBytes
	case UnitBitsGigaBit: // Gigabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesGigaBytes
	case UnitBitsTeraBit: // Terabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesTeraBytes
	case UnitBitsPetaBit: // Petabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesPetaBytes
	case UnitBitsExaBit: // Exabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesExaBytes
	case UnitBitsZettaBit: // Zettabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesZettaBytes
	case UnitBitsYottaBit: // Yottabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesYottaBytes
	}

	// 将比特转换为字节
	bits.Quo(bits, big.NewFloat(8))

	// 设定字节值
	returnBytes.BytesValue = bits

	return returnBytes
}
