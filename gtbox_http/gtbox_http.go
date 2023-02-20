package gtbox_http

import (
	"bytes"
	"errors"
	"io"

	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
)

// HttpClient 是一个高并发的 HTTP 客户端封装
type HttpClient struct {
	Client       *http.Client
	RequestCount int
	mutex        sync.Mutex
}

// NewHttpClient 返回一个新的 HttpClient 对象
func NewHttpClient(timeout int) *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
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

// PostWithBasicAuth 带BasicAuth的Post
func PostWithBasicAuth(url string, authName string, authPwd string, data []byte, endFunc func(respData []byte, err error)) {
	customClient := NewHttpClient(5)
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

			body, err := ioutil.ReadAll(resp.Body)
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
		resp, err = hc.Client.Post(url, contentType, ioutil.NopCloser(bytes.NewReader(data)))
		if err != nil {
			// 请求失败，重试
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("http read error: " + err.Error())
		}

		return string(body), nil
	}

	return "", errors.New("http post error: max retry exceeded")
}
