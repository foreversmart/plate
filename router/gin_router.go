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
	Engine *gin.Engine

	wg      sync.WaitGroup
	isClose bool
}

func NewGinRouter() *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	return g
}

func Midlleware(handler Handler, v interface{}) {

}

func (r *GinRouter) Middleware(path string, handler Handler, v interface{}) {
	map[]
}

func (r *GinRouter) Handle(method, path string, handler Handler, v interface{}) {
	vv := reflect.ValueOf(v)
	vt := vv.Type().Elem()

	var midHandle Handler
	var midV interface{}

	// TODO check v type must be struct
	r.Engine.Handle(method, path, func(c *gin.Context) {
		// connection come after server is closed
		if r.isClose {
			c.JSON(500, nil)
			return
		}

		res, err := midHandle(midV)
		if err != nil {
			c.JSON(500, nil)
			return
		}



		// count connection num for close
		r.wg.Add(1)
		defer r.wg.Done()

		req := NewGinRequest(c)
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
