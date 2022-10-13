package middleware

import (
	"fmt"
	"github.com/foreversmart/plate/logger"
	"time"
)

type LogMidReq struct {
	ReqId     string `json:"req_id" plate:"Req_id,header"`
	Method    string `json:"method" plate:"method,mid"`
	Path      string `json:"path" plate:"path,mid"`
	ParamName string `json:"param_name" plate:"param_name,mid"`
}

type LogMidResp struct {
	Logger       logger.Logger `json:"logger" plate:"logger,mid"`
	ReqRawLog    *logger.Log   `json:"req_raw_log" plate:"req_raw_log,mid:full"`
	LogStartTime time.Time     `json:"log_start_time" plate:"log_start_time,mid:full"`
}

func LogStart(req interface{}) (resp interface{}, err error) {

	arg := req.(*LogMidReq)

	start := time.Now()
	//path := c.Request.URL.Path
	//sourceIP := getRemoteIP(c)
	//method := c.Request.Method
	requestID := arg.ReqId

	requestLog := logger.NewEmptyLogger()

	// TODO	set level、format and output by config file
	requestLog.SetLogLevel("info")
	requestLog.SetFormat("json")
	requestLog.AddDefaultHook()

	requestLog.LogEntry = requestLog.LogEntry.WithFields(map[string]interface{}{
		"request_id": requestID,
		"method":     arg.Method,
		"path":       arg.Path,
		"paramName":  arg.ParamName,
	})

	requestLog.WithFieldsNewLog(map[string]interface{}{
		"time_start": start,
		"timedate":   time.Now(),
	}).Infof("start")

	//fmt.Println(requestLog.LogEntry)

	return &LogMidResp{
		Logger:       requestLog,
		LogStartTime: start,
		ReqRawLog:    requestLog,
	}, nil
}

type LogFinishReq struct {
	ReqRawLog    *logger.Log `json:"req_raw_log" plate:"req_raw_log,mid:full"`
	LogStartTime time.Time   `json:"log_start_time" plate:"log_start_time,mid:full"`
}

func LogFinish(req interface{}) (resp interface{}, err error) {
	arg := req.(*LogFinishReq)
	// after request
	latency := time.Since(arg.LogStartTime)
	// access the status we are sending
	arg.ReqRawLog.WithFieldsNewLog(map[string]interface{}{
		"time_latency": fmt.Sprint(latency),
		"latency":      latency.Seconds(),
		"time_finish":  time.Now(),
		"timedate":     time.Now(),
	}).Infof("finish")
	return
}
