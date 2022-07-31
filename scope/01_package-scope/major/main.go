package main

import (
	"fmt"

	"github.com/planet-i/goexample0/04_scope/01_package-scope/vis"
)

func main() {
	fmt.Println(vis.CatName + " and " + vis.MouseName)
	vis.PrintVar()
}
