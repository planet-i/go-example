package main

import (
	"fmt"

	"github.com/planet-i/go-example/scope/01_package-scope/vis"
)

func main() {
	fmt.Println(vis.CatName + " and " + vis.MouseName)
	vis.PrintVar()
}
