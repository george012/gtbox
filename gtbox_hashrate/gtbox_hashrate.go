package gtbox_hashrate

import "fmt"

// GTHashRateFormat 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位MH/s、GH/s等, 默认保留3位小数点
func GTHashRateFormat(hs float64) string {

	return GTHashRateFormatWithSed(hs, 3)
}

// GTHashRateFormatWithSed 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位PH/s、EH/s、ZH/s、YH/s等
// fSed 此参数标识保留机位小数点
func GTHashRateFormatWithSed(hs float64, fSed int) string {
	k := 1000.0
	m := k * 1000
	g := m * 1000
	t := g * 1000
	p := t * 1000
	e := p * 1000
	z := e * 1000
	y := z * 1000

	// 构造格式化字符串，例如 "%.2f" 或 "%.3f"，根据 fsed 的值来决定
	format := fmt.Sprintf("%%.%df", fSed)

	switch {
	case hs >= y:
		return fmt.Sprintf(format, hs/y) + " YH/s"
	case hs >= z:
		return fmt.Sprintf(format, hs/z) + " ZH/s"
	case hs >= e:
		return fmt.Sprintf(format, hs/e) + " EH/s"
	case hs >= p:
		return fmt.Sprintf(format, hs/p) + " PH/s"
	case hs >= t:
		return fmt.Sprintf(format, hs/t) + " TH/s"
	case hs >= g:
		return fmt.Sprintf(format, hs/g) + " GH/s"
	case hs >= m:
		return fmt.Sprintf(format, hs/m) + " MH/s"
	case hs >= k:
		return fmt.Sprintf(format, hs/k) + " KH/s"
	default:
		if hs <= 0 {
			return "0 H/s"
		}
		// 注意这里使用了 "%d" 而不是传入的 fsed，确保返回的是一个没有小数的整数
		return fmt.Sprintf("%.0f", hs) + " H/s"
	}
}
