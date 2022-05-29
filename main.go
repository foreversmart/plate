//package main
//
//import (
//	"flag"
//	"os"
//	"os/signal"
//	"path"
//	"server/application/conf"
//	"syscall"
//	"time"
//
//	"server/application"
//	"server/utils"
//)
//
//var (
//	runMode string // app run mode, available values are [development|test|production], default to development
//	srcPath string // app source path
//)
//
//func main() {
//	flag.StringVar(&runMode, "runMode", "development", "app -runMode=[development|test|production]")
//	flag.StringVar(&srcPath, "srcPath", "", "app -srcPath=/path/to/source")
//
//	flag.Parse()
//
//	mode := conf.ModeType(runMode)
//
//	// verify run mode
//	if !mode.IsValid() {
//		flag.PrintDefaults()
//		return
//	}
//
//	// adjust src path
//	if srcPath == "" {
//		var err error
//
//		srcPath, err = os.Getwd()
//		if err != nil {
//			panic(err)
//		}
//	} else {
//		srcPath = path.Clean(srcPath)
//	}
//
//	utils.StdLog.Infof("启动进程 , runMode: %s , srcPath: %s", runMode, srcPath)
//
//	// init application
//	_, err := application.NewApplication(mode, srcPath)
//	if err != nil {
//		panic(err)
//	}
//
//	go func() {
//		sigs := make(chan os.Signal, 1)
//		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
//
//		for {
//			select {
//			case <-sigs:
//				utils.StdLog.Info("接受到了结束进程的信号")
//
//				time.Sleep(5 * time.Second)
//
//				// close trace instance
//
//				os.Exit(0)
//			}
//		}
//
//	}()
//
//}
