package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

// server 1
func main() {
	server := gin.New()
	server.GET("/sleep/:time", func(c *gin.Context) {
		t := c.Param("time")
		t1, _ := strconv.Atoi(t)
		time.Sleep(time.Duration(t1) * time.Second)

		ctx, cancelFun := context.WithTimeout(c.Request.Context(), time.Second*5)
		defer cancelFun()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8888/sleep/5", nil)
		if err != nil {
			panic(err)
		}
		do, err := http.DefaultClient.Do(req)
		spew.Dump(do, err)

		c.String(http.StatusOK, "sleep "+t)
	})
	server.Run(":8887")
}
