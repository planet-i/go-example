package main

func count() int {
	print("count.")
	return 3
}
func main() {
	for i, c := 0, count(); i < c; i++ { //初始化语句的count()只执行一次
		println("i", i)
	}

	c := 0
	for c < count() { //条件表达式中的count()重复执行
		println("b", c)
		c++
	}
	println()
	data := [3]string{"a", "b", "c"}
	for i, s := range data {
		println(i, s)
	}
}
