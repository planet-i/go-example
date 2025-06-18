package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestSpecailLetter(t *testing.T) {
	str := `[]{}~/!:+_()"^\\/\\:*?"<>|()（）\'、;-=!@#$%^&`
	var chars []rune
	for _, letter := range str {
		ok, letters := SpecialLetters(letter)
		if ok {
			chars = append(chars, letters...)
		} else {
			chars = append(chars, letter)
		}
	}
	fmt.Println(string(chars))
}

func TestStringSplitScenarios(t *testing.T) {
	tests := []struct {
		name         string   // 测试场景描述
		input        string   // 输入字符串
		wantEmpty    bool     // 验证 input == "" 的预期结果
		wantZeroLen  bool     // 验证 len(elements) == 0 的预期结果
		wantElements []string // 期望的切分结果
	}{
		{
			name:         "empty string",
			input:        "",
			wantEmpty:    true,
			wantZeroLen:  false, // strings.Split("", ",") 返回 [""],此时 len(elements) == 1
			wantElements: []string{""},
		},
		{
			name:         "single comma",
			input:        ",",
			wantEmpty:    false,
			wantZeroLen:  false,
			wantElements: []string{"", ""},
		},
		{
			name:         "trailing comma",
			input:        "1,",
			wantEmpty:    false,
			wantZeroLen:  false,
			wantElements: []string{"1", ""},
		},
		{
			name:         "normal case",
			input:        "a,b,c",
			wantEmpty:    false,
			wantZeroLen:  false,
			wantElements: []string{"a", "b", "c"},
		},
		{
			name:         "no comma",
			input:        "hello",
			wantEmpty:    false,
			wantZeroLen:  false,
			wantElements: []string{"hello"},
		},
		{
			name:         "multiple commas",
			input:        "1,,3",
			wantEmpty:    false,
			wantZeroLen:  false,
			wantElements: []string{"1", "", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行切分操作
			elements := strings.Split(tt.input, ",")

			// 验证 input == ""
			if got := (tt.input == ""); got != tt.wantEmpty {
				t.Errorf("input == '' 结果为 %v，期望 %v", got, tt.wantEmpty)
			}

			// 验证 len(elements) == 0
			if got := (len(elements) == 0); got != tt.wantZeroLen {
				t.Errorf("len(elements) == 0 结果为 %v，期望 %v", got, tt.wantZeroLen)
			}

			// 验证切分结果
			if !reflect.DeepEqual(elements, tt.wantElements) {
				t.Errorf("切分结果 %v 不等于期望值 %v", elements, tt.wantElements)
			}
		})
	}
}
