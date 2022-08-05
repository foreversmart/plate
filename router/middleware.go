package router

type Middleware struct {
	H Handler
	V interface{}
}
