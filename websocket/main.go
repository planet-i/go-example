package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 创建一个 upgrader 实例，将其用于升级 HTTP 连接为 WebSocket 连接。
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main1() {
	// 注册服务器事件
	http.HandleFunc("/getTime", SendTime) // 支持多端同时访问,相当于收到不同请求后未每个请求分别创建线程处理
	fmt.Println("已向 DefaultServeMux 注册 SendTime 事件")

	// 打开要侦听的前端请求连接端口
	err := http.ListenAndServe(":13000", nil) // 这里是接受对本机所有的通过13000的请求,本机有多个IP则可在端口号前加IP地址区分,如"10.11.100.123:13000"
	if err != nil {
		fmt.Println("端口侦听错误:", err)
	}

}

func SendTime(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("升级websocket连接错误:", err)
		return
	}

	defer conn.Close()

	// 处理WebSocket连接
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("接收错误:", err)
	}

	fmt.Println("接收到前端消息:", string(message))

	// 发送问候
	var sendbuff string = "Hello Client!"
	conn.WriteMessage(messageType, []byte(sendbuff)) // 为阅读方便,省去错误处理

	// 不停地循环发送时间
	fmt.Println("开始发送时间")
	for {
		sendbuff = time.Now().Format("2006-01-02 15:04:05")
		err = conn.WriteMessage(messageType, []byte(sendbuff))
		if err != nil {
			fmt.Println("发送错误:", err)
			break
		}

		fmt.Println("发送时间:", sendbuff)
		time.Sleep(1 * time.Second) // 间隔1s,降低资源消耗,如果想提高实时性可减少时间间隔
	}

	fmt.Println("停止发送时间")

}
