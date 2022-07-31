package tcps

import (
	"fmt"
	"net"
)

var clients map[string]*Client

func init() {
	clients = make(map[string]*Client)
}

func Start() {
	port := 12345
	listenAddr := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("初始化失败", err.Error())
		return
	}
	fmt.Println("服务器开始运行, 端口", port)
	var tcpListener, ok = listen.(*net.TCPListener)
	if !ok {
		fmt.Println("listen error")
		return
	}
	for {
		con, err := tcpListener.AcceptTCP()
		if err != nil {
			return
		}
		client := NewClient(con)
		client.Run()

		clients[client.Uuid] = client
	}
}
