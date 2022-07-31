package main

func main(){
	a,b,c,d,x := 1,2,3,4,100
	switch x:=2;x{
	case a,b:
		println("a | b")
	case c:
		println("c")
	case d:
	case 5:
		println("4")
	default:
		println("inequality")
	}
	println(x)
}
/* switch同样支持初始化语句，按从上到下，从左到右的顺序匹配case执行，只有全部失败时，才执行default
无需显式执行break语句，case执行完毕后自动中断 */