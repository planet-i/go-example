package main

import (
	"fmt"
	"time"
)

var timeout <-chan time.Time  //声明了一个Time类型的单向通道
var result chan int
func main(){
	timeout = time.After(time.Second *6)//time.After创建了一个定时器，这个定时器会在6秒后向通道timeout写入当前的时间，
	result = make(chan int)

	//开启一个goroutine
	go func(){
		fmt.Println("--begin do task--")
		time.Sleep(time.Second*3)
		fmt.Println("--end do task")
		result <- 100
	}()

	select {
	case e := <-result:
		fmt.Println("get result:",e)
	case <-timeout:
		fmt.Println("get result timeout")//判断任务状态：6s内任务完成则输出get result;未完成，输出超时信息
	}
}