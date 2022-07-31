package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) { // /welcome?firstname=Bob&lastname=kai&hobby=swimming&hobby=sing
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s %s", firstname, lastname, c.QueryArray("hobby"))
	})
	r.GET("/map", func(c *gin.Context) { // /map/?ids[a]=123&ids[b]=456&ids[c]=789
		c.JSON(200, c.QueryMap("ids"))
	})
	r.Run(":8080")

}
