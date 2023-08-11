package main

import "fmt"

func main() {
	// map的三种声明和初始化方式
	var m1 map[string]int
	fmt.Println(m1, m1 == nil)
	m2 := make(map[string]int)
	fmt.Println(m2, m2 == nil)
	m3 := map[string]int{}
	fmt.Println(m3, m3 == nil)

	m := map[int]string{
		0: "abc",
		1: "def",
	}
	m[1] = "tyu"
	m[2] = "bnm"
	for i, s := range m {
		println("key: %d ==> value: %s", i, s)
	}
	fmt.Println("map的长度", len(m)) // len返回键值对数量，cap不支持map类型
	delete(m, 0)                  // 删除键值对，当其不存在时不会出错
	fmt.Println(m)
}
