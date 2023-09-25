package main

import (
	"fmt"
)

func main() {
	testVar()   // var 的使用
	testShort() // := 的使用
	testFunc()  // 一个变量作用域的坑
}

func testVar() {
	fmt.Println("-------testVar")
	// 显示声明、自动初始化为类型零值
	var a int
	// 显示定义
	var b int = 100
	// 隐式定义
	var c = false // 自动推断类型 c 为bool型

	// 并行声明
	var x, y int
	// 并行定义
	var z, h float32 = -1, -2

	// 块定义
	var ( //组的方式一次定义多个
		x1, y1     string
		a1, b1, c1 = 1.0, 100, "abc"
	)
	// 多变量赋值
	var (
		a2     = 1              // 1
		b2, c2 = a2 + 1, a2 + 2 // 2 3   先计算好右边的
	)

	fmt.Printf("a:%+v；b:%+v；c:%+v；\n", a, b, c)
	fmt.Printf("x:%+v；y:%+v；z:%+v；h:%+v；\n", x, y, z, h)
	fmt.Printf("x1:%+v；y1:%+v；a1:%+v；b1:%+v；c1:%+v；\n", x1, y1, a1, b1, c1)
	fmt.Printf("a2:%+v；b2:%+v；c2:%+v；\n", a2, b2, c2)
}

// := 只能用于局部变量、不能用于全局变量
func testShort() {
	fmt.Println("-------testShort")
	a := 100 // 隐式定义、自动推导类型

	// 并行定义
	x, y := 0, 10
	_, z, _ := 1, 2, 3 // 不准全为_,  左侧必须有新变量

	fmt.Printf("a:%+v；\n", a)
	fmt.Printf("x:%+v；y:%+v；z:%+v；\n", x, y, z)
}

func testFunc() {
	fmt.Println("-------testFunc")
	fmt.Printf("befor: %p, %T\n", p, p)
	p, err := foo() // 因为作用域不同，此局部变量p遮盖了全局变量p，全局变量p依旧为nil
	fmt.Printf("after: %p, %T\n", p, p)
	if err != nil {
		return
	}
	bar()
	fmt.Println("end", *p)
}

var p *int

func foo() (*int, error) {
	var i int = 5
	return &i, nil
}
func bar() {
	fmt.Printf("bar: %p, %T\n", p, p)
	fmt.Println(*p) // 使用的是全局变量p，nil造成panic
}
