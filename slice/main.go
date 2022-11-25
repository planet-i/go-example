package main

import (
	"fmt"
	"sort"
)

type a struct {
	aa int
	bb bool
}

func main() {
	// var c a
	// c.aa = 1
	// var s1 []int
	// a := append(s1, 2)
	// fmt.Println(s1)
	// fmt.Println(a)
	// s2 := []int{}
	// s3 := make([]int, 0)
	// fmt.Println(s1 == nil)
	// fmt.Println(s2 == nil)
	// fmt.Println(s3 == nil)
	// // fmt.Println(s2 == s3) 切片不可以做比较运算
	// s4 := make([]int, 5)
	// s5 := make([]int, 3, 5)
	// s6 := []int{10, 20, 4: 50, 5: 60}
	// s7 := [...]int{5: 2} // s7是个数组

	// fmt.Println(s1, len(s1), cap(s1))
	// fmt.Println(s2, len(s2), cap(s2))
	// fmt.Println(s3, len(s3), cap(s3))
	// fmt.Println(s4, len(s4), cap(s4))
	// fmt.Println(s5, len(s5), cap(s5))
	// fmt.Println(s6, len(s6), cap(s6))
	// fmt.Println(s7, len(s7), cap(s7))
	// // 切片反转
	// RecoverSlice(s6)
	// fmt.Println(s6, len(s6), cap(s6))
	createSlice()
	createSliceCopy()
}

// RecoverSlice 切片反转
func RecoverSlice(a []int) {
	n := len(a)
	for i := 0; i < n/2; i++ {
		a[i], a[len(a)-i-1] = a[len(a)-i-1], a[i]
	}
}

func createSlice() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	b := a[4:7]
	fmt.Printf("before sorting a, b = %v\n", b) // before sorting a, b = [5 6 7]

	sort.Slice(a, func(i, j int) bool {
		return a[i] > a[j]
	})
	fmt.Printf("after sorting a, b = %v\n", b) // after sorting a, b = [5 4 3]
}

func createSliceCopy() {
	c := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	d := make([]int, 3)
	copy(d, c[4:7])
	fmt.Printf("before sorting c, d = %v\n", d) // before sorting c, d = [5 6 7]

	sort.Slice(c, func(i, j int) bool {
		return c[i] > c[j]
	})
	fmt.Printf("after sorting c, d = %v\n", d) //after sorting c, d = [5 6 7]
}
