package router

import (
	"fmt"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	// set gin release mode
	gin.SetMode(gin.ReleaseMode)

}

type Handler func(e *gin.Engine) *gin.Engine

type Router struct {
	Engine *gin.Engine
	Config *ConfigType
}

func NewRouter(c *ConfigType, handle Handler) *Router {
	engine := NewRoute(handle)
	router := &Router{
		Engine: engine,
		Config: c,
	}

	return router
}

func (r *Router) Run() {
	logger.StdLog.Infof("server started at %s:%d", r.Config.Host, r.Config.Port)
	r.Engine.Run(fmt.Sprintf("%s:%d", r.Config.Host, r.Config.Port))
}

// NewRoute 初始化路由
func NewRoute(handle Handler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	return handle(r)
}
