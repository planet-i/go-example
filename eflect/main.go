package main

import (
	"fmt"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Student struct {
	Name string "学生姓名"
	Age  int    `a:"1111"b:"3333"`
}

func main() {
	s := Student{}
	rt := reflect.TypeOf(s)

	fieldName, ok := rt.FieldByName("Name")
	if ok {
		fmt.Println(fieldName.Tag)
	}

	fieldAge, ok2 := rt.FieldByName("Age")
	if ok2 {
		fmt.Println(fieldAge.Tag.Get("a"))
		fmt.Println(fieldAge.Tag.Get("b"))
	}

	fmt.Println("type_Name: ", rt.Name())
	fmt.Println("type_NumField: ", rt.NumField())
	fmt.Println("type_PkgPath: ", rt.PkgPath())
	fmt.Println("type_String: ", rt.String())

	fmt.Println("type.Kind.String: ", rt.Kind().String())
	fmt.Println("type.String() = ", rt.String())
	// 获取结构类型的字段名称
	for i := 0; i < rt.NumField(); i++ {
		fmt.Printf("type.Field[%d].Name: %v", i, rt.Field(i).Name)
		fmt.Println()
	}

	f, err := excelize.OpenFile("../34_excelize/Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	cell := f.GetCellValue("sheet1", "B2")
	fmt.Println(cell)
}
