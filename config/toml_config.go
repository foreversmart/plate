package config

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path"
)

// Config 公共配置
type TomlConfig struct {
	C map[string]interface{}
	t map[string]interface{}
}

func NewTomlConfig() *TomlConfig {
	return &TomlConfig{
		C: make(map[string]interface{}),
		t: make(map[string]interface{}),
	}
}

func (c *TomlConfig) Init(mode ModeType, path, configName, host, meta string) error {
	content, err := c.configContent(mode, path, configName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = toml.Decode(string(content), &c.t)
	if err != nil {
		return err
	}

	for k, v := range c.C {
		vv, ok := c.t[k]
		if !ok {
			continue
		}

		buff := &bytes.Buffer{}
		encoder := toml.NewEncoder(buff)
		e := encoder.Encode(vv)
		if e != nil {
			panic(e)
		}

		e = toml.Unmarshal(buff.Bytes(), &v)
		if e != nil {
			panic(e)
		}
	}

	return err
}

func (c *TomlConfig) Register(key string, conf interface{}) {
	c.C[key] = conf
}

// ConfigContent 读取配置文件
func (c *TomlConfig) configContent(mode ModeType, srcPath, configName string) (content []byte, err error) {
	filename := c.findModeConfigFilePath(mode, srcPath, configName)
	content, err = ioutil.ReadFile(filename)
	return
}

// FindModeConfigFilePath 确定配置文件
func (c *TomlConfig) findModeConfigFilePath(mode ModeType, srcPath, configName string) string {
	// adjust srcPath
	srcPath = path.Clean(srcPath)

	filename := "config"
	if configName != "" {
		filename = configName
	}

	switch mode {
	case Development:
		// try application.development.json
		filename = filename + ".development"

	case Test:
		// try application.test.json
		filename = filename + ".test"

	case Production:
		// skip
	}

	filename = filename + ".ini"

	file := path.Join(srcPath, "config", filename)

	// if mode config file not exist fail back to read standard config file
	if _, err := os.Stat(file); os.IsNotExist(err) {
		file = path.Join(srcPath, "config", "config.ini")
	}

	return file
}
