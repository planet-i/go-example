package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// 设置你的代理服务器地址
	proxyUrl, err := url.Parse("http://your-proxy-address:port")
	if err != nil {
		log.Fatal(err)
	}

	// 创建代理传输实例
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	// 创建http客户端，使用代理传输
	client := &http.Client{
		Transport: transport,
	}

	// 使用http客户端发起请求
	resp, err := client.Get("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 打印响应内容
	log.Println(string(body))
}
