package gtbox_hashrate

import (
	"fmt"
	"math/big"
)

// GTHashRateFormat 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位MH/s、GH/s等, 默认保留3位小数点
// hs Hash算力值
func GTHashRateFormat(hs *big.Float) string {

	return GTHashRateFormatWithSed(hs, 3)
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
