package main

import (
	"context"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8887/sleep/3", nil)
	if err != nil {
		panic(err)
	}
	do, err := http.DefaultClient.Do(req)
	spew.Dump(do, err)
}
