package main

import "fmt"

func main() {
	fmt.Println("Bin        \tOct \t Dec \t Hex \t ASCII")
	for i := 0; i < 127; i++ {
		fmt.Printf("%b        \t%o \t %d \t %#X \t %q \n", i, i, i, i, i)
	}
}
