package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	//b := []int{-1, -2, -3}
	var b = make([]int, 3)
	fmt.Printf("%p and %p \n", a, b)
	//copy(b, a)
	b = a
	a[0] = 100
	fmt.Println(a, b)
	fmt.Printf("%p and %p \n", a, b)

	test()
}

func test() {
	// 复制指针类型：b跟着a变
	// 复制值类型：b的值类型变与a无关，b的指针类型变化与a相关
	a := T{
		v: 0,
		p: new(int),
	}
	b := a
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d\n", a.v, *a.p, b.v, *b.p)
	// b.p = new(int)
	a.Set(1) // 等价于(&a).Set(1)
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d\n", a.v, *a.p, b.v, *b.p)
}

type T struct {
	v int
	p *int
}

func (t *T) Set(v int) {
	t.v = v
	*t.p = v
}
