/*
Package gtbox_reflect 反射方法工具库
*/
package gtbox_reflect

import (
	"reflect"
)

// GetNumberOfFieldWithModel 获取Model的字段数量支持指针、非指针传入
func GetNumberOfFieldWithModel(customModel interface{}) int {
	typ := reflect.TypeOf(customModel)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Struct {
		numFields := typ.NumField()
		return numFields
	} else {
		return 0
	}
}
