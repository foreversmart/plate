package router

import (
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/utils/request"
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
	"time"
)

type GinRouter struct {
	Engine     *gin.Engine
	middleware []*Middleware
	path       string
	wg         sync.WaitGroup
	isClose    bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.path = "/"
	g.middleware = make([]*Middleware, 0, 5)
	return g
}

func (r *GinRouter) Group(relativePath string) *GinRouter {
	ng := &GinRouter{
		Engine: r.Engine,
		path:   joinPaths(r.path, relativePath),
	}

	copy(ng.middleware, r.middleware)

	return ng
}

func (r *GinRouter) AddMiddle(handler Handler, v interface{}) {
	r.middleware = append(r.middleware, &Middleware{
		handler,
		v,
	})
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
		parser, err := request.NewParser(req)
		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		for _, mid := range r.middleware {
			mv := reflect.ValueOf(mid.V)
			mt := mv.Type().Elem()
			nmv := reflect.New(mt)
			err = parser.Parse(nmv.Elem())

			if err != nil {
				logger.StdLog.Error(err)
				c.JSON(400, nil)
				return
			}

			res, err := mid.H(nmv.Interface())
			if err != nil {
				logger.StdLog.Error(err)
				c.JSON(400, nil)
				return
			}

			parser.WithMid(res)
		}

		err = parser.Parse(nv.Elem())
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
