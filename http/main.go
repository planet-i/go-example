package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

func getUrlSize(url string) string {
	url = "https://imgcdn.umcasual.com/creative/648299135952750/1742461697.glb"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "发送 GET 请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "错误状态码: %s\n", resp.Status)
		os.Exit(1)
	}

	// 获取 Content-Length 头信息
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		fmt.Fprintln(os.Stderr, "响应头中未找到 Content-Length")
		os.Exit(1)
	}
	return contentLength
}

func url() {
	// 目标URL
	url := "https://cdn.shopify.com/3d/models/o/e38516a2d6824caa/watch.glb"
	size := getUrlSize2(url)
	fmt.Println(size)
	// 准备表单字段
	fields := map[string]string{
		"GoogleAccessId": "video-production@video-production-225115.iam.gserviceaccount.com",
		"key":            "c/o/v/123bbb4321f4d40a101mi1fd3c32aa7.mp4",
		"policy":         "eyJjb25kaXRpb25zIjpbWyJlcSIsIiRidWNrZXQiLCJzaG9waWZ5LXZpZGVvLWRldmVsb3BtZW50LWdlbmVyYWwtb3JpZ2luYWxzIl0sWyJlcSIsIiRrZXkiLCJkZXYvby92L2Y5MzdlZmM0MDExZjRkNDBhMTAxYWY4ZWQzYzU2Y2U3Lm1wNCJdLFsikj23423kj123kjahsdbaxNSw4NTA2MTVdXSwiZXhwaXJhdGlvbiI6IjIwMjItMDgtMDFUMjM6NTM6MjNaIn0=",
		"signature":      "vD7N/vHO4MS0EpG,ms@DSF@sfsdlkasn21D5+AuQXP2naBXU1mTr7K9EelXXufl/52lDvzgxJmQvgpUWVZ9tmNtxMjEj7uiL7dUzTs1vxQC7G/fWODk43bzX54Q6Xe2+BgBNp+fK4p9zM51+XZS9SrHcoTVaoqmGdYSWtu+ABOKRObQAf5hVm6AjKphB0hqWHxfLyk+/9MCnpXjdJzUrzNDnOAMVQYV7sBBNXS123123DuLQDn7lH8CFImsC3AVnB4nGoZpV2JhPko0teoogw7umfXrRZYB8NeTr2bOdnsFzJYdlXZvhbgUW3BjDQ==",
	}

	// 文件路径
	filePath := "/Users/shopifyemployee/watches_comparison.mp4"

	// 创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	// 添加文件字段
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		panic(fmt.Sprintf("Error creating form file: %v", err))
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(fmt.Sprintf("Error copying file content: %v", err))
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		panic(fmt.Sprintf("Error closing writer: %v", err))
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(fmt.Sprintf("Error creating request: %v", err))
	}

	// 设置Content-Type（包含boundary）
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error sending request: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nResponse: %s\n", resp.StatusCode, string(respBody))
}

func getUrlSize2(url string) int64 {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0
	}
	defer resp.Body.Close()

	if resp.ContentLength != -1 {
		return resp.ContentLength
	}

	var totalBytes int64
	buffer := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buffer)
		totalBytes += int64(n)
		if err != nil {
			break
		}
	}

	return totalBytes
}

type ApiCall struct {
	REQ_URL   string
	APP_KEY   string
	APP_TOKEN string
}

func NewApiCall() *ApiCall {
	return &ApiCall{
		REQ_URL:   "https://gwapi.mabangerp.com/api/v2",
		APP_KEY:   "200636",                           // Replace with your APP_KEY
		APP_TOKEN: "58ab5a9a70484ab9b0209b012636c384", // Replace with your APP_TOKEN
	}
}

func HMAC256(c, key string) string {
	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(c))
	return hex.EncodeToString(sig.Sum(nil))
}

// API Request Call
func (api *ApiCall) Call(apiName string, reqParams map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"api":       apiName,
		"appkey":    api.APP_KEY,
		"data":      reqParams,
		"timestamp": int(time.Now().Unix()),
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %v", err)
	}

	// Generate Authorization header using HMAC-SHA256
	authorization := HMAC256(string(dataJSON), api.APP_TOKEN)

	// Prepare headers
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Requested-With": "XMLHttpRequest",
		"Authorization":    authorization,
	}

	// Use Resty client to send POST request
	client := resty.New()
	resp, err := client.R().
		SetHeaders(headers).
		SetBody(dataJSON).
		Post(api.REQ_URL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

// func main() {
// 	// Create new API client
// 	apicall := NewApiCall()

// 	// Make the API call
// 	result, err := apicall.Call("stock-do-search-sku-list-new", map[string]interface{}{})
// 	if err != nil {
// 		log.Fatalf("API call failed: %v", err)
// 	}

// 	// Print the result
// 	fmt.Println(result)
// }
