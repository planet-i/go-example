package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func main() {
	imagePath := "https://img-va.myshopline.com/image/store/1699869479480/Frog-Family-Wooden-Jigsaw-Puzzle.jpeg?w=1080&h=1080"
	ctid, err := GenerateNewImage(imagePath)
	if err != nil {
		fmt.Println("GenerateNewImage:", err)
	}
	fmt.Println("ctid", ctid)
	url, err := GetNewImageInfo(ctid)
	fmt.Println("url", url)
	if err != nil {
		fmt.Println("GetNewImageInfo:", err)
	}
}

type CreateImageResp struct {
	C int `json:"c"`
	D struct {
		Ctid int64 `json:"ctid"`
	} `json:"d"`
	M string `json:"m"`
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

type PullImageReq struct {
	Link string `json:"link"`
}

func GenerateNewImage(imageUrl string) (int64, error) {
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

	//req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,or;q=0.7")
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
	fmt.Sprintln(string(body))
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

func GetNewImageInfo(ctid int64) (string, error) {
	var newUrl string
	url := fmt.Sprintf("https://imgsrv2.umcasual.com/api/v2/get_creative?ctid=%d", ctid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return newUrl, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return newUrl, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newUrl, err
	}
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

func GetFilenameFromUrl(url string) (filename, filename2 string) {
	arr := strings.Split(url, "/")
	if len(arr) >= 1 {
		filename = arr[len(arr)-1]

		filename2 = strings.Split(filename, ".")[0]
	}
	return
}

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

// 测试函数
func TestGenerateNewImage(t *testing.T) {
	// 调用函数
	imagePath := "https://img-va.myshopline.com/image/store/1699869479480/Frog-Family-Wooden-Jigsaw-Puzzle.jpeg?w=1080&h=1080"
	result, err := GenerateNewImage(imagePath)
	fmt.Println(result)
	// 验证结果
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestGetNewImage(t *testing.T) {
	// 调用函数
	var ctid int64 = 67192694769465
	result, err := GetNewImageInfo(ctid)
	fmt.Println(result)
	// 验证结果
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
