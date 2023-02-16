package gtbox_number

import "math"

// GTToolsNumberFloat64ToInt64	将float64转成精确的int64
func GTToolsNumberFloat64ToInt64(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}
