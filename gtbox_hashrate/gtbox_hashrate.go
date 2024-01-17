package gtbox_hashrate

import (
	"fmt"
	"math/big"
)

type HashRateUnitFormat int64

const (
	HashRateUnitFormatHs  HashRateUnitFormat = iota // H/s 默认值
	HashRateUnitFormatKHs                           // KH/s
	HashRateUnitFormatMHs                           // MH/s
	HashRateUnitFormatGHs                           // GH/s
	HashRateUnitFormatTHs                           // TH/s
	HashRateUnitFormatPHs                           // PH/s
	HashRateUnitFormatEHs                           // EH/s
	HashRateUnitFormatZHs                           // ZH/s
	HashRateUnitFormatYHs                           // YH/s
)

func (unitFormat HashRateUnitFormat) String() string {
	switch unitFormat {
	case HashRateUnitFormatYHs:
		return "YH/s"
	case HashRateUnitFormatZHs:
		return "ZH/s"
	case HashRateUnitFormatEHs:
		return "EH/s"
	case HashRateUnitFormatPHs:
		return "PH/s"
	case HashRateUnitFormatTHs:
		return "TH/s"
	case HashRateUnitFormatGHs:
		return "GH/s"
	case HashRateUnitFormatMHs:
		return "MH/s"
	case HashRateUnitFormatKHs:
		return "KH/s"
	case HashRateUnitFormatHs:
		return "H/s"
	default:
		return "H/s"
	}
}

type GTHashRate struct {
	Value    string
	UnitStr  string
	UnitFlag HashRateUnitFormat
}

// GTHashRate2Format 强制转换成某一个Hash单位类型
// baseHashRate 基础H/s为单位的 HashRate值
// toFormat 要转换成的HashRate计量单位类型
// fSed 保留小数点后几位
// Return[GTHashRate] 返回值
func GTHashRate2Format(baseHashRate *big.Float, toFormat HashRateUnitFormat, fSed int) *GTHashRate {
	k := big.NewFloat(1000)
	m := new(big.Float).Mul(k, k)
	g := new(big.Float).Mul(m, k)
	t := new(big.Float).Mul(g, k)
	p := new(big.Float).Mul(t, k)
	e := new(big.Float).Mul(p, k)
	z := new(big.Float).Mul(e, k)
	y := new(big.Float).Mul(z, k)

	// 构造格式化字符串，例如 "%.2f" 或 "%.3f"，根据 fsed 的值来决定
	format := fmt.Sprintf("%%.%df", fSed)

	cmHS := baseHashRate
	cmpUnit := HashRateUnitFormatHs
	cmpUnitStr := HashRateUnitFormatHs.String()

	switch toFormat {
	case HashRateUnitFormatYHs:
		cmHS = new(big.Float).Quo(cmHS, y)
		cmpUnit = HashRateUnitFormatYHs
		cmpUnitStr = HashRateUnitFormatYHs.String()
	case HashRateUnitFormatZHs:
		cmHS = new(big.Float).Quo(cmHS, z)
		cmpUnit = HashRateUnitFormatZHs
		cmpUnitStr = HashRateUnitFormatZHs.String()
	case HashRateUnitFormatEHs:
		cmHS = new(big.Float).Quo(cmHS, e)
		cmpUnit = HashRateUnitFormatEHs
		cmpUnitStr = HashRateUnitFormatEHs.String()
	case HashRateUnitFormatPHs:
		cmHS = new(big.Float).Quo(cmHS, p)
		cmpUnit = HashRateUnitFormatPHs
		cmpUnitStr = HashRateUnitFormatPHs.String()
	case HashRateUnitFormatTHs:
		cmHS = new(big.Float).Quo(cmHS, t)
		cmpUnit = HashRateUnitFormatTHs
		cmpUnitStr = HashRateUnitFormatTHs.String()
	case HashRateUnitFormatGHs:
		cmHS = new(big.Float).Quo(cmHS, g)
		cmpUnit = HashRateUnitFormatGHs
		cmpUnitStr = HashRateUnitFormatGHs.String()
	case HashRateUnitFormatMHs:
		cmHS = new(big.Float).Quo(cmHS, m)
		cmpUnit = HashRateUnitFormatMHs
		cmpUnitStr = HashRateUnitFormatMHs.String()
	case HashRateUnitFormatKHs:
		cmHS = new(big.Float).Quo(cmHS, k)
		cmpUnit = HashRateUnitFormatKHs
		cmpUnitStr = HashRateUnitFormatKHs.String()
	}
	return &GTHashRate{
		Value:    fmt.Sprintf(format, cmHS),
		UnitStr:  cmpUnitStr,
		UnitFlag: cmpUnit,
	}
}

