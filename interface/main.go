package main

import "fmt"

// type stringer interface{
// 	string() string
// }
// type tester interface{
// 	test()
// 	stringer
// }
// type data struct{}
// func (*data) test(){}
// func(data) string()string{
// 	return "???"
// }
// func main(){
// 	var d data
// 	var t tester = &d
// 	t.test()
// 	fmt.Println(t.string())
// }

type aaa interface {
	bbb()
}

// type b struct{}

func bbb() {
	fmt.Println("bbbb")
}
func main() {
	a := aa()
	fmt.Println(a == nil)
	a.bbb()
}

func aa() aaa {
	return nil
}
