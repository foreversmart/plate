package route

type Server interface {
	// Route return server root route
	Route() Router
	// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
	// Note: this method will block the calling goroutine indefinitely unless an error happens.
	Run(addr ...string)
	// Close the server close, after this method called all request all failed with 500 status code
	Close()
	// Wait all server request handle event finish or timeout
	Wait(timeout int)
}
