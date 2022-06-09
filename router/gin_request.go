package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
