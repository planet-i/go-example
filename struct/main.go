package main

import "fmt"

type node struct {
	_    bool //忽略值输出的是类型默认值
	id   int
	next *node
}
type aa struct {
	a int
	b string
}
type bb struct {
	aa
	c int
}
type cc struct {
	a int
	b string
	c int
}

func testContain() {
	x := aa{
		a: 1,
		b: "ioio",
	}
	y := bb{
		c: 6,
	}
	z := cc{
		a: x.a,
		b: x.b,
		c: y.c,
	}
	y.aa = x
	fmt.Println(y)
	fmt.Println(z)
}
func main() {
	//testContain()
	//n1 := node{
	//	id : 100,
	//}
	//n2 := node{
	//	id : 200,
	//	next : &n1,
	//}
	//n3 := node{true,300,&n2} //按顺序初始化全部字段，不能遗漏
	//fmt.Println(n1,n2,n3)
	s := structParam()
	fmt.Println(s.a, s.b)
}

func structParam() (avar aa) {
	return avar
}
