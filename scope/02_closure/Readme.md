## 闭包

##### 什么是闭包

```
// 1 1 2 3 5 8 13 21 34 55
//   a b
//     a b
func fibonacci() func() int {
​    a, b := 0, 1
​    return fnuc() int {
​        a, b = b, a+
​        return a
​    }
}

func TestFibonacci(t *testing.T){
​    f := fifibonacci()
​    for i := 0;i < 10;i++{
​        fmt.Println(f()," ")
​    }
​    fmt.Println()
}
```
Go语言支持连续赋值，无需中间变量（a, b = b , a   表示变量值直接交换）