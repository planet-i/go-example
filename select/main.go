package main

import "fmt"

func main() {
	var c1, c2, c3 chan int
	var i1, i2 int
	/* 	c1 = make(chan int)
	   	go func() {
	   		c1 <- 100
	   	}() */
	select {
	case i1 = <-c1:
		fmt.Print("received ", i1, " from c1\n")
	case c2 <- i2:
		fmt.Print("sent ", i2, " to c2\n")
	case i3, ok := <-c3:
		if ok {
			fmt.Print("received ", i3, " from c3\n")
		} else {
			fmt.Printf("c3 is closed\n")
		}
	default:
		fmt.Printf("no communication\n")
	}
}

/* select 随机执行一个可运行的 case；
如果没有 case 可运行，select将阻塞，直到有 case 可运行,default默认执行 ；
每个 case 必须是一个通信操作，要么是发送要么是接收。
每个 case 必须是一个通信操作，要么是发送要么是接收。 */
