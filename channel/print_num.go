package main

// N个协程合作输出1到M的数字

import (
	"fmt"
	"sync"
)

func WPrintNum() {
	var wg sync.WaitGroup
	wg.Add(1) // 仅需等待1个协程

	go func() {
		defer wg.Done() // 协程结束时标记完成
		for i := 0; i <= 100; i++ {
			fmt.Println(i) // 子协程独立执行循环，主协程仅负责启动和同步
		}
	}()

	wg.Wait() // 主协程通过sync.WaitGroup等待子协程完成，避免主程序提前退出
}

func CPrintNum() {
	ch := make(chan struct{}) // 信号通道

	go func() {
		for i := 0; i <= 100; i++ {
			fmt.Println(i)
		}
		ch <- struct{}{} // 发送结束信号
	}()

	<-ch // 通过无缓冲通道ch实现主协程等待子协程完成
}

// 通过N个channel形成环形触发链路，每个协程接收信号后打印对应数值，并通知下一个协程。此方案完全依赖通道通信，无需锁机制，逻辑清晰且线程安全。
func GPrintNum() {
	const N = 3   // 协程数量
	const M = 100 // 总打印数

	channels := make([]chan int, N)
	for i := 0; i < N; i++ {
		channels[i] = make(chan int)
	}

	var wg sync.WaitGroup
	wg.Add(N)

	// 创建N个协程，形成环形链路
	for i := 0; i < N; i++ {
		nextIndex := (i + 1) % N
		// 0 -> 1 -> 2 -> 0
		go printNum(i, channels[i], channels[nextIndex], &wg, M)
	}

	channels[0] <- 1 // 主协程向第一个channel发送初始值1，触发协程链式反应。
	wg.Wait()
}

// printNum 将当前通道接收的数据加一然后输出到下一个通道
func printNum(id int, chCurrent, chNext chan int, wg *sync.WaitGroup, max int) {
	defer wg.Done()
	for {
		num, ok := <-chCurrent
		if !ok || num > max { // 终止条件：通道关闭或超过最大值
			close(chNext)
			fmt.Printf("协程%d退出,原因-> 通道关闭:%v,超过最大值:%v\n", id+1, !ok, num > max)
			return // 收到M+1数据的g因为数值超过M退出。再触发协程链式反应，使后面的协程接连退出
		}
		fmt.Printf("协程%d打印: %d\n", id+1, num)
		chNext <- num + 1 // 传递下一个数字给后续协程
	}
}
