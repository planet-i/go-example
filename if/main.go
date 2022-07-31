package main

import "fmt"

func main() {
	a := true
	//-----可定义块局部变量或执行初始化函数，局部变量的有效范围包含整个if/else块
	if a, b, c := 1, 2, 3; a+b+c > 6 {
		fmt.Println("大于六")
	} else {
		fmt.Println("小于等于六")
		fmt.Println(a)
	}
	fmt.Println(a)
}
