package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	testnil("aaa")
	// err := handlePanic()
	// fmt.Println("main 函数得到错误:", err)
	// panic被捕获之后,可以继续往下执行
	// testError()
	// afterErrorfunc()
	// Compile()
}

func handlePanic() (err error) {
	//go func() { // recover只能捕获同一个goroutine中的panic
	//defer recover() // defer不可直接调用recover函数
	defer func() {
		// defer func() { // 不可用嵌套的defer
		if r := recover(); r != nil { // 要避免以nil为参数抛出panic，这样recover不能顺利捕获
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("Unknown panic:%v", r)
			}
		}
		//}()
	}()
	// defer func() {
	// 	recover捕获的是它最近执行的那一个panic
	// 	defer中又遇见了panic的话，则会释放这个defer，去执行下一个defer。
	// 	panic("第二次panic")
	// }()
	// }()
	panic("第一次panic")
}

func testError() {
	defer catch()
	panic(" \"panic 错误\"")
	fmt.Println("抛出一个错误后继续执行代码")
}

func catch() {
	var err error
	// 错误被recover 函数接收，转化为error类型的错误
	if r := recover(); r != nil {
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = fmt.Errorf("Unknown panic:%v", r)
		}
	}
	if err != nil {
		fmt.Println("recover后的错误:", err)
	}
}

func afterErrorfunc() {
	fmt.Println("遇到错误之后 func ")
}

// effective go示例
type Error string

func (e Error) Error() string {
	return string(e)
}

type Regexp struct{}

func (regexp *Regexp) error(err string) {
	panic(Error(err))
}
func (regexp *Regexp) doParse() *Regexp {
	panic("出现panic")
	return &Regexp{}
}

func new() *Regexp {
	return &Regexp{}
}

// Compile returns a parsed representation of the regular expression.
func Compile() (regexp *Regexp, err error) {
	regexp = new()
	// doParse will panic if there is a parse error.
	defer func() {
		if e := recover(); e != nil {
			regexp = nil    // Clear return value.
			err = e.(Error) // Will re-panic if not a parse error.
		}
	}()
	return regexp.doParse(), nil
}

var (
	m sync.Map
)

type a struct {
	x      int
	y      string
	status chan string
}

func testnil(s string) {
	v, ok := m.Load(1)
	var down *a
	if !ok {
		aa := &a{
			x: 22,
			y: "sl",
		}
		m.Store(aa, 1)
		fmt.Println("存进去了")
	} else {
		if down, ok = v.(*a); !ok {
			fmt.Println("拿到了")
		}
	}
	down.status <- s
}
