package main

import (
	"fmt"

	"github.com/pkg/errors"
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

var ErrCreat = errors.New("错误")

func TestErrors() {
	num := 1
	fmt.Printf("%+v", caller(num))
}
func caller(num int) error {
	switch num {
	case 1:
		return a()
	case 2:
		return b()
	case 3:
		return c()
	}
	return nil
}

func a() error {
	return b()
}
func b() error {
	return c()
}

func c() error {
	return ErrCreat
	//return errors.WithStack(ErrCreat)
}
