package router

import (
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/utils/tagger"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type GinRouter struct {
	Engine  *gin.Engine
	Config  *ConfigType
	wg      sync.WaitGroup
	isClose bool
}

func NewGinRouter(c *ConfigType) *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.Config = c
	return g
}

type GinRequest struct {
	c *gin.Context
}

func NewGinRequest(c *gin.Context) *GinRequest {
	return &GinRequest{
		c,
	}
}

func (g *GinRequest) Request() *http.Request {
	return g.c.Request
}

func (g *GinRequest) Params() map[string][]string {
	res := make(map[string][]string)
	for _, param := range g.c.Params {
		if v, ok := res[param.Key]; ok {
			v = append(v, param.Value)
			continue
		}

		// first init
		res[param.Key] = []string{param.Value}
	}

	return res
}

func (r *GinRouter) Handle(method, path string, handler Handler, v interface{}) {
	vv := reflect.ValueOf(v)
	vt := vv.Type().Elem()

	// TODO check v type must be struct

	r.Engine.Handle(method, path, func(c *gin.Context) {
		// connection come after server is closed
		if r.isClose {
			c.JSON(500, nil)
			return
		}

		// count connection num for close
		r.wg.Add(1)
		defer r.wg.Done()

		req := NewGinRequest(c)
		nv := reflect.New(vt)

		err := tagger.ParseRequest(req, nv.Elem())
		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		resp, _ := handler(nv.Interface())
		c.JSON(200, resp)
	})
}

func (r *GinRouter) Run(addr ...string) {
	r.Engine.Run(addr...)
}

func (r *GinRouter) Wait(timeout int) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		r.wg.Wait()
	}()

	select {
	case <-c:
		return // completed normally
	case <-time.After(time.Second * time.Duration(timeout)):
		return // timed out
	}
}
