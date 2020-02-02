package main

import (
	"github.com/foreversmart/plate/application"
	"github.com/foreversmart/plate/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	app, err := application.NewApplication(config.Development, pwd, func(e *gin.Engine) *gin.Engine {
		e.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

		return e
	})

	if err != nil {
		panic(err)
	}

	app.Run()
}
