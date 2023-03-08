package main

import "fmt"

func main() {
	testPointer()
}

type T struct {
	v int
	p *int
}

func (t T) Set(value int) {
	t.v = value
	*t.p = value
	fmt.Printf("T方法内 t.v=%d, t.p=%d;\n", t.v, *t.p)
}

func (t *T) Set2(value int) {
	t.v = value
	*t.p = value
	fmt.Printf("*T方法内 t.v=%d, t.p=%d;\n", t.v, *t.p)
}

func testPointer() {
	a := T{
		v: 0,
		p: new(int),
	}
	// --- 变量被复制的是值类型还是指针类型的区别
	b := a  // 复制值类型：b的值类型变与a无关，b的指针类型变化与a相关
	c := &a // 复制指针类型：c跟着a变
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d;  c.v=%d, c.p=%d  [原始的]\n", a.v, *a.p, b.v, *b.p, c.v, *c.p)
	a.v = 888
	*b.p = 777
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d;  c.v=%d, c.p=%d  [ab变化后]\n", a.v, *a.p, b.v, *b.p, c.v, *c.p)

	// --- 调用的方法接收者是T还是*T的区别
	a.Set(1) // 等价于a.Set(1)
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d  [调用T方法]\n", a.v, *a.p, b.v, *b.p)
	a.Set2(2) // (&a).Set2(2) 复制了一份a去调用的
	fmt.Printf("a.v=%d, a.p=%d;  b.v=%d, b.p=%d  [调用*T方法]\n", a.v, *a.p, b.v, *b.p)
}
