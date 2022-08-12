package mongo

import (
	"github.com/foreversmart/plate/config"
	"github.com/foreversmart/plate/logger"
	"testing"
)

var (
	testDb *Model
)

func TestMain(m *testing.M) {

	//var (
	//	currentPath, _ = os.Getwd()
	//	srcPath        string
	//)
	//
	//srcPaths := strings.Split(currentPath, "model")
	//if len(srcPaths) > 0 {
	//	srcPath = srcPaths[0]
	//}

	var modelConfig *Config

	conf := config.NewTomlConfig()
	conf.Register("mongo", &modelConfig)
	conf.Init(config.Test, "", "", "", "")

	testDb = NewModel(modelConfig, logger.StdLog)
}
