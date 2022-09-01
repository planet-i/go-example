package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/planet-i/go-example/826/controller/user_controller"
)

func Route() *gin.Engine {
	e := gin.Default()
	e.Handle("GET", "/hello", ErrorWrap(user_controller.Get))
	return e
}
