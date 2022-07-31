package main

import (
	"fmt"
	"net/http"
)

func SayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}

type MyHandler struct {
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World")
}
func main() {
	//创建web服务端
	handler := MyHandler{}
	http.HandleFunc("/hello", SayHello) //HandleFunc注册一个处理器函数handler和对应的模式pattern。
	//http.ListenAndServe(" ", nil) //使用指定的监听地址和处理器启动一个HTTP服务端
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &handler,
	}
	server.ListenAndServe()
	/* 	server := http.Server{
	   		Addr:    "127.0.0.1:80",
	   		Handler: nil,
	   	}
	   	server.ListenAndServeTLS("cert.pem","key.pem") */
}
