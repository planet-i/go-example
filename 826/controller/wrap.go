package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
)

type handler func(c *gin.Context) error

func ErrorWrap(h handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := h(c)
		if err != nil {
			log.Errorf("%+v\n", err)
		}
		return
	}
}
