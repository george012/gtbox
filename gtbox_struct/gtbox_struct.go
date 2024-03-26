package gtbox_struct

import "reflect"

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	// 我们需要传入一个结构体的指针，然后获取其元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()

		// 获取 json tag 作为 map 的键，如果没有 json tag，则使用字段名
		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}
		result[key] = value
	}

	return result
}
