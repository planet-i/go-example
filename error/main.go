package main

import (
	"errors"
	"fmt"

	errors1 "github.com/pkg/errors"
)

func main() {
	// NewError()
	TestErrors()
}

// err的创建
func New(text string) error {
	return &errorsString{text}
}

type errorsString struct {
	s string
}

func (e errorsString) Error() string {
	return e.s
}

func NewError() {
	e1 := New("生成一个错误")
	e2 := New("生成一个错误")
	e3 := errors.New("错误")
	e4 := errors.New("错误")
	fmt.Println(e1 == e2)
	fmt.Println(e3 == e4)
}

// errors第三方库的使用
func TestErrors() {
	num := 1
	// 输出错误信息，不包含堆栈
	fmt.Printf("%s\n", a(num))
	fmt.Printf("%T %v\n", errors1.Cause(a(num)), errors1.Cause(a(num)))
	// 输出带引号的错误信息，不包含堆栈
	fmt.Printf("%q\n", a(num))
	// 输出错误信息和堆栈
	fmt.Printf("%T %+v\n", a(num), a(num))
}
func a(num int) error {
	switch num {
	case 1:
		return b()
	case 2:
		return c()
	case 3:
		return d()
	}
	return nil
}

func b() error {
	// return errors1.New("new error")
	aa := "hh"
	return errors1.Errorf("new error %s", aa)
}

func c() error {
	err := errors.New("aaa")
	return errors1.Wrap(err, "warp error")
}
func d() error {
	err := errors.New("aaa")
	return err
}
