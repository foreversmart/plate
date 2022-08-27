package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	WithHook bool
	LogHook  Hook
)

var (
	StdLog *Log
)

func init() {
	// set stdlog and set level
	StdLog = NewEmptyLogger()
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	ReqId() string
}

type Log struct {
	mu       sync.Mutex
	out      io.Writer
	LogEntry *log.Entry `json:"log_entry" plate:"log_entry,mid"`
	fields   map[string]interface{}
}

func (logger *Log) WithFields(fields map[string]interface{}) {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}

	logger.LogEntry = logger.LogEntry.WithFields(f)
}

func (logger *Log) WithFieldsNewLog(fields map[string]interface{}) (entry *log.Entry) {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}

	return logger.LogEntry.WithFields(f)
}

func (logger *Log) SetLevel(level log.Level) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.LogEntry.Logger.SetLevel(level)
}

// SetOutput sets the standard logger output.
func (logger *Log) SetOutputWriter(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	logger.LogEntry.Logger.Out = out
}

func NewLog(entry *log.Entry) *Log {
	l := &Log{
		LogEntry: entry,
		fields:   make(map[string]interface{}),
	}
	l.SetFormat("json")
	return l
}

func NewEmptyLogger() *Log {
	logger := log.New()
	l := &Log{
		LogEntry: log.NewEntry(logger),
		fields:   make(map[string]interface{}),
	}
	l.SetFormat("json")
	return l
}

type LogFields map[string]interface{}

func NewEmptyLoggerWithFields(fields map[string]interface{}) *Log {
	l := NewEmptyLogger()
	l.fields = fields
	l.SetFormat("json")
	return l
}

func DecorateLog(logger *log.Entry) *log.Entry {
	var (
		fileName, funcName string
	)
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		fileName = "???"
		funcName = "???"
		line = 0
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		fileSlice := strings.Split(file, path.Dir(funcName))
		fileName = filepath.Join(path.Dir(funcName), fileSlice[len(fileSlice)-1]) + ":" + strconv.Itoa(line)
	}

	return logger.WithField("file", fileName).WithField("func", funcName)

}

// hook
type Hook interface {
	Levels() []log.Level
	Fire(*log.Entry) error
}

func (logger *Log) AddHook(hook Hook) *Log {
	if WithHook && hook != nil {
		logger.LogEntry.Logger.AddHook(hook)
		return logger
	}
	return logger
}

func (logger *Log) AddDefaultHook() *Log {
	if WithHook && LogHook != nil {
		logger.LogEntry.Logger.AddHook(LogHook)
		return logger
	}
	return logger
}

func (logger *Log) SetLogLevel(level string) {
	switch level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}

func (logger *Log) SetOutput(output string) {
	switch output {
	case "stderr":
		logger.SetOutputWriter(os.Stderr)
	case "stdout":
		logger.SetOutputWriter(os.Stdout)
	case "null":
		logger.SetOutputWriter(ioutil.Discard)
	default:
		logger.SetOutputWriter(os.Stderr)
	}
}

func (logger *Log) SetFormat(format string) {
	switch format {
	case "json":
		logger.LogEntry.Logger.Formatter = &log.JSONFormatter{
			DisableTimestamp: true,
		}
	case "text":
		logger.LogEntry.Logger.Formatter = &log.TextFormatter{
			DisableTimestamp: true,
		}
	default:
		logger.LogEntry.Logger.Formatter = &log.JSONFormatter{
			DisableTimestamp: true,
		}
	}
}

func (logger *Log) Debug(args ...interface{}) {
	initLog(logger).Debug(args...)
}

func (logger *Log) Info(args ...interface{}) {
	initLog(logger).Info(args...)
}

func (logger *Log) Warn(args ...interface{}) {
	initLog(logger).Warn(args...)
}

func (logger *Log) Error(args ...interface{}) {
	initLog(logger).Error(args...)
}

func (logger *Log) Fatal(args ...interface{}) {
	initLog(logger).Fatal(args...)
}

func (logger *Log) Panic(args ...interface{}) {
	initLog(logger).Panic(args...)
}

// Entry Printf family functions
func (logger *Log) Debugf(format string, args ...interface{}) {
	initLog(logger).Debugf(format, args...)
}

func (logger *Log) Infof(format string, args ...interface{}) {
	initLog(logger).Infof(format, args...)
}

func (logger *Log) Warnf(format string, args ...interface{}) {
	initLog(logger).Warnf(format, args...)
}

func (logger *Log) Errorf(format string, args ...interface{}) {
	initLog(logger).Errorf(format, args...)
}

func (logger *Log) Fatalf(format string, args ...interface{}) {
	initLog(logger).Fatalf(format, args...)
}

func (logger *Log) Panicf(format string, args ...interface{}) {
	initLog(logger).Panicf(format, args...)
}

func (logger *Log) ReqId() string {
	reqId, ok := logger.LogEntry.Data["request_id"]
	if ok {
		return fmt.Sprintf("%v", reqId)
	}
	return ""
}

func initLog(logger *Log) *log.Entry {
	entry := logger.LogEntry
	f := logger.fields
	return DecorateLog(entry.WithFields(log.Fields(f)).WithFields(log.Fields{"timedate": time.Now()}))
}
