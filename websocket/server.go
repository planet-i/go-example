package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "192.168.51.201:30883", Path: "/api/v1/ocean/aichat/files/_add"}
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("kkk", err)
	}
	defer conn.Close()
	// // 定义关闭连接的逻辑
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)
	// go func() {
	// 	<-interrupt
	// 	log.Println("接收到中断信号，正在关闭连接...")
	// 	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 	if err != nil {
	// 		log.Println("发送关闭消息失败:", err)
	// 	}
	// }()

	data := map[string]interface{}{
		"name":      "Alice",   // 知识库ID
		"namespace": "default", // 命名空间，默认default，可不传
		"files": []interface{}{
			map[string]interface{}{
				"file_name":     "2023121515_bofxooyzwzwyy_lv0.txt", // 文件ID
				"s3_object":     "2023121515_bofxooyzwzwyy_lv0.txt", // 文件ID
				"s3_hosts":      "192.168.51.201:50001",
				"s3_access_key": "datatom",
				"s3_secret_key": "datatom.com",
				"s3_bucket":     "dense-test", // 源码桶ID
				"s3_prefix":     "",
				" ":             "minio",
			},
		},
	}
	err = conn.WriteJSON(data)
	if err != nil {
		log.Println("发送数据失败:", err)
	}
	for {
		// 读取消息
		typeuu, p, err := conn.ReadMessage()
		fmt.Println(typeuu)
		if err != nil {
			log.Fatal("222", err)
		}
		var res map[string]interface{}
		err = json.Unmarshal(p, &res)
		if err != nil {
			log.Fatal("222", err)
		}
		fmt.Println(res["2023121515_bofxooyzwzwyy_lv0.txt"])
		log.Println("Received message:", string(p))
	}
}
