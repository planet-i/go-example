package main

import "fmt"

func main() {
	var a, b string
	a, b = "Hello", "Go! "
	var c, d, e, f int = 2, 0, 2, 0
	var g, h = " from ", 2.13 //auto infer type
	i, j := "learning", 1
	var (
		k int
		l string
		m float32
		n bool
	)
	name := "Tom"
	fmt.Println("begin", a, b, c, d, e, f, g, h, i, j)
	fmt.Println(k, l, m, n)
	fmt.Println("Hello", name)
	fmt.Printf("%T \n", i)
	fmt.Printf("%T \n", j)
}
