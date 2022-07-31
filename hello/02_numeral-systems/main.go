package main

import "fmt"

func main() {
	fmt.Println(42)                                        //decimal
	fmt.Printf("%d \t %b \n", 42, 42)                      //binary
	fmt.Printf("%d \t %b \t %#o\n", 42, 42, 42)            //octonary
	fmt.Printf("%d \t %b \t %#o \t %#X\n", 42, 42, 42, 42) //hexadecimal
	//	fmt.Printf("%X \n", 42)
	//	fmt.Printf("%x \n", 42)
	//	fmt.Printf("%#x \n", 42)
}

/*
import导入标准库或第三方包；
未使用的导入会被认定为错误；
*/
