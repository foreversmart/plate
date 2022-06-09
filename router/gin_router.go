package router

import (
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/utils/tag"
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
	"time"
)

type GinRouter struct {
	Engine     *gin.Engine
	middleware map[string][]*Middleware
	wg         sync.WaitGroup
	isClose    bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.middleware = make(map[string][]*Middleware)
	return g
}

type Middleware struct {
	H Handler
	V interface{}
}

func (r *GinRouter) Middle(path string, handler Handler, v interface{}) {
	if m, ok := r.middleware[path]; ok {
		m = append(m, &Middleware{
			handler,
			v,
		})
		return
	}

	r.middleware[path] = []*Middleware{{handler, v}}
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

		middlewares := r.middleware[path]
		for _, mid := range middlewares {
			res, err := mid.H(mid.V)
			if err != nil {
				c.JSON(400, nil)
				return
			}

		}

		nv := reflect.New(vt)
		err := tag.ParseRequest(req, nv.Elem())
		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		// do before

		resp, _ := handler(nv.Interface())

		// do after
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

func (r *GinRouter) Close() {
	r.isClose = true
}
