package gtbox_http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// RequestConfig 请求配置
// Url
// ExtendHeaders 说明：
// Content-Type  or User-Agent or Authorization ...
type RequestConfig struct {
	Url           string
	ExtendHeaders map[string]string
}

// DoRequestWithPost 发送 POST 请求，支持 JSON、map、struct、string []byte
func DoRequestWithPost(reqCfg RequestConfig, data interface{}, endFunc func(resp []byte, err error)) {
	var buffer io.Reader
	contentType := "application/x-www-form-urlencoded"
	// 根据 data 类型执行不同的序列化操作
	switch d := data.(type) {
	case string:
		// 字符串直接作为 body 发送
		buffer = strings.NewReader(d)
	case []byte:
		buffer = bytes.NewReader(d)
	default:
		contentType = "application/json"
		// 其他类型，尝试将其序列化为 JSON
		jsonData, err := json.Marshal(d)
		if err != nil {
			endFunc(nil, err)
			return
		}
		buffer = bytes.NewBuffer(jsonData)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", reqCfg.Url, buffer)
	if err != nil {
		endFunc(nil, err)
		return
	}

	req.Header.Set("Content-Type", contentType)

	for key, val := range reqCfg.ExtendHeaders {
		req.Header.Set(key, val)
	}

	// 发送 HTTP POST 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		endFunc(nil, err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		endFunc(nil, err)
		return
	}

	endFunc(body, nil)
}
