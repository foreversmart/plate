package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToml(t *testing.T) {
	type TestConfig struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		UserName string `toml:"userName"`
		Password string `toml:"password"`
	}

	var testConfig *TestConfig
	c := NewTomlConfig()
	c.Register("test", &testConfig)
	err := c.Init(Test, "", "", "", "")
	assert.Nil(t, err)

	fmt.Println(testConfig, "0000")

}
