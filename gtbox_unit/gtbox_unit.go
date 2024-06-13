package gtbox_unit

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_unit/gtbox_unit_bit"
	"github.com/george012/gtbox/gtbox_unit/gtbox_unit_bytes"
	"math/big"
)

// UnitType
/*
* Bit 比特 显示单位 Model
 * @enum UnitTypeNone	不做转换
 * @enum UnitTypeAuto 自动根据位数换算单位
 * @enum UnitTypeKilo Kilo 千
 * @enum UnitTypeMega Mega 兆
 * @enum UnitTypeGiga Giga 吉
 * @enum UnitTypeTera Tera 太
 * @enum UnitTypePeta Peta 拍
 * @enum UnitTypeExa Exa 艾
 * @enum UnitTypeZetta Zetta 泽
 * @enum UnitTypeYotta Yotta 尧
*/
type UnitType int

const (
	UnitTypeNone  UnitType = iota // 不做转换
	UnitTypeAuto                  // 自动根据位数换算单位
	UnitTypeKilo                  // Kilo 千
	UnitTypeMega                  // Mega 兆
	UnitTypeGiga                  // Giga 吉
	UnitTypeTera                  // Tera 太
	UnitTypePeta                  // Peta 拍
	UnitTypeExa                   // Exa 艾
	UnitTypeZetta                 // Zetta 泽
	UnitTypeYotta                 // Yotta 尧

)

func (aType UnitType) String() string {
	switch aType {
	case UnitTypeKilo:
		return "Kilo"
	case UnitTypeMega:
		return "Mega"
	case UnitTypeGiga:
		return "Giga"
	case UnitTypeTera:
		return "Tera"
	case UnitTypePeta:
		return "Peta"
	case UnitTypeExa:
		return "Exa"
	case UnitTypeZetta:
		return "Zetta"
	case UnitTypeYotta:
		return "Yotta"
	case UnitTypeAuto:
		return "Auto"
	case UnitTypeNone:
		return "None"
	default:
		return "None"
	}
}

type GTUnit struct {
	BitInfo   *gtbox_unit_bit.GTUnitBit
	BytesInfo *gtbox_unit_bytes.GTUnitBytes
}

