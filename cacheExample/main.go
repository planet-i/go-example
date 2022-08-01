package main

import (
	"github.com/planet-i/go-example/cacheExample/cache"
	http "github.com/planet-i/go-example/cacheExample/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
