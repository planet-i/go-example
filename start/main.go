package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planet-i/goexample0/start/model"
)

func main() {
	fmt.Println("Hello")
	var people model.App
	people.ID = 100
	people.Name = "strr"
	people.UUID = "sty"
	fmt.Println(people)

	router := gin.Default()

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}

//报错在于未上传
