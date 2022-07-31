package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()                         //实例化一个gin示例
	r.GET("/users/:id", func(c *gin.Context) { // 为首页(/)的GET访问注册一个处理函数
		id := c.Param("id")
		c.String(200, "The user id is  %s", id)
	})
	r.GET("/users/:id/*name", func(c *gin.Context) {
		id := c.Param("id")
		name := c.Param("name")
		message := id + " belong " + name
		c.String(200, message)
	})
	r.Run(":8080") //启动了一个HTTP服务，端口是8080
}

//:id就是一个路由参数，我们可以通过c.Param("id")获取定义的路由参数的值
