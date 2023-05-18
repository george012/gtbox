package gtbox_unit_bytes

import (
	"github.com/george012/gtbox/gtbox_unit/gtbox_unit_bit"
	"math/big"
)

// UnitBytes Bytes 字节 显示单位 Model
type UnitBytes int

const (
	UnitBytesByte       UnitBytes = iota // Byte 字节
	UnitBytesKiloBytes                   // KB or KBytes 千字节
	UnitBytesMegaBytes                   // MB or MBytes 兆字节
	UnitBytesGigaBytes                   // GB or GBytes 吉字节
	UnitBytesTeraBytes                   // TB or TBytes 太字节
	UnitBytesPetaBytes                   // PB or PBytes 拍字节
	UnitBytesExaBytes                    // EB or EBytes 艾字节
	UnitBytesZettaBytes                  // ZB or ZBytes 泽字节
	UnitBytesYottaBytes                  // YB or YBytes 尧字
)

func (aBytes UnitBytes) String() string {
	switch aBytes {
	case UnitBytesKiloBytes:
		return "KB"
	case UnitBytesMegaBytes:
		return "MB"
	case UnitBytesGigaBytes:
		return "GB"
	case UnitBytesTeraBytes:
		return "TB"
	case UnitBytesPetaBytes:
		return "PB"
	case UnitBytesExaBytes:
		return "EB"
	case UnitBytesZettaBytes:
		return "ZB"
	case UnitBytesYottaBytes:
		return "YB"
	case UnitBytesByte:
		return "Bytes"
	default:
		return "Bytes"
	}
}

type GTUnitBytes struct {
	BytesValue *big.Float // 字节值
	Unit       UnitBytes  // 单位
}

// Covert2Bit 以2进制计算方式 将 Byte (字节)  换算成 Bit(比特)
func (aBytes *GTUnitBytes) Covert2Bit() *gtbox_unit_bit.GTUnitBit {
	// 初始字节值
	bytes := aBytes.BytesValue
	returnBit := &gtbox_unit_bit.GTUnitBit{}

	// 将字节转换为比特
	bytes.Mul(bytes, big.NewFloat(8))

	// 根据字节的单位进行相应的换算
	switch aBytes.Unit {
	case UnitBytesByte: // Byte
		returnBit.Unit = gtbox_unit_bit.UnitBitsBit
	case UnitBytesKiloBytes: // KiloBytes
		bytes.Mul(bytes, big.NewFloat(1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsKiloBit
	case UnitBytesMegaBytes: // MegaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsMegaBit
	case UnitBytesGigaBytes: // GigaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsGigaBit
	case UnitBytesTeraBytes: // TeraBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsTeraBit
	case UnitBytesPetaBytes: // PetaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsPetaBit
	case UnitBytesExaBytes: // ExaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsExaBit
	case UnitBytesZettaBytes: // ZettaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsZettaBit
	case UnitBytesYottaBytes: // YottaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsYottaBit
	}

	// 设定比特值
	returnBit.BitValue = bytes

	return returnBit
}
