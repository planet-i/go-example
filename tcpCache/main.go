package main

import (
	"github.com/planet-i/goexample1/cacheExample/cache"
	http "github.com/planet-i/goexample1/cacheExample/http"
	tcp "github.com/planet-i/goexample1/tcpCache/tcp"
)

func main() {
	ca := cache.New("inmemory")
	go tcp.New(ca).Listen() //创建一个TCP服务监听的goroutine
	http.New(ca).Listen()
}
