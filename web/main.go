package main

import (
	"fmt"
	"net/http" //能让程序与HTTP进行交互
)

func main() {
	http.HandleFunc("/", handler)     //把定义的handler函数设置为/被访问时的处理器
	http.ListenAndServe(":8080", nil) //启动服务，让它监听8080端口
}

//从request结构中提取相关信息，创建一个HTTP响应，通过ResponseWriter接口将响应返回给客户端
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, " Hello World,%s!", request.URL.Path[1:]) //对I/O进行格式化
}
