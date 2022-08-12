package model

import (
	"os"
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

	SetTest(m, srcPath)
}
