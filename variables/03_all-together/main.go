package main

var x = 100

func main() {
	//-------变量赋值的注意操作
	println(&x, x) //全局变量
	x := "abc"     //重新定义和初始化同名局部变量
	println(&x, x)

	{
		i := 100
		println(&i, i)
		i, j := 200, "zxc" //i退化为赋值操作，仅有j是变量重定义
		println(&i, i, j)
	}
	//--------多变量赋值操作：先计算右值，再赋值
	a, b := 1, 2
	a, b = b+3, a+2
	println(a, b)
	//--------变量的类型转换
	var y float32 = 3.1415
	z := int(y)
	println(z) //类型转换需要显式声明，只能发生在两种相互兼容的类型之间
}
