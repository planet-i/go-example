package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/mozillazg/go-pinyin"
)

func main() {
	a := firstLetter("世c界yug")
	fmt.Println(a)
	b := FirstLetterOfPinYin("世界shijirhijie")
	fmt.Println(b)
	c := FirstLetterOfString("g世界hi")
	fmt.Println(c)
}

// 获取字符串的第一个字
func firstLetter(s string) string {
	_, size := utf8.DecodeRuneInString(s)
	return s[:size]
}

func firstLetter1(s string) string {
	for _, l := range s {
		return string(l)
	}
	return ""
}

func FirstLetterOfPinYin(r string) string {
	var a = pinyin.NewArgs()
	result := pinyin.Pinyin(r, a)
	fmt.Println(result)
	return string(result[0][0][0])
}

// 获取字符串首字母
func FirstLetterOfString(s string) string {
	l := []rune(s)[:1][0]
	if unicode.Is(unicode.Han, l) {
		fmt.Println("是一个汉字", string(l))
		var a = pinyin.NewArgs()
		result := pinyin.Pinyin(string(l), a)
		return string(result[0][0][0])
	} else if unicode.IsNumber(l) {
		fmt.Println("是一个数字", string(l))
		return "~"
	} else if unicode.IsLetter(l) {
		fmt.Println("是一个字母", string(l))
		a := string(unicode.ToUpper(l))
		b := string(unicode.ToLower(l))
		fmt.Println(string(a), string(b))
		return a
	}
	return ""
}
