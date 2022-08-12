package ginroute

import (
	"fmt"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/logger/stack"
	"github.com/foreversmart/plate/route"
	"github.com/foreversmart/plate/utils/request"
	"github.com/foreversmart/plate/utils/val"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"reflect"
	"sync"
	"time"
)

type GinRouter struct {
	Engine     *gin.Engine
	middleware []*route.Middleware
	path       string
	subRoute   []*GinRouter
	wg         sync.WaitGroup
	recover    route.Handler
	isClose    bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.path = "/"
	g.middleware = make([]*route.Middleware, 0, 5)
	return g
}

func (g *GinRouter) Sub(relativePath string) *GinRouter {
	ng := &GinRouter{
		Engine: g.Engine,
		path:   joinPaths(g.path, relativePath),
	}

	copy(ng.middleware, g.middleware)
	g.subRoute = append(g.subRoute, ng)

	return ng
}

func (g *GinRouter) AddMiddle(handler route.Handler, v interface{}) {
	g.middleware = append(g.middleware, &route.Middleware{
		handler,
		v,
	})

}

func (g *GinRouter) Handle(method, path string, handler route.Handler, v interface{}) {
	if v == nil {
		panic("router handle v interface{} cant be nil")
	}

	vv := reflect.ValueOf(v)
	vv = val.SettableValue(vv)
	vt := vv.Type()

	if vt.Kind() != reflect.Struct {
		panic("router handle v interface{} must be struct kind !")
	}

	// TODO check v type must be struct
	g.Engine.Handle(method, path, func(c *gin.Context) {
		// connection come after server is closed
		if g.isClose {
			c.JSON(500, nil)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				g.recover(nil)
				if l != nil {
					stack := stack.Stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					msg := fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
					l.Println(msg)
					panicLogger.Error(msg)
				}
				c.AbortWithStatus(500)
			}
		}()

		// count connection num for close
		g.wg.Add(1)
		defer g.wg.Done()
		req := NewGinRequest(c)

		nv := reflect.New(vt)
		parser, err := request.NewParser(req)
		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		for _, mid := range g.middleware {
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

func (g *GinRouter) Wait(timeout int) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		for _, sub := range g.subRoute {
			sub.Wait(timeout)
		}
		g.wg.Wait()
	}()

	select {
	case <-c:
		return // completed normally
	case <-time.After(time.Second * time.Duration(timeout)):
		return // timed out
	}
}

func (g *GinRouter) Close() {
	g.isClose = true
	for _, sub := range g.subRoute {
		sub.Close()
	}
}
func (g *GinRouter) Get(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodGet, path, handler, req)
}
func (g *GinRouter) Head(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodHead, path, handler, req)
}
func (g *GinRouter) Post(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodPost, path, handler, req)
}
func (g *GinRouter) Put(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodPut, path, handler, req)
}
func (g *GinRouter) Patch(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodPatch, path, handler, req)
}
func (g *GinRouter) Delete(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodDelete, path, handler, req)
}
func (g *GinRouter) Connect(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodConnect, path, handler, req)
}
func (g *GinRouter) Options(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodOptions, path, handler, req)
}
func (g *GinRouter) Trace(path string, handler route.Handler, req interface{}) {
	g.Handle(http.MethodTrace, path, handler, req)
}
