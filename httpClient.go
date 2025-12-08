package WdaGo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Ning9527fff/MyLog"
)

// HTTPClient HTTP
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetRequest 发送GET请求
func (h *HTTPClient) GetRequest(url string, headers map[string]string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(" Error in create http request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" Error in send request : %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(" Error in read message from response : %v", err)
	}

	if resp.StatusCode >= 400 {
		return body, fmt.Errorf(" Error in check status code : %s", resp.Status)
	}

	return body, nil
}

// PostRequest 发送POST请求
func (h *HTTPClient) PostRequest(url string, data interface{}, headers map[string]string) ([]byte, error) {
	var body io.Reader

	// 处理请求数据
	if data != nil {
		jsonData, err := json.Marshal(data)

		log.DebugF("Request body is : %v", string(jsonData))

		if err != nil {
			return nil, fmt.Errorf(" Format json failed : %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf(" Create POST requestd failed: %v", err)
	}

	// 设置默认Content-Type
	if headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" Send POST failed : %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.DebugF("response status code is %v ", resp.StatusCode)
		log.DebugF("response body is %v ", respBody)

		return nil, fmt.Errorf(" Read  message from response failed: %v", err)
	}

	log.DebugF("response body is %v\n", string(respBody))

	// 检查状态码
	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf(" Error, http status code is not correct : %s", resp.Status)
	}

	return respBody, nil
}

// DeleteRequest 发送DELETE请求
func (h *HTTPClient) DeleteRequest(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf(" Error in Create Delete Request: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" Error in sending Delete Request : %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(" Error in reading message from response : %v", err)
	}

	// 检查状态码
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf(" Error in checking http status code : %s", resp.Status)
	}

	return body, nil
}

// 便捷函数 - 使用默认客户端

// Get 使用默认客户端发送GET请求
func Get(url string, headers map[string]string) ([]byte, error) {
	client := NewHTTPClient(30 * time.Second)
	return client.GetRequest(url, headers)
}

// Post 使用默认客户端发送POST请求
func Post(url string, data interface{}, headers map[string]string) ([]byte, error) {
	client := NewHTTPClient(30 * time.Second)
	return client.PostRequest(url, data, headers)
}

// Delete 使用默认客户端发送DELETE请求
func Delete(url string, headers map[string]string) ([]byte, error) {
	client := NewHTTPClient(30 * time.Second)
	return client.DeleteRequest(url, headers)
}
