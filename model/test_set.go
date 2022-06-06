package model

import (
	"github.com/foreversmart/plate/config"
	"math/rand"
	"os"
	"testing"
	"time"
)

// SetTest 测试设置,一般只能在同服务内使用，因为同一服务的Config解析相同
// 对于在TestMain中需要额外设置client的则不能调用该方法
func SetTest(m *testing.M, srcPath string) {
	rand.Seed(time.Now().Unix())

	var (
		runMode = config.Test
	)

	config.NewTomlConfig()
	conf, err := config.NewConfig(runMode, srcPath)
	if err != nil {
		panic(err)
	}

	config.SetupConfig(configig)

	code := m.Run()
	MongoModel().Session().DB(MongoModel().Database()).DropDatabase()

	os.Exit(code)
}
