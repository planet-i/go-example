package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var InitFileDownload = "http://192.168.51.206:30805/api/as/file/InitFileDownload"
var token = "c0b7f435a059c7bbb7a91411ecde09eb"

type InitDownReq struct {
	FileID        string
	Status        int
	TranscodeType int
}

func InitDownload(param InitDownReq) (signatureUrl string, fileName string, fileSize int, err error) {
	data := map[string]interface{}{
		"fileId":        param.FileID,
		"status":        0,
		"transcodeType": 0,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Println("请求参数", string(dataJson))
	resp, err := HttpFunc(http.MethodPost, token, InitFileDownload, bytes.NewBuffer(dataJson))
	if err != nil {
		return "", "", -1, err
	}
	js := jsoniter.Get(resp)
	code := js.Get("code").ToInt()
	msg := js.Get("msg").ToString()
	if code != 200 || msg != "" {
		return
	}
	fmt.Println("请求结果", string(resp))
	signatureUrl = js.Get("result").Get("signatureUrl").ToString()
	fileName = js.Get("result").Get("fileName").ToString()
	fileSize = js.Get("result").Get("file_size").ToInt()

	return signatureUrl, fileName, fileSize, nil
}

// InitFileDownloads 签发文件下载签名
func InitFileDownloads(param InitDownReq) (signatureUrl string, fileName string, fileSize int64, err error) {
	data := map[string]interface{}{
		"fileId":        param.FileID,
		"status":        0,
		"transcodeType": 0,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Println("请求InitDownload", string(dataJson))
	body, err := HttpPost(token, InitFileDownload, "application/json; charset=utf-8", bytes.NewBuffer(dataJson))
	if err != nil {
		return
	}
	js := jsoniter.Get(body)
	code := js.Get("code").ToInt()
	msg := js.Get("msg").ToString()
	if code != 200 || msg != "" {
		return
	}
	fmt.Println("返回的结果", string(body))
	signatureUrl = js.Get("result").Get("signatureUrl").ToString()
	fileName = js.Get("result").Get("fileName").ToString()
	fileSize = js.Get("result").Get("file_size").ToInt64()
	return
}

// HttpPost 发送POST请求
func HttpPost(token, url, contentType string, bodyReader io.Reader) (result []byte, errRet error) {
	timeout := time.Minute
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Access-Token", token)
	req.Header.Set("Connection", "keep-alive")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: timeout}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
