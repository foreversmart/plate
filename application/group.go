package application

import (
	"github.com/foreversmart/plate/router"
	"net/http"
	"path"
)

type Group struct {
	r            router.Router
	relativePath string
}

func (a *Application) Group(relativePath string) *Group {
	return &Group{
		a.Router,
		relativePath,
	}
}

func (g *Group) Group(relativePAth string) {

}

func (g *Group) Get(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodGet, path, handler, req)
}
func (g *Group) Head(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodHead, path, handler, req)
}
func (g *Group) Post(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodPost, path, handler, req)
}
func (g *Group) Put(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodPut, path, handler, req)
}
func (g *Group) Patch(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodPatch, path, handler, req)
}
func (g *Group) Delete(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodDelete, path, handler, req)
}
func (g *Group) Connect(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodConnect, path, handler, req)
}
func (g *Group) Options(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodOptions, path, handler, req)
}
func (g *Group) Trace(path string, handler router.Handler, req interface{}) {
	g.r.Handle(http.MethodTrace, path, handler, req)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
