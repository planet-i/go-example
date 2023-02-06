package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var (
	name   string
	age    int
	weight *int
	height *int
	number *int
)

func init() {
	// 支持长选项、短选项、默认值和使用文本，并将标志的值绑定到变量。
	pflag.StringVarP(&name, "name", "n", "null", "Input Your Name")
	// 支持长选项、默认值和使用文本，并将标志的值绑定到变量。
	pflag.IntVar(&age, "age", 0, "Input Your Age")
	// 支持长选项、默认值和使用文本，并将标志的值存储在指针中
	weight = pflag.Int("weight", 0, "Input Your Weight")
	// 支持长选项、短选项、默认值和使用文本，并将标志的值存储在指针中
	height = pflag.IntP("height", "h", 0, "Input Your Weight")
	// 定义非选项参数
	number = pflag.Int("", 1234, "help message for number")
}

func main() {

	pflag.Parse()

	// 如果有必须填的参数没有填，打印帮助文档
	if age == 0 || *weight == 0 {
		fmt.Println("here")
		pflag.PrintDefaults()
		return
	}

	// go run main.go -n xx --weight 10 -h 20 --age 99 短选项用-，长选项用--
	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(*weight)
	fmt.Println(*height)
	// go run main.go xx yy zz -n xx --weight 10 -h 20 --age 99
	fmt.Printf("argument number is: %v\n", pflag.NArg())
	fmt.Printf("argument list is: %v\n", pflag.Args())
	fmt.Printf("the first argument is: %v\n", pflag.Arg(0))
	fmt.Println(*number)
}
