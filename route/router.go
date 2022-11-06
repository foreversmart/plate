package route

import "net/http"

type Router interface {
	// Handle register a new handle logic for given path, method and object v
	// v must be a struct or struct pointer
	Handle(method, path string, handler Handler, v interface{})
	// Sub return the relative sub router, sub router will inherit parent router's middleware
	Sub(relativePath string) Router
	// AddMiddleBefore will add mid handler before handle's handler called, v is mid request
	// v must be a struct or struct pointer
	// sub before middle exec after parent before middles
	AddMiddleBefore(mid Handler, v interface{})
	// AddMiddleAfter will add mid handler before handle's handler called, v is mid request
	// v must be a struct or struct pointer
	// sub after middle exec before parent after middles
	AddMiddleAfter(mid Handler, v interface{})
	// SetRecover set router default recover
	SetRecover(rec Recover)
	// SetAutoCors set router auto add cors methods
	SetAutoCors(bool)
	// Wait will wait all the connection logic close or timeout
	Wait(timeout int)
	// Close will close the router later connection will get failed
	Close()

	Get(path string, handler Handler, req interface{})
	Head(path string, handler Handler, req interface{})
	Post(path string, handler Handler, req interface{})
	Put(path string, handler Handler, req interface{})
	Patch(path string, handler Handler, req interface{})
	Delete(path string, handler Handler, req interface{})
	Connect(path string, handler Handler, req interface{})
	Options(path string, handler Handler, req interface{})
	Trace(path string, handler Handler, req interface{})
}

type Handler func(req interface{}) (resp interface{}, err error)

// Recover is function when handle panic will recover the handle controller function
// recV is recover() return value which catch from panic
// req means handle http request
// args means handle args and the middleware args and args are array by called sequence
type Recover func(recV interface{}, req *http.Request, args []interface{}) (resp interface{}, err error)
