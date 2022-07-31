package main

import (
	"log"
	"sync"
	"time"

	"github.com/planet-i/goexample1/work"
)

var names = []string{
	"steve",
	"bob",
	"mary",
	"therese",
	"jason",
}

//实现了Worker接口的工作
type namePrinter struct {
	name string
}

func (m *namePrinter) Task() {
	log.Println(m.name)
	time.Sleep(time.Second)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100 * len(names))
	//使用两个goroutine来创建工作池
	p := work.New(2)
	for i := 0; i < 100; i++ {
		for _, name := range names { //声明并创建100*5个goroutine
			//当场实例化各个work
			np := namePrinter{
				name: name,
			}
			go func() {
				p.Run(&np)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	p.Shutdown()
}
