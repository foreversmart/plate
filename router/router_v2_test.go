package router

import (
	"math/rand"
	"net/http"
	"testing"
)

type Business struct {
	//Get string `plage:"get:user/"`
}

type Request struct {
	A int `json:"a"`
}

type Resp struct {
	Out int `json:"out"`
}

var (
	_Business *Business
)

func (b *Business) Get(req interface{}) (resp interface{}, err error) {
	r := req.(Request)
	res := &Resp{}
	res.Out = r.A + 10 + rand.Int()
	return res, nil
}

func TestNewGinRouter(t *testing.T) {
	route := NewGinRouter(nil)
	route.Handle(http.MethodGet, "/", _Business.Get, &Request{})
	route.Run("8080")
}
