package main

import (
	"fmt"
	"time"
)

var c int

func counter() int {
	c++
	return c
}
func testFuncParam() {
	a := 100
	go func(x, y int) {
		time.Sleep(time.Second) // 让goroutine在main逻辑之后执行
		println("go:", x, y)
	}(a, counter()) // 立即计算并复制参数 100,1

	a += 100
	println("main:", a, counter()) // 200,2
	time.Sleep(time.Second * 3)    // 等待goroutine结束
	go println("hello!world")
	go func(s string) {
		println(s)
	}("hello,go")
	time.Sleep(5 * time.Second)
}

func main() {
	// testNoParam(5)
	// testFuncParam()
	// testPrintNum()
	testPool()
	// 阻塞主goroutine，保持程序运行
	//select {}
}

func testNoParam(n int) {
	// 当goroutine被调度执行时，for循环早已完成，此时i的值已经变成了5。但是如果for循环很大，当g被调度的时候，for循环未完成，输出的就是i的当前值
	for i := 0; i < n; i++ {
		go func() {
			fmt.Println(i) // 始终输出i的当前值
		}()
	}
	time.Sleep(5 * time.Second)
}

func testPrintNum() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})
	go printNum(ch1, ch2, 1)
	go printNum(ch2, ch3, 2)
	go printNum(ch3, ch1, 3)

	ch1 <- struct{}{}

	time.Sleep(1 * time.Second)
}

func printNum(in, out chan struct{}, num int) {
	for {
		<-in
		fmt.Println(num)
		out <- struct{}{}
	}
}

// 实现一个协程池
type Pool struct {
	work chan func()   // 接收task任务
	sem  chan struct{} // 设置协程池大小，即可同时执行的协程数量
}

// New 创建一个协程池对象 size 为协程池大小
func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

// worker 执行任务
func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()
	for {
		task()
		task = <-p.work
	}
}

// NewTask 添加任务
func (p *Pool) NewTask(task func()) {
	select {
	case p.work <- task: // work无缓冲，sem有缓冲，第一次添加任务的时候会先走第二个case
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func testPool() {
	pool := New(10)
	pool.NewTask(func() {
		fmt.Println("run task")
	})
	time.Sleep(10 * time.Second)
}
