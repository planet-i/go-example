package main

import "fmt"

func main() {
	m := make(map[string]int)
	fmt.Println(m)

	m2 := map[int]string{
		0: "abc",
		1: "def",
	}
	m2[1] = "tyu"
	m2[2] = "bnm"
	for i, s := range m2 {
		println(i, s)
	}
	fmt.Println(len(m2)) //len返回键值对数量，cap不支持map类型
	delete(m2, 0)        //删除键值对，当其不存在时不会出错
	fmt.Println(m2)
}
