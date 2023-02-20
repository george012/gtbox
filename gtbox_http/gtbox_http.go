package gtbox_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type GTToolsHttpRequest struct {
	HttpClient          *http.Client
	CurrentRequest      *http.Request
	CurrentResponse     *http.Response
	CurrentResponseBody []byte
	CurrentError        error
}

var (
	HttpRequest *GTToolsHttpRequest
	Once        sync.Once
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	mutex       sync.Mutex
)

func (httpReq *GTToolsHttpRequest) SetUp() {
	httpReq.HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

// Instance 单例
func Instance() *GTToolsHttpRequest {
	Once.Do(func() {
		HttpRequest = &GTToolsHttpRequest{}
		HttpRequest.SetUp()
	})

	//	每次调用请求置Nil
	HttpRequest.CurrentRequest = nil
	HttpRequest.CurrentResponse = nil
	HttpRequest.CurrentResponseBody = nil
	HttpRequest.CurrentError = nil

	return HttpRequest
}
func RequestGet(url string, successFunc func(respData []byte), errorFuc func()) {
	Instance().ToRequest(url, "", "", "GET", nil, successFunc, errorFuc)
}

// RequestPost POST请求
func RequestPost(url string, data []byte, successFunc func(respData []byte), errorFuc func()) {
	Instance().ToRequest(url, "", "", "POST", data, successFunc, errorFuc)
}

// RequestPostWithBasicAuth GET请求
func RequestPostWithBasicAuth(url string, authName string, authPwd string, data []byte, successFunc func(respData []byte), errorFuc func()) {
	Instance().ToRequest(url, authName, authPwd, "POST", data, successFunc, errorFuc)
}

func (httpReq *GTToolsHttpRequest) ToRequest(url string, authName string, authPwd string, method string, data []byte, successFunc func(respData []byte), errorFuc func()) {
	mutex.Lock()
	defer mutex.Unlock()
	//http.post等方法只是在NewRequest上又封装来了一层而已
	httpReq.CurrentRequest, httpReq.CurrentError = http.NewRequest(method, url, bytes.NewBuffer(data))
	if httpReq.CurrentError != nil {
		errors.New(fmt.Sprintf("[请求]---错误----[%s]", httpReq.CurrentError))
		errorFuc()
		return
	}
	//设置Header
	httpReq.CurrentRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

	//	设置basicAuth
	if authName != "" && authPwd != "" {
		httpReq.CurrentRequest.SetBasicAuth(authName, authPwd)
	}

	httpReq.CurrentResponse, httpReq.CurrentError = httpReq.HttpClient.Do(httpReq.CurrentRequest)
	if httpReq.CurrentError != nil {
		errors.New(fmt.Sprintf("[网络请求][返回数据]---错误----[%s]", httpReq.CurrentError))
		errorFuc()
		return
	}
	defer httpReq.CurrentResponse.Body.Close()

	httpReq.CurrentResponseBody, httpReq.CurrentError = ioutil.ReadAll(httpReq.CurrentResponse.Body)
	if httpReq.CurrentError != nil {
		errors.New(fmt.Sprintf("[读取Body]---错误----[%s]", httpReq.CurrentError))
		errorFuc()
		return
	}

	temp := make(map[string]interface{}, 0)
	httpReq.CurrentError = json.Unmarshal(httpReq.CurrentResponseBody, &temp)
	if httpReq.CurrentError != nil {
		errors.New(fmt.Sprintf("[解析JSON]---错误----[%s]", httpReq.CurrentError))
		errorFuc()
		temp = nil
		return
	}
	//	测试完JSON格式后置空 可能影响性能，但是不要紧
	temp = nil
	successFunc(httpReq.CurrentResponseBody)
}
