package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CreateImageResp struct {
	C int `json:"c"`
	D struct {
		Ctid int64 `json:"ctid"`
	} `json:"d"`
	M string `json:"m"`
}

type PullImageReq struct {
	Link string `json:"link"`
}

// GenerateNewImage 函数用于向指定 API 发送图片链接，生成新图片并返回其 CTID
// 参数 imageUrl 是图片的远程链接
// 返回值为生成图片的 CTID 和可能出现的错误
func GenerateNewImage(imageUrl string) (int64, error) {
	// 初始化 CTID 变量，用于存储生成图片的 CTID
	var ctid int64
	url := "https://imgsrv2.umcasual.com/api/v2/creative_creative"

	reqData := map[string]string{
		"link": imageUrl,
	}
	reqJs, err := json.Marshal(reqData)
	if err != nil {
		return ctid, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJs))
	if err != nil {
		return ctid, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctid, err
	}
	defer resp.Body.Close()

	// 输出响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ctid, err
	}
	fmt.Printf("GenerateNewImage response body: %s\n", string(body))
	if resp.StatusCode == http.StatusOK {
		var respData CreateImageResp
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return ctid, err
		}
		ctid = respData.D.Ctid
	}
	return ctid, nil
}

type GetImageResp struct {
	C int `json:"c"`
	D struct {
		Creative struct {
			URL string `json:"url"`
		} `json:"creative"`
	} `json:"d"`
	M string `json:"m"`
}

// GetNewImageInfo 函数根据 CTID 从指定 API 获取新图片的 URL
// 参数 ctid 是生成图片的 CTID
// 返回值为新图片的 URL 和可能出现的错误
func GetNewImageInfo(ctid int64) (string, error) {
	// 初始化新图片 URL 变量，用于存储获取到的新图片 URL
	var newUrl string
	url := fmt.Sprintf("https://imgsrv2.umcasual.com/api/v2/get_creative?ctid=%d", ctid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return newUrl, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return newUrl, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newUrl, err
	}
	fmt.Printf("GetNewImageInfo response body: %s\n", string(body))
	if resp.StatusCode == http.StatusOK {
		var respData GetImageResp
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return newUrl, err
		}
		newUrl = respData.D.Creative.URL
	}
	return newUrl, nil
}

// ReadRemoteFile 函数从指定的远程 URL 读取文件内容
// 参数 url 是文件的远程链接
// 返回值为文件内容的字节切片和可能出现的错误
// 可以通过这个，读取图片的数据然后导出到文件的时候显示图片
func ReadRemoteFile(url string) ([]byte, error) {
	result := make([]byte, 0)
	if url == "" {
		return result, errors.New("URL不能为空")
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// 使用一个缓冲区
	buffer := bytes.NewBuffer(make([]byte, 0, 2048))
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
