package gtbox_decimal

import (
	"github.com/shopspring/decimal"
	"math/big"
)

func Decimal2BigFloat(decimalVal decimal.Decimal) *big.Float {
	bigFloat := new(big.Float)
	bigFloat.SetString(decimalVal.String())
	return bigFloat
}

func BigFLoat2Decimal(bigFloatVal *big.Float) decimal.Decimal {
	bigFloatStr := bigFloatVal.String()
	decimalVal, _ := decimal.NewFromString(bigFloatStr)
	return decimalVal
}

func DecimalToFloat64(decimalVal decimal.Decimal) float64 {
	float64Val, _ := decimalVal.Float64()
	return float64Val
}

func Float64ToDecimal(float64Val float64) decimal.Decimal {
	decimalVal := decimal.NewFromFloat(float64Val)
	return decimalVal
}
