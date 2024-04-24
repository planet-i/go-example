package main

import (
	"fmt"
	"sort"
)

func main() {
	var target = 3
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(sort.Find(len(nums), func(i int) int {
		if nums[i] < target {
			return 1
		} else if nums[i] == target {
			return 0
		}
		return -1
	})) // 2 true

	nums2 := []int{7, 6, 5, 4, 2, 1}
	fmt.Println(sort.Find(len(nums2), func(i int) int {
		if nums2[i] > target {
			return 1
		} else if nums2[i] == target {
			return 0
		}
		return -1
	})) // 4 false
}

func Find(n int, cmp func(int) int) (i int, found bool) {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if cmp(h) > 0 {
			i = h + 1
		} else {
			j = h
		}
	}
	return i, i < n && cmp(i) == 0
}
