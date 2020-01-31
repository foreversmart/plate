package model

import (
	"os"
	"server/utils"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	var (
		currentPath, _ = os.Getwd()
		srcPath        string
	)

	srcPaths := strings.Split(currentPath, "model")
	if len(srcPaths) > 0 {
		srcPath = srcPaths[0]
	}

	utils.StdLog = utils.NewEmptyLogger()

	SetTest(m, srcPath)
}
