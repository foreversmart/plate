package application

import (
	"github.com/foreversmart/plate/config"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/router"
	"io"
)

type Application struct {
	config.Configer
	*router.Router
	TraceCloser io.Closer
}

func NewApplication(mode config.ModeType, srcPath string, handle router.Handler) (app *Application, err error) {
	app = &Application{}
	c := config.NewTomlConfig()
	app.Configer = c

	app.Register("server", router.Config)
	app.Register("log", &logger.Config)

	err = c.Init(mode, srcPath, "", "", "")
	if err != nil {
		return nil, err
	}

	// set config

	// set logger , log level and log formatter
	logger.StdLog.SetLogLevel(logger.Config.Level)
	logger.StdLog.SetFormat(logger.Config.Format)
	logger.StdLog.SetOutput(logger.Config.Output)

	logger.StdLog.Debug("std logger setup completed.")

	// route
	app.Router = router.NewRouter(router.Config, handle)
	return
}
