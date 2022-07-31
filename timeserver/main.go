//一个定期报告时间的 TCP 服务器
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000") //创建一个 net.Listener 对象在一个网络端口上监听进来的连接
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept() //监听器的 Accept 方法被阻塞，直到有连接请求进来，然后返回 net.Conn 对象来代表一个连接。
		if err != nil {
			log.Print(err) //例如，连接中止
			continue
		}
		handleConn(conn) // 一次处理一个连接
	}
}

func handleConn(c net.Conn) {
	defer c.Close() //当写入失败时循环结束，很多时候是客户端断开连接，这时 handleconn 函数使用延迟的 Close 调用关闭自己这边的连接，然后继续等待下一个连接请求。
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return //例如，连接断开
		} //将 time.Now() 获取的当前时间发送给客户端。因为 net.Conn 满足 io.Writer 接口，所以可以直接向它进行写入。
		time.Sleep(1 * time.Second)
	}
}

//为了连接到服务器，需要一个像 nc("netcat") 这样的程序，以及一个用来操作网络连接的标准工具
