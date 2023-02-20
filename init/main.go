package main

import "fmt"

type ErrorImpl struct{}

func (e *ErrorImpl) Error() string {
	return ""
}

var ei *ErrorImpl

func ErrorImplFun() error {
	return ei
}

func main() {
	f := ErrorImplFun()
	fmt.Println(f == nil)
}
