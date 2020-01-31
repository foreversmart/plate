package application

import (
	"github.com/foreversmart/plate/config"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/router"
	"io"
)

type Application struct {
	Config      config.Configer
	Router      *router.Router
	TraceCloser io.Closer
}

func NewApplication(mode config.ModeType, srcPath string, handle router.Handler) (app *Application, err error) {
	app = &Application{}
	c := config.NewTomlConfig()
	c.Register("server", &router.Config)
	c.Register("log", &logger.Config)
	err = c.Init(mode, srcPath, "", "", "")
	if err != nil {
		return nil, err
	}

	// set config
	app.Config = c

	// set logger , log level and log formatter
	logger.StdLog.SetLogLevel(logger.Config.Level)
	logger.StdLog.SetFormat(logger.Config.Format)
	logger.StdLog.SetOutput(logger.Config.Output)

	logger.StdLog.Debug("Log setup completed.")

	// route
	app.Router = router.NewRouter(router.Config, handle)
	return
}

func (app *Application) Run() {
	if app == nil || app.Router == nil {
		return
	}

	app.Router.Run()
}
