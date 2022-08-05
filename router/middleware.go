package router

// Middleware middleware is pre handler for common logic for concrete logic
// H is the middleware handler, handler's resp tag `plate:"x,mid"`
// means resp will pass this value to parser's mid view
// V is Handler request concrete type
type Middleware struct {
	H Handler
	V interface{}
}
