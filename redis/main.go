package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c1, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c1.Close()
	/* 	c2, err := redis.DialURL("redis://127.0.0.1:6379")
	   	if err != nil {
	   		log.Fatalln(err)
	   	}
	   	defer c2.Close() */

	rec1, err := c1.Do("GET", "go")
	fmt.Println(string(rec1.([]byte)))

	// c2.Send("Get", "name")
	/* 	c2.Flush()
	   	rec2, err := c2.Receive()
	   	fmt.Println(string(rec2.([]byte))) */
}
