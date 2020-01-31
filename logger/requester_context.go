package logger

import (
	"errors"
	"github.com/foreversmart/plate/middleware"

	"github.com/gin-gonic/gin"
)

// Logger function
func LoggerFromContext(c *gin.Context) Logger {
	logger, ok := c.MustGet(middleware.RequestLoggerKey).(*Log)
	if !ok {
		panic(errors.New("can't fetch logger from context"))
	}

	return logger
}
