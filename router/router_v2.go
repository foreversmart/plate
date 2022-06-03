package router

import (
	"encoding/json"
	"github.com/foreversmart/plate/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"reflect"
)

type V2Router interface {
	Handle(method, path string, handler HanldeV2, v interface{})
	Run(addr ...string)
}

type GinRouter struct {
	Engine *gin.Engine
	Config *ConfigType
}

func NewGinRouter(c *ConfigType) *GinRouter {
	g := &GinRouter{}
	g.Engine = gin.New()
	g.Config = c
	return g
}

type HanldeV2 func(req interface{}) (resp interface{}, err error)

//

func (r *GinRouter) Handle(method, path string, handler HanldeV2, v interface{}) {
	vv := reflect.ValueOf(v)
	vt := vv.Type().Elem()

	r.Engine.Handle(method, path, func(c *gin.Context) {
		// TODO
		nv := reflect.New(vt)
		nv.Elem().Field(i
		nv.Elem().Set(vv.Elem())
		nvi := nv.Interface()

		body, err := ioutil.ReadAll(c.Request.Body)
		c.Request.
		defer c.Request.Body.Close()
		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		err = json.Unmarshal(body, &nvi)

		if err != nil {
			logger.StdLog.Error(err)
			c.JSON(400, err)
			return
		}

		resp, _ := handler(nvi)
		c.JSON(200, resp)
	})
}

func (r *GinRouter) Run(addr ...string) {
	r.Engine.Run(addr...)
}
