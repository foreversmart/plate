package main

import (
	"github.com/foreversmart/plate/application"
	"github.com/foreversmart/plate/config"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/middleware"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	app, err := application.NewApplication(config.Development, pwd)
	if err != nil {
		panic(err)
	}

	type Req struct {
		Logger logger.Logger `json:"logger" plate:"logger,mid"`
	}

	type Resp struct {
		Res string `json:"res"`
	}

	app.Route().AddMiddleBefore(middleware.LogStart, &middleware.LogMidReq{})
	app.Route().AddMiddleAfter(middleware.LogFinish, &middleware.LogFinishReq{})

	app.Route().Get("ping", func(r interface{}) (resp interface{}, err error) {
		arg := r.(*Req)
		arg.Logger.Info("hello")

		return &Resp{
			Res: "pong",
		}, nil
	}, Req{})

	app.Run()

}
