package main

import "fmt"

func main() {
	var a [4]int
	b := [4]int{2, 5}      //未提供初始值的元素自动初始化为0
	c := [4]int{2, 3: 10}  //指定索引位置初始化
	d := [...]int{1, 5: 9} //编译器通过初始值数量推断数组长度
	e := [...][2]int{
		{1, 2},
		{3, 4},
	}

	fmt.Println(len(a), cap(a))
	fmt.Println(b == c)
	b[1] = 8
	fmt.Println(b[1])
	fmt.Println(a, b, c, d, e)

	type user struct {
		name string
		age  int
	}
	f := [...]user{
		{"Tom", 8},
		{"Jerry", 7},
	}
	fmt.Printf("%#v \n", f)
}
