package ginroute

import (
	"encoding/json"
	"fmt"
	"github.com/foreversmart/plate/errors"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/logger/stack"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func GinRecovery(recV interface{}, request *http.Request, args []interface{}) (resp interface{}, err error) {
	stack := stack.Stack(4)

	// 原生logger
	var l *log.Logger
	l = log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)

	httprequest, _ := httputil.DumpRequest(request, false)

	argStr := "[Recovery] args: \n"
	for i, arg := range args {
		s, _ := json.Marshal(arg)
		argStr = argStr + fmt.Sprintf("%d :%s \n", i, string(s))
	}

	msg := fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n%s%s", recV, string(httprequest), stack, argStr)
	logger.StdLog.Error(msg)
	l.Println(msg)

	return nil, errors.ErrorUnknownError
}
