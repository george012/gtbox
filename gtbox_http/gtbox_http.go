/*
Package gtbox_http http客户端工具
*/
package gtbox_http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	DefaultTimeout = 30 * time.Second // default second 30s
)

// HttpClient 是一个高并发的 HTTP 客户端封装
type HttpClient struct {
	Client       *http.Client
	RequestCount int
	mutex        sync.Mutex
}

// NewHttpClient 返回一个新的 HttpClient 对象
func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get 发送带 Header 的 GET 请求
func (hc *HttpClient) Get(url string) ([]byte, error) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("http new request error: " + err.Error())
	}
	req.Header.Add("User-Agent", UserAgent)

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, errors.New("http do error: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("http read error: " + err.Error())
	}

	return body, nil
}

// GetWithTimeOut Get请求,timeout 超时时间
func GetWithTimeOut(url string, timeout time.Duration, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(timeout)
	resp, err := customClient.Get(url)
	endFunc(resp, err)
}

// Get Get请求
func Get(url string, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(DefaultTimeout)
	resp, err := customClient.Get(url)
	endFunc(resp, err)
}

// PostWithTimeOut Post请求,timeout 超时时间
func PostWithTimeOut(url string, timeout time.Duration, data []byte, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(timeout)
	resp, err := customClient.Post(url, "", "", data)
	endFunc(resp, err)
}

// Post Post请求
func Post(url string, data []byte, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(DefaultTimeout)
	resp, err := customClient.Post(url, "", "", data)
	endFunc(resp, err)
}

// PostWithBasicAuth 带BasicAuth的Post
func PostWithBasicAuth(url string, authName string, authPwd string, data []byte, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(DefaultTimeout)
	resp, err := customClient.Post(url, authName, authPwd, data)
	endFunc(resp, err)
}

// Post 发送带 Header 的 POST 请求
func (hc *HttpClient) Post(url string, authName string, authPwd string, data []byte) ([]byte, error) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, errors.New("http new request error: " + err.Error())
	}

	//	设置basicAuth
	if len(authName) > 1 && len(authPwd) > 1 {
		req.SetBasicAuth(authName, authPwd)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", UserAgent)
	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, errors.New("http do error: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("http read error: " + err.Error())
	}

	return body, nil
}

// GetWithRetry 带重试的 GET 请求，当请求失败时会自动重试，直到达到最大重试次数
func (hc *HttpClient) GetWithRetry(url string, maxRetry int) (string, error) {
	var err error
	var resp *http.Response
	for i := 0; i <= maxRetry; i++ {
		hc.mutex.Lock()
		resp, err = hc.Client.Get(url)
		if err == nil {
			hc.mutex.Unlock()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", errors.New("http read error: " + err.Error())
			}
			return string(body), nil
		}

		hc.RequestCount++
		hc.mutex.Unlock()
		time.Sleep(time.Duration(1) * time.Second)
	}

	return "", errors.New("http get error: " + err.Error())
}

// PostWithRetry 带重试的 Post 请求，当请求失败时会自动重试，直到达到最大重试次数
func (hc *HttpClient) PostWithRetry(url string, contentType string, data []byte, maxRetry int) (string, error) {
	var err error
	var resp *http.Response
	for i := 0; i <= maxRetry; i++ {
		resp, err = hc.Client.Post(url, contentType, io.NopCloser(bytes.NewReader(data)))
		if err != nil {
			// 请求失败，重试
			continue
		}
		defer resp.Body.Close()

		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return "", errors.New("http read error: " + err.Error())
		}

		return string(body), nil
	}

	return "", errors.New("http post error: max retry exceeded")
}
