package application

import (
	"github.com/foreversmart/plate/config"
	"github.com/foreversmart/plate/logger"
	"github.com/foreversmart/plate/router"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	config.Configer
	router.V2Router
	TraceCloser io.Closer
}

func NewApplication(mode config.ModeType, srcPath string, r router.V2Router) (app *Application, err error) {
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
	app.V2Router = r

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

		for {
			select {
			case <-sigs:
				logger.StdLog.Info("接受到了结束进程的信号")

				time.Sleep(5 * time.Second)

				// close trace instance

				os.Exit(0)
			}
		}

	}()
	return
}
