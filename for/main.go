package main

import "fmt"

func main() {
	testForRangeType()
	testRangeString()
	testExecuteTimes()
	testRangeKey()
	testForType()
}

// testForRangeType (测试for Range边读边写三种类型)
func testForRangeType() {
	// slice类型，边读边写只会遍历 len(原始slice) 次
	var count = 0
	s := []string{"aa", "bb", "cc"}
	for _, v := range s {
		count++
		fmt.Println("range slice", count) // range slice 0\1\2
		s = append(s, v)
	}
	fmt.Println("最终的slice: ", s)

	// chan的边读边写，一直运行
	count = 0
	c := make(chan int, 1)
	c <- 0
	for i := range c {
		count++
		fmt.Println("range channel", count)
		c <- i // 一直运行
	}

	// map的边读边写，运行次数不定
	count = 0
	m := map[int]int{
		1: 1,
		2: 2,
	}
	for k, v := range m {
		count++
		fmt.Println("range map", count)
		k += 100
		v += 100
		m[k] = v
	}
	fmt.Printf("最终的map: %v\n", m)
}

// testRangeString (字符串的遍历时，使用 for range 循环会以 rune（Unicode 码点）为单位遍历字符串，确保正确处理多字节字符；而使用简单的 for 循环则会以字符串的字节数为单位遍历字符串，可能导致在处理多字节字符时出现问题。)
func testRangeString() {
	fmt.Println("// 字符串的遍历时，使用 for range 循环会以 rune（Unicode 码点）为单位遍历字符串，确保正确处理多字节字符；而使用简单的 for 循环则会以字符串的字节数为单位遍历字符串，可能导致在处理多字节字符时出现问题。")
	count := 0
	st := "中国abc"
	fmt.Println("字符串长度", len(st))
	for i, v := range st {
		count++
		fmt.Printf("循环次数%d: 字符初始位置%d,ASCII码表示%c\n", count, i, v)
	}
	count = 0
	for i := 0; i < len(st); i++ {
		count++
		fmt.Printf("循环次数%d: 每个字符%c \n", count, st[i])
	}
}

// testExecuteTimes (测试for控制结构中，条件的执行次数)
func testExecuteTimes() {
	fmt.Println("测试for控制结构中条件的执行次数")
	for i, c := 0, record(); i < c; i++ { // 初始化语句的record()只执行一次
		println("i", i)
	}

	c := 0
	for c < record() { // 条件表达式中的record()重复执行
		println("b", c)
		c++
	}
	println()
}

func record() int {
	print("record.")
	return 3
}

// testRangeKey (v 始终为集合中对应索引的值拷贝、不断被重复赋值，v的地址始终都是不变的。)
func testRangeKey() {
	fmt.Println("v 始终为集合中对应索引的值拷贝、不断被重复赋值，v的地址始终都是不变的。")
	arr := []int{1, 2, 3}
	newArr := []*int{}
	oldArr := []*int{}
	// range的v地址一直不变，每遍历一个值，把值赋值给v
	for i, v := range arr {
		oldArr = append(oldArr, &arr[i])
		newArr = append(newArr, &v)
	}
	fmt.Println("原始数据", arr) // 1 2 3
	for i, v := range oldArr {
		fmt.Printf("oldArr[%d] ==> %d\n", i, *v) // 1 2 3
	}
	for i, v := range newArr {
		fmt.Printf("newArr[%d] ==> %d\n", i, *v) // 3 3 3
	}
}

// testForType (for range也能用于实现了range方法的自定义类型，因为struct中有类型是slice。)
type MyType struct {
	data []int
}

func (m MyType) Range() []int {
	return m.data
}

func testForType() {
	fmt.Println("for range也能用于实现了range方法的自定义类型,因为struct中有类型是slice。")
	m := MyType{data: []int{1, 2, 3, 4, 5}}
	for index, value := range m.Range() {
		fmt.Printf("[%d] ==> %d \n", index, value)
	}
}
