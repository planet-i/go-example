package main

import "fmt"

func main() {
	x, y := 1, 2
	defer func(a int) {
		println("defer:x,y = ", a, y) //y为闭包引用
	}(x) //注册时复制调用函数
	x += 100
	y += 200
	println(x, y)
	//----延迟函数的运行时机
	DeferAndReturn()
	DeferAndReEnd()
	DeferAndPanic()
}

func DeferAndReturn() {
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
	return //就近原则，多个defer时，离return越近的先执行
	//fmt.Println(4)
}
func DeferAndReEnd() {
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
}
func DeferAndPanic() {
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
	panic("panic")
	//fmt.Println(4)
}
