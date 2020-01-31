package logger

type Requester interface {
	RequestID() string
}
