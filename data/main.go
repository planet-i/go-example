package main

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var a string
	println(a == "") //true

	b := "雨痕\x61\142\u0041" //字面量里允许十六进制、八进制、UTF格式
	fmt.Printf("%s\n", b)
	fmt.Println(b[1])
	//b[6] = 3
	fmt.Printf("% x,len: %d\n", b, len(b)) // len为9

	c := `line\r\n,      //这是字符串定义
    line2`
	println(c)

	d := "ab" + "cdf"
	println(&d)
	d += "cf"
	println(&d)
	println(d)
	println(d == "abcdfcf")
	println(d > "abcde") //

	s := "雨痕"
	println(s[1])
	println(&s)
	//println(&s[0])
	for i := 0; i < len(s); i++ { //byte 一个中文字符占3个字节
		fmt.Printf("%d: [%c]\n", i, s[i])
	}
	for i, c := range s { //rune 返回数组索引号，以及Unicode字符
		fmt.Printf("%d: [%c]\n", i, c)
	}

	s = "abcdefg" //字符串是可以变的呀？
	println(&s)
	s1 := s[:3]
	println(&s1)
	s2 := s[1:4]
	println(&s2)
	s3 := s[2:]
	println(&s3)
	println(s1, s2, s3)

	r := '我'
	i := string(r)
	k := byte(r)
	i1 := string(k)
	r2 := rune(k)
	fmt.Println(i, k, i1, r2, r)

	x, y := 10, 20
	a1 := [...]*int{&x, &y}
	p := &a1
	fmt.Printf("%T,%v\n", a1, a1)
	fmt.Printf("%T,%v\n", p, p)
	b1 := [...]int{1, 2}
	p1 := &b1
	p1[1] += 10                  //数组指针可以直接用来操作元素，而不是用先解成值
	println(&b1, &b1[0], &b1[1]) //任意获取元素地址， 字符串不可以获取元素地址
	println(p1[1], b1[0])

	sliceTest()
	stackTest()
	mapTest()
	interfaceTest()
}
func mapTest() {
	//nil字典
	var m1 map[string]int
	//m2、m3都是空字典
	m2 := map[string]int{}
	m3 := make(map[string]int)
	m4 := make(map[string]int, 0)
	//引用类型用make或者初始化表达语句
	m5 := make(map[string]int, 3)
	m6 := map[string]int{
		"chen": 21,
		"si":   25,
		"lin":  29,
	}
	fmt.Println(m1, len(m1), m1 == nil)
	fmt.Println(m2, len(m2), m2 == nil)
	fmt.Println(m3, len(m3), m3 == nil)
	fmt.Println(m4, len(m4), m4 == nil)
	fmt.Println(m5, len(m5), m5 == nil)
	fmt.Println(m6, len(m6), m6 == nil)
	fmt.Printf("m1: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m1)))
	fmt.Printf("m2: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m2)))
	fmt.Printf("m3: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m3)))
	fmt.Printf("m4: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m4)))
	fmt.Printf("m5: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m5)))
	fmt.Printf("m6: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&m6)))
	delete(m6, "chen")
	fmt.Println(m6, len(m6), m6 == nil)

}

//数组名字代表的是数组这个整体，不再代表这个数组的首地址。
func sliceTest() {
	s1 := make([]int, 3, 5)
	s2 := make([]int, 3)
	s3 := []int{10, 20, 5: 30}
	s4 := []int{} //内部指针被赋值，指向runtime.zerobase,依然完成了初始化操作
	s6 := []int{}
	s8 := make([]int, 0)
	s6 = append(s6, 1, 2, 3, 4)
	var s5 []int //定义了未初始化，默认赋值为nil，会分配所需内存
	var s7 []int
	s7 = append(s7, 1, 2, 3)
	fmt.Println(s1, len(s1), cap(s1))
	fmt.Println(s2, len(s2), cap(s2))
	fmt.Println(s3, len(s3), cap(s3))
	fmt.Println(s4, len(s4), cap(s4), s4 == nil)
	fmt.Printf("s4: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s4)))
	fmt.Printf("s6: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s6)))
	//&reflect.SliceHeader{Data:0x1191390, Len:0, Cap:0}   指向runtime.zerobase,不为nil
	fmt.Println(s5, len(s5), cap(s5), s5 == nil)
	fmt.Printf("s5: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s5)))
	fmt.Printf("s7: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s7)))
	fmt.Printf("s8: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s8)))
	//fmt.Printf("s9: %#v\n", (*reflect.SliceHeader)(unsafe.Pointer(&s9)))
	//&reflect.SliceHeader{Data:0x0, Len:0, Cap:0}    未初始化的切片对象， == nil

	s := []int{0, 1, 2, 3, 4}
	p := &s
	p0 := &s[0]
	p1 := &s[1]
	println(p, p0, p1)

}

//在切片的基础上创建切片，实现栈式数据结构
func stackTest() {
	stack := make([]int, 0, 6)
	push := func(x int) error {
		n := len(stack)
		if n == cap(stack) {
			return errors.New("stack is full")
		}
		stack = stack[:n+1] //先扩栈再放值
		stack[n] = x
		return nil
	}
	pop := func() (int, error) {
		n := len(stack)
		if n == 0 {
			return 0, errors.New("stack is empty")
		}
		x := stack[n-1]
		stack = stack[:n-1]
		return x, nil
	}

	for i := 0; i < 7; i++ {
		fmt.Printf("push %d: %v,%v\n", i, push(i), stack)
	}

	for i := 0; i < 7; i++ {
		x, err := pop()
		fmt.Printf("pop %d: %v,%v\n", x, err, stack)
	}
}
func interfaceTest() {
	var a interface{}
	var b interface{} = nil
	var c interface{} = (*int)(nil)
	println(a == nil, b == nil, c == nil, a == b)
	println(a, b, c)
}
