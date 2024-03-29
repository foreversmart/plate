package ginroute

import (
	"github.com/foreversmart/plate/errors"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/route"
	"github.com/foreversmart/plate/utils/request"
	"github.com/foreversmart/plate/utils/val"
	"github.com/foreversmart/plate/utils/view"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type GinRouter struct {
	engine            *gin.Engine
	beforeMid         []*route.Middle
	afterMid          []*route.Middle
	parentAfterMidNum int          // mark after mid is parent router or sub itself
	path              string       // url path
	subs              []*GinRouter // sub routers
	wg                sync.WaitGroup
	recover           route.Recover
	autoCors          bool                // auto handle cors options request
	routeMap          map[string][]string // record route path map

	isClose bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.engine = gin.New()
	g.path = "/"
	g.beforeMid = make([]*route.Middle, 0, 5)
	g.afterMid = make([]*route.Middle, 0, 5)
	g.routeMap = make(map[string][]string)

	// set router default recover
	g.recover = GinRecovery
	return g
}

func (g *GinRouter) SetRecover(res route.Recover) {
	g.recover = res
}

func (g *GinRouter) SetAutoCors(isOpen bool) {
	g.autoCors = isOpen
}

func (g *GinRouter) addCors() {
	if g.autoCors {
		for path, methods := range g.routeMap {
			methodsStr := strings.Join(methods, ", ")
			g.engine.Handle(http.MethodOptions, path, func(c *gin.Context) {
				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Methods", methodsStr)
				c.Header("Access-Control-Allow-Headers", "x-token, content-type")
				c.JSON(200, nil)
			})
		}
	}

	for _, sub := range g.subs {
		sub.addCors()
	}
}

func (g *GinRouter) Sub(relativePath string) route.Router {
	ng := &GinRouter{
		engine:   g.engine,
		path:     joinPaths(g.path, relativePath),
		routeMap: make(map[string][]string),
	}

	ng.beforeMid = make([]*route.Middle, len(g.beforeMid))
	ng.afterMid = make([]*route.Middle, len(g.afterMid))
	copy(ng.beforeMid, g.beforeMid)
	copy(ng.afterMid, g.afterMid)

	g.subs = append(g.subs, ng)
	ng.parentAfterMidNum = len(g.afterMid)
	ng.autoCors = g.autoCors

	return ng
}

func (g *GinRouter) AddMiddleBefore(handler route.Handler, v interface{}) {
	g.beforeMid = append(g.beforeMid, &route.Middle{
		H: handler,
		V: v,
	})

}

// AddMiddleAfter add after middle
func (g *GinRouter) AddMiddleAfter(handler route.Handler, v interface{}) {
	g.afterMid = append(g.afterMid, &route.Middle{
		H: handler,
		V: v,
	})

	if len(g.afterMid) == 1 {
		return
	}

	// move the new middle after to the right place
	for i := len(g.afterMid) - 1; i > len(g.afterMid)-g.parentAfterMidNum-1; i-- {
		g.afterMid[i], g.afterMid[i-1] = g.afterMid[i-1], g.afterMid[i]
	}
}

type ReqMeta struct {
	Method    string `json:"method" plate:"method,mid"`
	Path      string `json:"path" plate:"path,mid"`
	ParamName string `json:"param_name" plate:"param_name,mid"`
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

	meta := &ReqMeta{
		Method:    method,
		Path:      path,
		ParamName: vt.Name(),
	}

	// TODO uniform different path form
	absPath := joinPaths(g.path, path)
	g.routeMap[absPath] = append(g.routeMap[absPath], method)

	g.engine.Handle(method, absPath, func(c *gin.Context) {
		if g.autoCors {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", method)
			c.Header("Access-Control-Allow-Headers", "x-token, content-type")
		}

		handleArgs := make([]interface{}, 0, 5)

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

		// add request meta info
		parser.WithMid(meta)

		// do before mid
		for _, mid := range g.beforeMid {
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

			err = parser.WithMid(res)
			if err != nil {
				logger.StdLog.Error(err)
				handleError(c, err)
				return
			}
		}

		err = parser.Parse(nv.Elem())
		if err != nil {
			logger.StdLog.Error(err)
			handleError(c, err)
			return
		}

		reqArg := nv.Interface()
		handleArgs = append(handleArgs, reqArg)

		// do handle function
		resp, err := handler(reqArg)
		if err != nil {
			handleError(c, err)
			return
		}

		// do after mid
		for _, mid := range g.afterMid {
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

			err = parser.WithMid(res)
			if err != nil {
				logger.StdLog.Error(err)
				handleError(c, err)
				return
			}
		}

		if resp != nil {
			respView, err := view.FetchViewFromStruct(reflect.ValueOf(resp), false, request.TagNameFetch, request.LocResp)
			if err != nil {
				logger.StdLog.Error(err)
				c.JSON(200, resp)
				return
			}

			// return handler resp
			ctrResp := &CommonRespCtr{
				Code: http.StatusOK,
			}
			err = respView.SetObjectValue(ctrResp, false, request.TagNameFetch, request.LocResp)
			if err != nil {
				logger.StdLog.Error(err)
				c.JSON(200, resp)
				return
			}

			for k, v := range ctrResp.Header {
				c.Header(k, v)
			}
			c.JSON(ctrResp.Code, resp)
			return
		}

		c.JSON(http.StatusOK, resp)

	})
}

type CommonRespCtr struct {
	Code   int               `json:"-" plate:"code,resp"`
	Header map[string]string `json:"-" plate:"header,resp"`
}

func handleResp(c *gin.Context, resp interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	}

	handleError(c, err)
}

func handleError(c *gin.Context, err error) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(e.Code, e)
		return
	}

	ne := errors.BadRequestError(err.Error())
	c.JSON(ne.Code, ne)
}

func (g *GinRouter) Wait(timeout int) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		for _, sub := range g.subs {
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
	for _, sub := range g.subs {
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
