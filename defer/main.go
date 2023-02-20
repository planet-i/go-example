package main

import "fmt"

func main() {
	//----延迟函数的运行时机
	DeferAndPanic()
	testParam()
	fmt.Println("return", test())
	fmt.Println("return", testName())
	fmt.Println("return", a())
}

// ---- 测试defer和panic的执行顺序
func DeferAndPanic() {
	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)
	panic("panic") // 碰上panic，先执行defer再来执行panic
}

// ---- 测试defer的参数
func testParam() {
	x, y := 1, 2
	defer func(a int) {
		fmt.Println("defer: x,y = ", a, y) // y为闭包引用
	}(x) // 注册时复制调用函数
	x += 100
	y += 200
	fmt.Println("当前: x,y = ", x, y)
}

// ---- 测试defer和return的顺序
func a() (i int) {
	i = 1
	defer func() {
		fmt.Println("defer", i)
		i = 2
	}()

	return 3
} // i=1 i=3 defer 3  i=2 return 2

// ---- 测试defer顺序与是否命名返回值的关系
func test() int {
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

func testName() (i int) {
	i = 0
	defer func() {
		i += 1
		fmt.Println("defer", i)
	}()
	return i // 有名返回值的函数，执行 return 语句时，不会再创建临时变量保存，defer修改了i，会对返回值造成影响
} // defer 1 return 1
