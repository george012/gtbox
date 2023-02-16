package gtbox_data

import (
	"bytes"
	"encoding/json"
)

// GTStructToMap struct 转 map
func GTStructToMap(content interface{}) map[string]interface{} {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		//d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name
}
