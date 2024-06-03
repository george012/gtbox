/*
Package gtbox_decimal en: Decimal ToolBox, zh-cn: Decimal高精度运算常用工具
*/
package gtbox_decimal

import (
	"github.com/shopspring/decimal"
	"math/big"
)

// Decimal2BigFloat en: decimal covert to big.Float, zh-cn: decimal 转 big.Float
func Decimal2BigFloat(decimalVal decimal.Decimal) *big.Float {
	bigFloat := new(big.Float)
	bigFloat.SetString(decimalVal.String())
	return bigFloat
}

// BigFloat2Decimal en: big.Float covert to decimal, zh-cn: big.Float 转 decimal
func BigFloat2Decimal(bigFloatVal *big.Float) decimal.Decimal {
	bigFloatStr := bigFloatVal.String()
	decimalVal, _ := decimal.NewFromString(bigFloatStr)
	return decimalVal
}

// DecimalToFloat64 en: decimal covert to Float64, zh-cn: decimal 转 Float64
func DecimalToFloat64(decimalVal decimal.Decimal) float64 {
	float64Val, _ := decimalVal.Float64()
	return float64Val
}

// Float64ToDecimal en: Float64 covert to decimal, zh-cn: Float64 转 decimal
func Float64ToDecimal(float64Val float64) decimal.Decimal {
	decimalVal := decimal.NewFromFloat(float64Val)
	return decimalVal
}
