package main

func main() {
	a := make(chan int)
	b := make(chan int, 3) //创建带三个缓冲槽的异步通道
	done := make(chan struct{})

	b <- 1
	b <- 2
	println(<-b)
	println(<-b) //异步模式在缓冲区未满或者数据未读完之前，都不会阻塞
	println("a:", len(a), cap(a))
	println("b:", len(b), cap(b))

	go func() {
		defer close(done)
		for {
			x, ok := <-a
			if !ok {
				println("error")
				return
			}
			println(x)
		}
	}()

	a <- 100
	a <- 200
	a <- 300
	close(a) //同步模式必须有配对操作的goroutine出现，否则会一直阻塞
	<-done
}
