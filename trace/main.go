package main

import (
	"runtime/debug"
)

func test1() {
	test2()
}

func test2() {
	test3()
}

func test3() { // 可以通过 debug.PrintStack() 直接打印，也可以通过 debug.Stack() 方法获取堆栈然后自己打印fmt.Printf("%s", debug.Stack())
	debug.PrintStack()
}

func main() {
	test1()
}
