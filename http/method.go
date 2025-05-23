package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	getUrl()
	get("")
}

func getUrl() {
	apiUrl := "http://127.0.0.1:9090/get"
	// URL param
	data := url.Values{}
	data.Set("name", "枯藤")
	data.Set("age", "18")
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed,err:%v\n", err)
	}
	u.RawQuery = data.Encode() // URL encode
}

func get(urlStr string) {
	if urlStr == "" {
		urlStr = "http://www.5lmh.com"
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	defer resp.Body.Close() // 程序在使用完response后必须关闭回复的主体。
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read from resp.Body failed,err:", err)
		return
	}
	fmt.Print(string(body))
}
