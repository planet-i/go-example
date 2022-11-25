package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
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

/**
判断是否为字母： unicode.IsLetter(v)
判断是否为十进制数字： unicode.IsDigit(v)
判断是否为数字： unicode.IsNumber(v)
判断是否为空白符号： unicode.IsSpace(v)
判断是否为Unicode标点字符 :unicode.IsPunct(v)
判断是否为中文：unicode.Han(v)
*/
func SpecialLetters(letter rune) (bool, []rune) {
	if unicode.IsPunct(letter) || unicode.IsSymbol(letter) {
		var chars []rune
		chars = append(chars, '\\', '\\', letter)
		return true, chars
	}
	return false, nil
}
