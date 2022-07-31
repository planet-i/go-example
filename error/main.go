package main

import (
	"errors"
	"fmt"
)

func main() {
	//NewError()
	err := handlePanic()
	fmt.Println(err)
}

func handlePanic() (err error) {
	//go func() {
	defer func() {
		//defer func() {
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
		//}()
	}()
	//}()
	panic("异常信息")
}

func Arecover() (err error) {
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
	return err
}

func New(text string) error {
	return errorsString{text}
}

// errorString is a trivial implementation of error.
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
