package main

import (
	"fmt"
	"unsafe"
)

func main() {
	number := 5
	pointer := &number
	fmt.Printf("number:addr:%p, value: %d\n", pointer, *pointer)
	float32Number := (*float32)(unsafe.Pointer(pointer))
	*float32Number = *float32Number + 3
	fmt.Printf("float64:addr:%p, value:%f\n", float32Number, *float32Number)
}
