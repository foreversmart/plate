package middleware

import (
	"fmt"
	"github.com/foreversmart/plate/common"
	"github.com/foreversmart/plate/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func SetLog() gin.HandlerFunc {
	return Log(commonDecorateLogFunc())
}

func Log(decorateLogFunc func(requestLog *logger.Log, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		sourceIP := getRemoteIP(c)
		method := c.Request.Method
		requestID := c.GetString(common.RequestID)

		requestLog := logger.NewEmptyLogger()

		// 	set level、format and output with logger
		decorateLogFunc(requestLog, c)

		requestLog.LogEntry = requestLog.LogEntry.WithFields(map[string]interface{}{
			"request_id": requestID,
		})

		// Set example variable
		c.Set(common.RequestLogger, requestLog)

		requestLog.WithFieldsNewLog(map[string]interface{}{
			"time_start": start,
			"method":     method,
			"source_ip":  sourceIP,
			"path":       path,
			"timedate":   time.Now(),
		}).Infof("start")

		// before request

		c.Next()

		// after request
		latency := time.Since(start)

		// access the status we are sending
		status := c.Writer.Status()
		if len(c.Errors) > 0 {
			requestLog.LogEntry = requestLog.LogEntry.WithField("error", c.Errors.String())
		}
		requestLog.WithFieldsNewLog(map[string]interface{}{
			"time_latency": fmt.Sprint(latency),
			"latency":      latency.Seconds(),
			"status":       status,
			"method":       method,
			"source_ip":    sourceIP,
			"path":         path,
			"time_finish":  time.Now(),
			"timedate":     time.Now(),
		}).Infof("finish")

	}
}


func commonDecorateLogFunc() func(requestLog *logger.Log, c *gin.Context) {
	return func(requestLog *logger.Log, c *gin.Context) {
		requestLog.SetLogLevel("info")
		requestLog.SetFormat("json")
		requestLog.AddDefaultHook()
	}
}

func getRemoteIP(c *gin.Context) string {
	return c.ClientIP()
}
