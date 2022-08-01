package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type StructA struct {
	A string `yam1:"a"`
}
type StructB struct {
	StructA `yam1:".inline"`
	B       string `yam1:"b"`
}

var b StructB

func main() {
	var data string
	err := yaml.Unmarshal([]byte(data), &b)
	fmt.Println(err)
}
