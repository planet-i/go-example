package main

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/planet-i/go-example/pool"
)

const (
	maxGoroutines  = 25
	pooledResourse = 2
)

//实现了io.Closer接口的资源
type dbConnection struct {
	ID int32
}

func (dbConn *dbConnection) Close() error {
	log.Println("Close: Connection", dbConn.ID)
	return nil
}

var idCounter int32 // 用于生成资源的唯一ID
func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create:New Connection", id)
	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines) //设置为将要执行的goroutine数量
	// 创建用来管理连接的池
	p, err := pool.New(createConnection, pooledResourse)
	if err != nil {
		log.Println(err)
	}
	//使用池中的连接来完成查询
	for query := 0; query < maxGoroutines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}
	wg.Wait()

	log.Println("Shutdown Program.")
	p.Close()
}
func performQueries(query int, p *pool.Pool) {
	//从池里请求一个连接
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	//将该连接释放回池里
	defer p.Release(conn)
	//用等待来模拟查询响应
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)

}
