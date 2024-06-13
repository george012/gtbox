package gtbox_unit_bytes

import (
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
