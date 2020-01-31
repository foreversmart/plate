package controller

import (
	"github.com/foreversmart/plate/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	// set gin release mode
	gin.SetMode(gin.ReleaseMode)

}

type Handler func(e *gin.Engine) *gin.Engine

// NewRoute 初始化路由
func NewRoute(subRoute Handler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	return subRoute(r)
}
