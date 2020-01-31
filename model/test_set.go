package model

import (
	"math/rand"
	"os"
	"server/application/conf"
	"testing"
	"time"
)

// SetTest 测试设置,一般只能在同服务内使用，因为同一服务的Config解析相同
// 对于在TestMain中需要额外设置client的则不能调用该方法
func SetTest(m *testing.M, srcPath string) {
	rand.Seed(time.Now().Unix())

	var (
		runMode = conf.Test
	)

	config, err := conf.NewConfig(runMode, srcPath)
	if err != nil {
		panic(err)
	}

	conf.SetupConfig(config)

	code := m.Run()
	MongoModel().Session().DB(MongoModel().Database()).DropDatabase()

	os.Exit(code)
}
