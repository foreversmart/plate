package application

import (
	"github.com/foreversmart/plate/router"
	"net/http"
)

func (a *Application) Get(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodGet, path, handler, req)
}
func (a *Application) Head(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodHead, path, handler, req)
}
func (a *Application) Post(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodPost, path, handler, req)
}
func (a *Application) Put(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodPut, path, handler, req)
}
func (a *Application) Patch(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodPatch, path, handler, req)
}
func (a *Application) Delete(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodDelete, path, handler, req)
}
func (a *Application) Connect(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodConnect, path, handler, req)
}
func (a *Application) Options(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodOptions, path, handler, req)
}
func (a *Application) Trace(path string, handler router.Handler, req interface{}) {
	a.Handle(http.MethodTrace, path, handler, req)
}