// UnitBitCovertToUnitBytes 以2进制计算方式 将 Bit(比特) 换算成 Byte (字节)
func UnitBitCovertToUnitBytes(aBit *gtbox_unit_bit.GTUnitBit) *gtbox_unit_bytes.GTUnitBytes {
	// 初始比特值
	bits := aBit.BitValue
	returnBytes := &gtbox_unit_bytes.GTUnitBytes{}

	// 根据比特的单位进行相应的换算
	switch aBit.Unit {
	case gtbox_unit_bit.UnitBitsBit:
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesByte
	case gtbox_unit_bit.UnitBitsKiloBit: // Kilobit
		bits.Mul(bits, big.NewFloat(1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesKiloBytes
	case gtbox_unit_bit.UnitBitsMegaBit: // Megabit
		bits.Mul(bits, big.NewFloat(1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesMegaBytes
	case gtbox_unit_bit.UnitBitsGigaBit: // Gigabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesGigaBytes
	case gtbox_unit_bit.UnitBitsTeraBit: // Terabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesTeraBytes
	case gtbox_unit_bit.UnitBitsPetaBit: // Petabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesPetaBytes
	case gtbox_unit_bit.UnitBitsExaBit: // Exabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesExaBytes
	case gtbox_unit_bit.UnitBitsZettaBit: // Zettabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesZettaBytes
	case gtbox_unit_bit.UnitBitsYottaBit: // Yottabit
		bits.Mul(bits, big.NewFloat(1024*1024*1024*1024*1024*1024*1024*1024))
		returnBytes.Unit = gtbox_unit_bytes.UnitBytesYottaBytes
	}

	// 将比特转换为字节
	bits.Quo(bits, big.NewFloat(8))

	// 设定字节值
	returnBytes.BytesValue = bits

	return returnBytes
}

// UnitBytesCovertToUnitBit 以2进制计算方式 将 Byte (字节)  换算成 Bit(比特)
func UnitBytesCovertToUnitBit(aBytes *gtbox_unit_bytes.GTUnitBytes) *gtbox_unit_bit.GTUnitBit {
	// 初始字节值
	bytes := aBytes.BytesValue
	returnBit := &gtbox_unit_bit.GTUnitBit{}

	// 将字节转换为比特
	bytes.Mul(bytes, big.NewFloat(8))

	// 根据字节的单位进行相应的换算
	switch aBytes.Unit {
	case gtbox_unit_bytes.UnitBytesByte: // Byte
		returnBit.Unit = gtbox_unit_bit.UnitBitsBit
	case gtbox_unit_bytes.UnitBytesKiloBytes: // KiloBytes
		bytes.Mul(bytes, big.NewFloat(1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsKiloBit
	case gtbox_unit_bytes.UnitBytesMegaBytes: // MegaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsMegaBit
	case gtbox_unit_bytes.UnitBytesGigaBytes: // GigaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsGigaBit
	case gtbox_unit_bytes.UnitBytesTeraBytes: // TeraBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsTeraBit
	case gtbox_unit_bytes.UnitBytesPetaBytes: // PetaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsPetaBit
	case gtbox_unit_bytes.UnitBytesExaBytes: // ExaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsExaBit
	case gtbox_unit_bytes.UnitBytesZettaBytes: // ZettaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsZettaBit
	case gtbox_unit_bytes.UnitBytesYottaBytes: // YottaBytes
		bytes.Mul(bytes, big.NewFloat(1024*1024*1024*1024*1024*1024*1024*1024))
		returnBit.Unit = gtbox_unit_bit.UnitBitsYottaBit
	}

	// 设定比特值
	returnBit.BitValue = bytes

	return returnBit
}

// NewWithBitFormat 通过 bit(比特) 创建 Model
// bitString bit单位的 bit值
// resultUnitType 返回值的单位 默认不做处理，推荐自动格式化
func NewWithBitFormat(bitString string, resultUnitType UnitType) *GTUnit {
	bigNum, _ := new(big.Float).SetString(bitString)
	unitType := gtbox_unit_bit.UnitBitsBit

	abitInfo := &gtbox_unit_bit.GTUnitBit{
		BitValue: bigNum,
		Unit:     unitType,
	}

	// 根据结果单位类型来进行单位换算
	switch resultUnitType {
	case UnitTypeKilo:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsKiloBit
	case UnitTypeMega:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsMegaBit
	case UnitTypeGiga:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsGigaBit
	case UnitTypeTera:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsTeraBit
	case UnitTypePeta:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsPetaBit
	case UnitTypeExa:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsExaBit
	case UnitTypeZetta:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsZettaBit
	case UnitTypeYotta:
		abitInfo.Unit = gtbox_unit_bit.UnitBitsYottaBit
	case UnitTypeAuto:
		// 递减处理，直到找到最合适的单位
		for bigNum.Cmp(big.NewFloat(1024)) >= 0 && unitType < gtbox_unit_bit.UnitBitsYottaBit {
			unitType++
			bigNum.Quo(bigNum, big.NewFloat(1024))
		}
		abitInfo.Unit = unitType
	}

	abytesInfo := UnitBitCovertToUnitBytes(abitInfo)

	return &GTUnit{
		abitInfo,
		abytesInfo,
	}
}

// NewWithBytesFormat 通过 Bytes(字节) 创建 Model
// BytesString Byte 单位的 Bytes值
// resultUnitType 返回值的单位 默认不做处理，推荐自动格式化
func NewWithBytesFormat(BytesString string, resultUnitType UnitType) *GTUnit {
	bigNum, _ := new(big.Float).SetString(BytesString)
	unitType := gtbox_unit_bytes.UnitBytesByte
	aBytes := &gtbox_unit_bytes.GTUnitBytes{
		BytesValue: bigNum,
		Unit:       unitType,
	}

	// 根据结果单位类型来进行单位换算
	switch resultUnitType {
	case UnitTypeKilo:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesKiloBytes
	case UnitTypeMega:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesMegaBytes
	case UnitTypeGiga:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesGigaBytes
	case UnitTypeTera:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesTeraBytes
	case UnitTypePeta:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesPetaBytes
	case UnitTypeExa:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesExaBytes
	case UnitTypeZetta:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesZettaBytes
	case UnitTypeYotta:
		aBytes.Unit = gtbox_unit_bytes.UnitBytesYottaBytes
	case UnitTypeAuto:
		// 递减处理，直到找到最合适的单位
		for bigNum.Cmp(big.NewFloat(1024)) >= 0 && unitType < gtbox_unit_bytes.UnitBytesYottaBytes {
			unitType++
			bigNum.Quo(bigNum, big.NewFloat(1024))
		}
		aBytes.Unit = unitType
	}

	aBits := UnitBytesCovertToUnitBit(aBytes)

	return &GTUnit{
		aBits,
		aBytes,
	}
}

// UnitFormatWith1024 单位计算公式 以byte 为单位传入，根据数值不同换算成进制单位MB、GB等
// fSed 小数点保留位数
func UnitFormatWith1024(baseWithByte *big.Float, fSed int) string {
	k := big.NewFloat(1024)
	m := new(big.Float).Mul(k, k)
	g := new(big.Float).Mul(m, k)
	t := new(big.Float).Mul(g, k)
	p := new(big.Float).Mul(t, k)
	e := new(big.Float).Mul(p, k)
	z := new(big.Float).Mul(e, k)
	y := new(big.Float).Mul(z, k)

	format := fmt.Sprintf("%%.%df", fSed)

	newStr := ""
	switch {
	case baseWithByte.Cmp(y) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, y)), "YB")
	case baseWithByte.Cmp(z) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, z)), "ZB")
	case baseWithByte.Cmp(e) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, e)), "EB")
	case baseWithByte.Cmp(p) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, p)), "PB")
	case baseWithByte.Cmp(t) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, t)), "TB")
	case baseWithByte.Cmp(g) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, g)), "GB")
	case baseWithByte.Cmp(m) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, m)), "MB")
	case baseWithByte.Cmp(k) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, k)), "KB")
	default:
		newStr = fmt.Sprintf("%s", baseWithByte)
	}

	return newStr
}
