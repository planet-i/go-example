package main

import (
	"errors"
	"fmt"
)

func main() {
	NewError()
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
