package ginroute

import (
	"github.com/foreversmart/plate/route"
)

type GinServer struct {
	Root    *GinRouter
	isClose bool
}

func NewGinServer() route.Server {
	return &GinServer{
		Root: NewGinRouter(),
	}
}

func (g *GinServer) Route() route.Router {
	return g.Root
}

func (g *GinServer) Run(addr ...string) {
	g.Root.addCors()
	panic(g.Root.engine.Run(addr...))
}

func (g *GinServer) Wait(timeout int) {
	g.Root.Wait(timeout)
}

func (g *GinServer) Close() {
	g.isClose = true
	g.Root.Close()
}
