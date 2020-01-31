package middleware

import (
	"fmt"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/logger/stack"
	"log"
	"net/http/httputil"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow       = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	// 原生logger
	var l *log.Logger
	l = log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)

	// Discovery logger
	fields := make(map[string]interface{})
	panicLogger := logger.NewEmptyLoggerWithFields(fields).AddDefaultHook()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if l != nil {
					stack := stack.Stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					msg := fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
					l.Println(msg)
					panicLogger.Error(msg)
				}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
