package main

import "fmt"

type node struct{
	_ bool      //忽略值输出的是类型默认值
	id int
	next *node
}

func main(){
	n1 := node{
		id : 100,
	}
	n2 := node{
		id : 200,
		next : &n1,
	}
	n3 := node{true,300,&n2} //按顺序初始化全部字段，不能遗漏
	fmt.Println(n1,n2,n3)
}