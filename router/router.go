package router

type Router interface {
	// Handle register a new handle logic for given path, method and object v
	// v must be a struct or struct pointer
	Handle(method, path string, handler Handler, v interface{})
	Middle(path string, handler Handler, v interface{})
	// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
	// Note: this method will block the calling goroutine indefinitely unless an error happens.
	Run(addr ...string)
	// Wait will wait all the connection logic close or timeout
	Wait(timeout int)
	// Close will close the router later connection will get failed
	Close()
}

type Handler func(req interface{}) (resp interface{}, err error)

func mid(req interface{}) (resp interface{}, err error) {

}
