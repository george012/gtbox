package gtbox_number

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// GTToolsNumberFloat64ToInt64	将float64转成精确的int64
func GTToolsNumberFloat64ToInt64(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// GTFloat64GetLengthSegmentLength 获取`float64`总位数、小数点前位数和小数点后位数
func GTFloat64GetLengthSegmentLength(f float64) (total, before, after int) {
	// 将 float64 转换为字符串
	str := strconv.FormatFloat(f, 'f', -1, 64)

	// 如果小数点前有负号，总位数加 1
	if strings.HasPrefix(str, "-") {
		total++
		before++
	}
	// 计算总位数、小数点前位数和小数点后位数
	for i, s := range str {
		if s == '.' {
			after = len(str) - i - 1
			break
		} else {
			total++
			before++
		}
	}
	return
}

// GetFloat64GetLengthTotal 获取`float64`总位数
func GetFloat64GetLengthTotal(f float64) int {
	total, _, _ := GTFloat64GetLengthSegmentLength(f)
	return total
}

// GTFloat64GetLengthBefore 获取`float64`小数点前位数
func GTFloat64GetLengthBefore(f float64) int {
	_, before, _ := GTFloat64GetLengthSegmentLength(f)
	return before
}

// GTFloat64GetLengthAfter 获取`float64`小数点后位数
func GTFloat64GetLengthAfter(f float64) int {
	_, _, after := GTFloat64GetLengthSegmentLength(f)
	return after
}

// GTFloat64Sum Float64超高精度运算--加法
func GTFloat64Sum(a, b float64) float64 {
	f1 := big.NewFloat(a)
	f2 := big.NewFloat(b)
	result := new(big.Float).Add(f1, f2)

	f64, _ := result.Float64()

	return f64
}

// GTFloat64Sub Float64超高精度运算--减法
func GTFloat64Sub(a, b float64) float64 {
	f1 := big.NewFloat(a)
	f2 := big.NewFloat(b)
	result := new(big.Float).Sub(f1, f2)

	f64, _ := result.Float64()

	return f64
}

// GTFloat64Mul Float64超高精度运算--乘法
func GTFloat64Mul(a, b float64) float64 {
	f1 := big.NewFloat(a)
	f2 := big.NewFloat(b)
	result := new(big.Float).Mul(f1, f2)

	f64, _ := result.Float64()

	return f64
}

// GTFloat64Div Float64超高精度运算--除法
func GTFloat64Div(a, b float64) float64 {
	f1 := big.NewFloat(a)
	f2 := big.NewFloat(b)
	result := new(big.Float).Quo(f1, f2)

	f64, _ := result.Float64()

	return f64
}
