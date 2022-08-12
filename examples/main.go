package main

import (
	"github.com/foreversmart/plate/application"
	"github.com/foreversmart/plate/config"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	app, err := application.NewApplication(config.Development, pwd)
	if err != nil {
		panic(err)
	}

	type Req struct {
	}

	type Resp struct {
		Res string `json:"res"`
	}
	app.Route().Get("ping", func(r interface{}) (resp interface{}, err error) {
		return &Resp{
			Res: "pong",
		}, nil
	}, Req{})

	app.Run()
}
