package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	sText := "配额容量超出"
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted1 := textQuoted[1 : len(textQuoted)-1]
	fmt.Println(textUnquoted1)
	textUnquoted := "\u914d\u989d\u5bb9\u91cf\u8d85\u51fa"
	fmt.Println(textUnquoted)
	sUnicodev := strings.Split(textUnquoted1, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
	fmt.Println(context)
}

// s := []byte(`{"code":200,"msg":"\u914d\u989d\u5bb9\u91cf\u8d85\u51fa"}`)
// 	// v, _ := UnescapeUnicode(s)
// 	// fmt.Println(v)
// }

func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
