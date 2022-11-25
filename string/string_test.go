package main

import (
	"fmt"
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
