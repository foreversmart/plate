package ginroute

import (
	"github.com/foreversmart/plate/errors"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/route"
	"github.com/foreversmart/plate/utils/request"
	"github.com/foreversmart/plate/utils/val"
	"github.com/gin-gonic/gin"
	"net/http"
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
	recover    route.Recover
	isClose    bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.path = "/"
	g.middleware = make([]*route.Middleware, 0, 5)

	// set router default recover
	g.recover = GinRecovery
	return g
}

func (g *GinRouter) SetRecover(res route.Recover) {
	g.recover = res
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

	handleArgs := make([]interface{}, 0, 5)

	// TODO check v type must be struct
	g.Engine.Handle(method, path, func(c *gin.Context) {
		// connection come after server is closed
		if g.isClose {
			c.JSON(500, nil)
			return
		}

		// recover
		defer func() {
			if recV := recover(); recV != nil {
				resp, err := g.recover(recV, c.Request, handleArgs)
				if err != nil {
					handleError(c, err)
					return
				}

				c.JSON(200, resp)

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
			handleError(c, err)
			return
		}

		for _, mid := range g.middleware {
			mv := reflect.ValueOf(mid.V)
			mv = val.SettableValue(mv)
			mt := mv.Type()
			nmv := reflect.New(mt)

			err = parser.Parse(nmv.Elem())

			if err != nil {
				logger.StdLog.Error(err)
				handleError(c, err)
				return
			}

			midArgs := nmv.Interface()
			handleArgs = append(handleArgs, midArgs)
			res, err := mid.H(midArgs)
			if err != nil {
				logger.StdLog.Error(err)
				handleError(c, err)
				return
			}

			parser.WithMid(res)
		}

		err = parser.Parse(nv.Elem())
		if err != nil {
			logger.StdLog.Error(err)
			handleError(c, err)
			return
		}

		reqArg := nv.Interface()
		handleArgs = append(handleArgs, reqArg)

		// do before
		resp, _ := handler(reqArg)

		// do after
		c.JSON(200, resp)
	})
}

func handleResp(c *gin.Context, resp interface{}, err error) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(e.Code, e)
	}

	ne := errors.BadRequestError(err.Error())
	c.JSON(ne.Code, ne)
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
