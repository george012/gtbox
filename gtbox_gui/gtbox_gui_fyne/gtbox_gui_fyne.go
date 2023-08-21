package gtbox_gui_fyne

import (
	"encoding/json"
	"fmt"
)

// GTGetDescription 获取描述
func GTGetDescription() string {
	des_map := map[string]map[string]string{
		"GTButton": {
			"控件名": "GTButton",
			"描述":  "[自定义封装]--[Fyne]---[Button]",
		},
		"GTLabel": {
			"控件名": "GTLabel",
			"描述":  "[自定义封装]--[Fyne]---[Label]",
		},
	}

	jsonBytes, err := json.Marshal(des_map)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return ""
	}
	return string(jsonBytes)
}
