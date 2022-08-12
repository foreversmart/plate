package route

type Router interface {
	// Handle register a new handle logic for given path, method and object v
	// v must be a struct or struct pointer
	Handle(method, path string, handler Handler, v interface{})
	// AddMiddle will add mid handler before handle's handler called, v is mid request
	// v must be a struct or struct pointer
	AddMiddle(mid Handler, v interface{})

	// SetRecovery is
	SetRecovery(rec Handler)
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
