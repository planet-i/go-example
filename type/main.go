package main

import "fmt"

type SourceType int
type NewTypeName SourceType

type (
	// 新定义类型a    源类型A
	NewTypeName1 SourceType
	NewTypeName2 SourceType
	// 新定义类型b    源类型A
)

// func main() {

// 	var a NewTypeName1 = 1
// 	var b NewTypeName2 = 2
// 	var c SourceType = 3
// 	a = NewTypeName1(b)
// 	a = NewTypeName1(c)
// 	fmt.Println(a, b, c)
// }

type ggg struct {
	g *string
}

func main() {
	ss := ff()
	fmt.Println(ss, ss.g)
}
func ff() (p *ggg) {
	return
}