// GTHashRateFormat 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位MH/s、GH/s等, 默认保留3位小数点
// hs Hash算力值
func GTHashRateFormat(hs *big.Float) string {
	return GTHashRateFormatWithSed(hs, 3)
}

// HashRateFormat 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位MH/s、GH/s等, 默认保留3位小数点
// hs Hash算力值
func HashRateFormat(hs *big.Float) string {
	k := big.NewFloat(1000)
	m := new(big.Float).Mul(k, k)
	g := new(big.Float).Mul(m, k)
	t := new(big.Float).Mul(g, k)
	p := new(big.Float).Mul(t, k)
	e := new(big.Float).Mul(p, k)
	z := new(big.Float).Mul(e, k)
	y := new(big.Float).Mul(z, k)

	var hsr *GTHashRate
	fSed := 3

	switch {
	case hs.Cmp(y) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatYHs, fSed)
	case hs.Cmp(z) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatZHs, fSed)
	case hs.Cmp(e) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatEHs, fSed)
	case hs.Cmp(p) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatPHs, fSed)
	case hs.Cmp(t) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatTHs, fSed)
	case hs.Cmp(g) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatGHs, fSed)
	case hs.Cmp(m) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatMHs, fSed)
	case hs.Cmp(k) >= 0:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatKHs, fSed)
	default:
		hsr = GTHashRate2Format(hs, HashRateUnitFormatHs, fSed)
	}

	return fmt.Sprintf("%s,%s", hsr.Value, hsr.UnitStr)
}

// GTHashRateFormatWithSed 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位PH/s、EH/s、ZH/s、YH/s等
// hs Hash算力值
// fSed 此参数标识保留机位小数点
func GTHashRateFormatWithSed(hs *big.Float, fSed int) string {
	k := big.NewFloat(1000)
	m := new(big.Float).Mul(k, k)
	g := new(big.Float).Mul(m, k)
	t := new(big.Float).Mul(g, k)
	p := new(big.Float).Mul(t, k)
	e := new(big.Float).Mul(p, k)
	z := new(big.Float).Mul(e, k)
	y := new(big.Float).Mul(z, k)

	// 构造格式化字符串，例如 "%.2f" 或 "%.3f"，根据 fsed 的值来决定
	format := fmt.Sprintf("%%.%df", fSed)

	switch {
	case hs.Cmp(y) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, y)) + " YH/s"
	case hs.Cmp(z) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, z)) + " ZH/s"
	case hs.Cmp(e) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, e)) + " EH/s"
	case hs.Cmp(p) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, p)) + " PH/s"
	case hs.Cmp(t) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, t)) + " TH/s"
	case hs.Cmp(g) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, g)) + " GH/s"
	case hs.Cmp(m) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, m)) + " MH/s"
	case hs.Cmp(k) >= 0:
		return fmt.Sprintf(format, new(big.Float).Quo(hs, k)) + " KH/s"
	default:
		if hs.Cmp(big.NewFloat(0)) <= 0 {
			return "0 H/s"
		}
		// 注意这里使用了 "%.0f" 而不是传入的 fsed，确保返回的是一个没有小数的整数
		return fmt.Sprintf("%.0f", hs) + " H/s"
	}
}
