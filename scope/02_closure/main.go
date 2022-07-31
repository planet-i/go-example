package main

import "fmt"

func main() {
	f := test(123)
	f()

	increment := wrapper()
	fmt.Println(increment()) //1
	fmt.Println(increment()) //2
	//increment := wrapper() 通过把这个函数变量赋值给increment，increment就成为了一个闭包。
	//i的生命周期没有随着它的作用域结束而结束
	fmt.Println(wrapper()()) //1
	fmt.Println(wrapper()()) //2
	//这里调用了三次wrapper()，返回了三个闭包，这三个闭包引用着三个不同的i，它们的状态是各自独立的。
}
func test(x int) func() {
	println(&x)
	return func() {
		x++
		println(&x, x)
	}
}
func wrapper() func() int {
	i := 0
	println(&i)
	return func() int {
		i++
		println(&i)
		return i
	}
}

//test函数或wrapper函数返回的是 所定义的匿名函数 和  所引用的环境变量的指针
