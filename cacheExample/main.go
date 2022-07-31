package main

import (
	"github.com/planet-i/goexample1/cacheExample/cache"
	http "github.com/planet-i/goexample1/cacheExample/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
