package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args)
	keyword := os.Args[1]
	test(keyword)

}

func test(keyword string) {
	keyword = strings.ReplaceAll(keyword, "\\", "\\\\")
	keyword = strings.ReplaceAll(keyword, "%", "\\%")
	keyword = strings.ReplaceAll(keyword, "'", "''")
	fmt.Println(keyword)
}
