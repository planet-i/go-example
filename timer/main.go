package main

import (
	"fmt"
	"time"
)

// func main() {
// 	fmt.Println("Starting")
// 	timer := time.NewTimer(2 * time.Second)
// 	<-timer.C
// 	fmt.Println("定时器触发")
// }

// var timeout <-chan time.Time  //声明了一个Time类型的单向通道
// var result chan int
// func main(){
// 	timeout = time.After(time.Second *6)//time.After创建了一个定时器，这个定时器会在6秒后向通道timeout写入当前的时间，
// 	result = make(chan int)

// 	//开启一个goroutine
// 	go func(){
// 		fmt.Println("--begin do task--")
// 		time.Sleep(time.Second*3)
// 		fmt.Println("--end do task")
// 		result <- 100
// 	}()

// 	select {
// 	case e := <-result:
// 		fmt.Println("get result:",e)
// 	case <-timeout:
// 		fmt.Println("get result timeout")//判断任务状态：6s内任务完成则输出get result;未完成，输出超时信息
// 	}
// }

func main() {
	fmt.Println("Starting", time.Now())
	syncSleep()
	fmt.Println("End", time.Now())
}

// func syncSleep() {
// 	ticker := time.Tick(5 * time.Second) // 每隔1分钟执行一次
// 	var sleepTime = 60
// 	var newTick = time.NewTicker(time.Duration(sleepTime) * time.Second)
// 	for {
// 		fmt.Println("循环啦")
// 		select {
// 		case <-ticker: // 每隔1分钟执行一次
// 			fmt.Println("5秒钟到啦", time.Now())
// 			newSleepTime := 60
// 			if newTick != nil && newSleepTime != sleepTime {
// 				newTick.Stop()           // 停止之前的定时器
// 				sleepTime = newSleepTime // 更新 sleepTime
// 			}
// 			// 重新创建定时器
// 			newTick = time.NewTicker(time.Duration(sleepTime) * time.Second)
// 			// 等待定时器触发或者超时
// 		case <-time.After(time.Duration(sleepTime) * time.Second): // 阻塞 sleepTime 时间后退出
// 			fmt.Println("阻塞时间到了", sleepTime, time.Now())
// 			return
// 		case <-newTick.C: // 定时器触发
// 			fmt.Println("定时器触发", time.Now())
// 			return
// 		}
// 	}
// }

func syncSleep() {
	sleepTime := 60 // 初始的 sleepTime
	ticker := time.NewTicker(time.Duration(sleepTime) * time.Second)
	for {
		select {
		case <-time.After(5 * time.Second): // 每隔5秒检查一次新的 sleepTime
			newSleepTime := getSyncSleepTime() // 获取新的 sleepTime 值
			fmt.Println("新的sleepTime", newSleepTime, time.Now())
			if newSleepTime != sleepTime {
				if ticker != nil {
					ticker.Stop() // 停止之前的定时器
				}
				sleepTime = newSleepTime // 更新 sleepTime
				// 创建新的定时器
				ticker = time.NewTicker(time.Duration(sleepTime) * time.Second)
			}
		case <-ticker.C: // 定时器触发
			fmt.Println("定时器触发", time.Now())
			return
		}
	}
}

func getSyncSleepTime() int {
	// rand.Seed(time.Now().UnixNano()) // 设置随机种子
	// return rand.Intn(100)            // 返回 0 到 99 之间的随机整数
	return 60
}
