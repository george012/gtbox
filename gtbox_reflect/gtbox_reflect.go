/*
Package gtbox_reflect 反射方法工具库
*/
package gtbox_reflect

import (
	"fmt"
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

// GetFieldNameAtIndex 获取指定索引位置的--字段名
func GetFieldNameAtIndex(customModel interface{}, index int) (string, error) {
	typ := reflect.TypeOf(customModel)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Struct {
		if index >= 0 && index < typ.NumField() {
			return typ.Field(index).Name, nil
		}
		return "", fmt.Errorf("Index out of range")
	}
	return "", fmt.Errorf("Provided interface is not a struct or pointer to struct")
}

// GetFieldValueAtIndex 获取指定索引位置的字段--Value值
func GetFieldValueAtIndex(customModel interface{}, index int) (interface{}, error) {
	val := reflect.ValueOf(customModel)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		if index >= 0 && index < val.NumField() {
			return val.Field(index).Interface(), nil
		}
		return nil, fmt.Errorf("Index Out Of Range")
	}
	return nil, fmt.Errorf("Provided interface is not a struct or pointer to struct")
}
