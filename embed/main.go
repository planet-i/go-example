package main

import "fmt"

func main() {
	StructToStruct()
}

// ---- 结构体嵌入结构
type Base struct {
	name string
	tag  string
}

func (base Base) DescribeName() string {
	return fmt.Sprintf("base name %s ", base.name)
}

func (base Base) DescribeTag() string {
	return fmt.Sprintf("base tag %s", base.tag)
}

type Container struct { // Container 是嵌入结构体 [外部类型]
	Base           // Base 是被嵌入的结构体  [内部类型]
	address string // 注意，直接嵌入可能导致内部类型的导出字段也可导出，
	tag     string
}

func (con Container) DescribeTag() string {
	return fmt.Sprintf("container tag %s", con.tag)
}

func StructToStruct() {
	co := Container{}
	co.name = "内部类型" // 外部类型可直接为内部类型的字段赋值
	co.address = "外部类型"
	fmt.Println(co.DescribeName()) // 外部类型可直接调用内部类型的方法
	fmt.Printf("co -> {name: %v, address: %v}\n", co.name, co.address)

	co1 := Container{Base: Base{name: "aa"}, address: "bb"} // 当使用结构体字面量时，我们需要将内部类型整体初始化，而不是单单对其字段初始化。
	fmt.Printf("co -> {name: %v, address: %v}\n", co1.name, co1.address)

	b2 := Base{name: "内部类型", tag: "b's tag"}
	co2 := Container{Base: b2, address: "外部类型", tag: "co's tag"}

	fmt.Println(co2.DescribeTag(), co2.tag)           // 直接用外部类型调用，覆盖的字段和方法都是用外部类型的
	fmt.Println(co2.Base.DescribeTag(), co2.Base.tag) // 通过外部类型.内部类型名.字段/方法名调用
	fmt.Println(b2.DescribeTag(), b2.tag)
	// 如果多个内部类型字段/方法重复，并且不与外部类型重复。直接外部类型调用会报错，不知道选择哪个字段/方法
}
