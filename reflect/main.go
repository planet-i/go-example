package main

import (
	"fmt"
	"reflect"
)

type user struct {
	name string
	age  int
}
type manager struct {
	user
	title string
}

func main() {
	var m manager
	t := reflect.TypeOf(&m)
	if t.Kind() == reflect.Ptr { //如果t是指针类型
		t = t.Elem() //获取指针的基类型
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Name, f.Type, f.Offset)

		if f.Anonymous {
			for x := 0; x < f.Type.NumField(); x++ {
				af := f.Type.Field(x)
				fmt.Println(" ", af.Name, af.Type)
			}
		}
	}

	t1 := reflect.TypeOf(m)
	name, _ := t1.FieldByName("name") //按名称查找
	fmt.Println(name.Name, name.Type)

	age := t1.FieldByIndex([]int{0, 1}) //按多级索引查找
	fmt.Println(age.Name, age.Type)

}
