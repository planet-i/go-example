package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "a:b:c"
	sep := ":"
	str := Split(s, sep)
	fmt.Println(str)
}

func Split(s, sep string) (result []string) {
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+1:]
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}
