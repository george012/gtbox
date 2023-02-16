package gtbox_hashrate

import "fmt"

// GTHashrateFormat 哈希显示单位计算公式 以H/s 为单位传入，根据数值不同换算成进制单位MH/s、GH/s等,[fsed] 此参数标识保留机位小数点
func GTHashrateFormat(hs float64) string {
	k := 1000.0
	m := k * 1000
	g := m * 1000
	t := g * 1000

	if hs > t {
		return fmt.Sprintf("%.2f", hs/t) + " TH/s"
	}
	if hs > g {
		return fmt.Sprintf("%.2f", hs/g) + " GH/s"
	}
	if hs > m {
		return fmt.Sprintf("%.2f", hs/m) + " MH/s"
	}
	if hs > k {
		return fmt.Sprintf("%.2f", hs/k) + " KH/s"
	}

	return "0 H/s"
}
