package gtbox_http

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	testUrlGet := "https://postman-echo.com/get"
	testUrlPost := "https://postman-echo.com/post"

	Get(testUrlGet, func(respData []byte, err error) {
		if err != nil {
			fmt.Printf("test [GET] https request failed, error: %s\n", err.Error())
			return
		}

		// 尝试将返回的数据解析为 JSON
		var parsedData interface{}
		if jsonErr := json.Unmarshal(respData, &parsedData); jsonErr != nil {
			fmt.Printf("Failed to parse response as JSON. Raw response: %s, error: %s\n", string(respData), jsonErr.Error())
			return
		}

		// 格式化输出 JSON 数据
		ajData, _ := json.MarshalIndent(parsedData, "", "    ")
		fmt.Printf("test [GET] https request, received data:\n%s\n", ajData)
	})

	Post(testUrlPost, nil, func(respData []byte, err error) {
		if err != nil {
			fmt.Printf("test [GET] https request failed, error: %s\n", err.Error())
			return
		}

		// 尝试将返回的数据解析为 JSON
		var parsedData interface{}
		if jsonErr := json.Unmarshal(respData, &parsedData); jsonErr != nil {
			fmt.Printf("Failed to parse response as JSON. Raw response: %s, error: %s\n", string(respData), jsonErr.Error())
			return
		}

		// 格式化输出 JSON 数据
		ajData, _ := json.MarshalIndent(parsedData, "", "    ")
		fmt.Printf("test [POST] https request, received data:\n%s\n", ajData)
	})
}
