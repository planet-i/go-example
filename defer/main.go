package main

import (
	"fmt"
	"os"
)

func main() {
	//----延迟函数的运行时机
	deferExit()
	testPanic()
	// return的返回值是有名返回值或者是指针的时候，defer才会对return的值造成影响
	testReturn()
	fmt.Println("return", testReturnNoName())
	fmt.Println("return", testReturnHaveName())
	fmt.Println("return", *testReturnPointer())
	// defer后的函数带参数的时候，会临时先将参数值记录下来
	testParam()
	// defer后的函数带的参数也是函数的话，也会先执行好
	testCalc()
}

// ---- 测试defer和panic的执行顺序
func testPanic() {
	fmt.Println("-------testPanic")
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
	panic("panic") // 碰上panic，先执行defer再来执行panic   3 2 1 panic
}

// ---- 测试defer和return的顺序
func testReturn() (i int) {
	fmt.Println("-------testReturn")
	i = 1
	defer func() {
		fmt.Println("defer", i)
		i = 2
	}()

	return 3
} // i=1 i=3 defer 3  i=2 return 2

// ---- 测试defer顺序与是否命名返回值的关系
func testReturnNoName() int {
	fmt.Println("-------testReturnNoName")
	i := 0
	defer func() {
		fmt.Println("defer先", i)
	}()
	defer func() {
		i += 1
		fmt.Println("defer后", i)
	}()
	return i // 执行 return 语句后，Go 会创建一个临时变量保存返回值，局部变量i的变化与返回值无关
} // defer后 1  defer先 1  return 0

func testReturnHaveName() (i int) {
	fmt.Println("-------testReturnHaveName")
	i = 0
	defer func() {
		i += 1
		fmt.Println("defer", i)
	}()
	return i // 有名返回值的函数，执行 return 语句时，不会再创建临时变量保存，defer修改了i，会对返回值造成影响
} // defer 1 return 1

// testReturnPointer 返回无名指针，defer也会对返回值造成影响
func testReturnPointer() *int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2:", i)
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	return &i
}

// ---- 测试defer函数带参数
func testParam() {
	fmt.Println("-------testParam")
	x, y := 1, 2
	defer func(x int) {
		fmt.Println("defer先: x,y = ", x, y) // y为闭包引用
	}(x) // 注册时复制调用函数
	defer func() {
		fmt.Println("defer后: x,y = ", x, y) // x,y为闭包引用
	}()
	x += 100
	y += 200
}

// 测试defer函数调用中参数执行情况
func testCalc() {
	fmt.Println("-------testCalc")
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b)) // 10 1 2 3      1 1 3 4
	a = 0                                //    |             |
	defer calc("2", a, calc("20", a, b)) // 20 0 2 2  ->  2 0 2 2
	b = 1
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

// deferExit 在os.Exit函数之前的defer不会执行
func deferExit() {
	defer func() {
		fmt.Println("defer")
	}()
	os.Exit(0)
}
