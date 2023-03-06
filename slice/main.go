package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	createSlice()
	// 切片反转
	s := []string{"a", "b", "c", "d", "e", "f", "g"}
	RecoverSlice(s)
	// 切片的深浅拷贝
	createBySlice()
	createByCopy()
	// 切片并发读写
	concurrentSliceNotForceIndex()
	concurrentSliceForceIndex()
	concurrentSliceWithMutex()
	concurrentSliceWithMutexPro()
	concurrentSliceWithChan()
	concurrentWriteMap()
}

// RecoverSlice 切片反转
func RecoverSlice(a []string) {
	var old []string
	copy(a, old)
	n := len(a)
	for i := 0; i < n/2; i++ {
		a[i], a[len(a)-i-1] = a[len(a)-i-1], a[i]
	}
	fmt.Printf("old=%v,new=%v", old, a)
}

func createSlice() {
	var s1 []int
	s2 := []int{}
	s3 := make([]int, 0)
	fmt.Println("var s1 []int 是nil?", s1 == nil)
	fmt.Println(":= []int{} 是nil?", s2 == nil)
	fmt.Println(":= make([]int, 0) 是nil?", s3 == nil)
	// fmt.Println(s2 == s3) 切片不可以做比较运算
	a := append(s1, 2) // 为nil的切片也可以append
	fmt.Println("往nil切片append后", a)
	s4 := make([]int, 5)
	s5 := make([]int, 3, 5)
	s6 := []int{10, 20, 4: 50, 5: 60}
	s7 := [...]int{5: 2} // s7是个数组

	fmt.Println("nil切片", s1, len(s1), cap(s1), s1 == nil)
	fmt.Println("[]int{}切片", s2, len(s2), cap(s2))
	fmt.Println("make([]int, 0) 切片", s3, len(s3), cap(s3))
	fmt.Println("make([]int, 5) 切片", s4, len(s4), cap(s4))
	fmt.Println("make([]int, 3,5) 切片", s5, len(s5), cap(s5))
	fmt.Println("切片定义", s6, len(s6), cap(s6))
	fmt.Println("数组", s7, len(s7), cap(s7))
}
func createBySlice() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	b := a[4:7] // 左闭右开、浅拷贝，共用同一个底层数组
	fmt.Printf("before sorting a = %v, b = %v\n", a, b)
	// slice排序
	sort.Slice(a, func(i, j int) bool {
		return a[i] > a[j]
	})
	fmt.Printf("after sorting a = %v, b = %v\n", a, b) //  a = [9 8 7 6 5 4 3 2 1], b = [5 4 3]
}

func createByCopy() {
	c := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	d := make([]int, 3)
	copy(d, c[4:7]) // copy是深拷贝，即不共用同一个底层数组
	fmt.Printf("before sorting c = %v, d = %v\n", c, d)

	sort.Slice(c, func(i, j int) bool {
		return c[i] > c[j]
	})
	fmt.Printf("after sorting c = %v, d = %v\n", c, d) // c = [9 8 7 6 5 4 3 2 1], d = [5 6 7]

}

// -----------验证Go 的切片是否原生支持并发

// 不指定索引，动态扩容并发向切片添加数据
func concurrentSliceNotForceIndex() {
	sl := make([]int, 0)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			sl = append(sl, num)
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d [NotForceIndex]\n", len(sl), cap(sl))
	// fmt.Println(sl)
}

// --- 指定索引，指定容量并发向切片添加数据
func concurrentSliceForceIndex() {
	sl := make([]int, 100)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			sl[num] = num
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d [ForceIndex]\n", len(sl), cap(sl))
	// fmt.Println(sl)
}

// --- 加互斥锁解决切片并发安全的问题,避免多个 goroutine 同时修改 slice
func concurrentSliceWithMutex() {
	sl := make([]int, 0)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			mu.Lock()
			sl = append(sl, num)
			mu.Unlock()
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d [WithMutex]\n", len(sl), cap(sl))
}

type SafeSlice struct {
	mu sync.Mutex
	sl []int
}

func (s *SafeSlice) Append(val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sl = append(s.sl, val)
}

func concurrentSliceWithMutexPro() {
	sl := SafeSlice{sl: make([]int, 0)}
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			sl.Append(num)
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d [WithMutexPro]\n", len(sl.sl), cap(sl.sl))
	// fmt.Println(sl)
}

// ----- 使用channel 串行操作来解决切片并发读写的问题
func concurrentSliceWithChan() {
	// 定义切片
	sl := make([]int, 0)
	// 定义 channel
	ch := make(chan func())
	// 定义一个协程从 channel 中取出函数并执行
	go func() {
		for f := range ch {
			f()
		}
	}()
	// 定义多个协程，将修改切片的操作封装成一个函数，并发送到 channel 中
	var wg sync.WaitGroup
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f := func() {
				sl = append(sl, num)
			}
			ch <- f
		}(k)
	}
	// 等待所有协程执行完成
	wg.Wait()
	// 关闭 channel
	close(ch)
	// 打印最终结果
	fmt.Printf("final len(sl)=%d cap(sl)=%d [WithChan]\n", len(sl), cap(sl))
	// fmt.Println(sl)
}

// ----- 使用sync.Map代替切片来解决并发读写的问题
func concurrentWriteMap() {
	var m sync.Map
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			// 存储键值对
			m.Store(num, num)
		}(k)
	}
	wg.Wait()
	// 遍历键值对
	len := 0
	m.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	fmt.Printf("final len(sl)=%d [WriteMap]\n", len)
}
