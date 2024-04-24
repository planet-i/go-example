package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func main() {
	inputStrings := []string{"your_name", "yourName", "Your_name"}
	for _, str := range inputStrings {
		result, err := processString(str)
		if err != nil {
			fmt.Printf("输入：%s，输出：%s\n", str, err.Error())
		}
		fmt.Printf("输入：%s，输出：%s\n", str, result)
	}
}

// 输入字符串 your_name yourName Your_name 输出结果 yourName your_name error
// 解析：下划线命名法就改成小驼峰、小驼峰命名法就改为下划线、
func processString(str string) (res string, err error) {
	if len(str) > 0 && unicode.IsLower(rune(str[0])) {
		if strings.Contains(str, "_") {
			res = snakeToCamel(str)
		} else {
			res = camelToSnake(str)
		}
	} else {
		return "", errors.New("格式错误")
	}
	return
}

// snakeToCamel 将蛇形命名转换为驼峰命名
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var camelCaseString string
	for _, part := range parts {
		camelCaseString += strings.Title(part)
	}
	return camelCaseString
}

// camelToSnake 将驼峰命名转换为蛇形命名
func camelToSnake(s string) string {
	var builder strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteByte('_')
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
