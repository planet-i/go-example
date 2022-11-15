package main

import "fmt"

func main() {
	//testExecuteTimes()
	testForRangeType()
	//testRangeKey()
}

// 测试forRange边读边写三种类型
func testForRangeType() {
	// count := 0
	// s := []int{1, 2, 3}
	// for _, v := range s {
	// 	fmt.Println("range slice", count)
	// 	s = append(s, v)
	// }
	// fmt.Println(s)

	// count = 0
	// c := make(chan int, 1)
	// c <- 0
	// for i := range c {
	// 	fmt.Println("range channel", count)
	// 	c <- i
	// }

	// count = 0
	// m := map[int]int{
	// 	1: 1,
	// 	2: 2,
	// }
	// for k, v := range m {
	// 	fmt.Println("range map", count)
	// 	k += 100
	// 	v += 100
	// 	m[k] = v
	// }
	// fmt.Printf("%v\n", m)

	// data := [3]string{"a", "b", "c"}
	// for i, s := range data {
	// 	println(i, s)
	// 	// 数组的长度不可变
	// }

	st := "中国abc"
	fmt.Println("字符串长度", len(st))
	for i, v := range st {
		println(i, v)
	}
	for i := 0; i < len(st); i++ {
		println(st[i])
	}
}

// 测试for控制结构中，条件的执行次数
func testExecuteTimes() {
	for i, c := 0, count(); i < c; i++ { //初始化语句的count()只执行一次
		println("i", i)
	}

	c := 0
	for c < count() { //条件表达式中的count()重复执行
		println("b", c)
		c++
	}
	println()
}
func count() int {
	print("count.")
	return 3
}

func testRangeKey() {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	oldArr := []*int{}
	// range的v地址一直不变，每遍历一个值，把值赋值给v
	for i, v := range arr {
		oldArr = append(newArr, &arr[i])
		newArr = append(newArr, &v)
	}
	for range arr {
		fmt.Println(arr)
	}
	for _, v := range newArr {
		fmt.Println(*v) // 3 3 3
	}
	for _, v := range oldArr {
		fmt.Println(*v) // 1 2 3
	}
}
