package main

import (
	"container/list"
	"fmt"
)

func main() {
	a := list.New() // var a list.List
	a.PushBack("add")
	a.PushFront(123)
	//这两个方法都会返回一个 *list.Element 结构，如果在以后的使用中需要删除插入的元素，则只能通过 *list.Element 配合 Remove() 方法进行删除
	element := a.PushBack("fist")
	a.InsertAfter("high", element)
	a.InsertBefore("noon", element)
	// 使用
	a.Remove(element)
	//Front() 函数获取头元素,遍历时只要元素不为空就可以继续进行
	for i := a.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}

}

//列表并没有具体元素类型的限制
