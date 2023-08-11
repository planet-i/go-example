package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	_, err := c.AddFunc("@every 1s", func() {
		fmt.Println("定时任务执行")
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}
