package logger

import (
	"errors"
	"github.com/foreversmart/plate/common"
	"github.com/gin-gonic/gin"
)

type Requester interface {
	RequestID() string
}

// Logger function
func LoggerFromContext(c *gin.Context) Logger {
	l, ok := c.MustGet(common.RequestLogger).(*Log)
	if !ok {
		panic(errors.New("can't fetch logger from context"))
	}

	return l
}
