package main
import "fmt"
type stringer interface{
	string() string
}
type tester interface{
	test()
	stringer
}
type data struct{}
func (*data) test(){}
func(data) string()string{
	return "???"
}
func main(){
	var d data
	var t tester = &d
	t.test()
	fmt.Println(t.string())
}