package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main222() {
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "192.168.51.201:30883", Path: "/api/v1/ocean/aichat/files/_add"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("连接到WebSocket服务器失败:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// 启动协程接收消息
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("读取消息失败:", err)
				return
			}
			fmt.Printf("收到消息: %s\n", message)
		}
	}()

	// 发送数据
	data := getData()
	err = c.WriteJSON(data)
	if err != nil {
		log.Println("发送数据失败:", err)
		return
	}

	// 定义关闭连接的逻辑
	go func() {
		select {
		case <-interrupt:
			log.Println("接收到中断信号，正在关闭连接...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("发送关闭消息失败:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}()

	// 保持主线程不退出
	select {}
}

func getData() map[string]interface{} {
	data := map[string]interface{}{
		"name":      "Alice",
		"namespace": "default",
		"files":     []interface{}{},
	}

	fileList := []string{"file1.txt"}
	s3Host := "192.168.51.201:50001"
	s3AccessKey := "datatom"
	s3SecretKey := "datatom.com"
	s3Bucket := "dense-test"
	s3Backend := "minio"

	for _, fileName := range fileList {
		fileItem := map[string]interface{}{
			"file_name":       fileName,
			"s3_object":       fileName,
			"s3_hosts":        s3Host,
			"s3_access_key":   s3AccessKey,
			"s3_secret_key":   s3SecretKey,
			"s3_bucket":       s3Bucket,
			"s3_prefix":       "",
			"backend_storage": s3Backend,
		}
		data["files"] = append(data["files"].([]interface{}), fileItem)
	}

	return data
}
