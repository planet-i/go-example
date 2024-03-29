package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 写一个go请求接口的典型示例
func HttpFunc(method, token, url string, body io.Reader) ([]byte, error) {
	// 生成请求实例
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	// 设置头部
	req.Header.Set("Access-Token", token)
	// 自定义client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	fmt.Println("实际请求结构", req)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("返回的请求结果", resp)
	// 返回请求结果
	return ioutil.ReadAll(resp.Body)
}
