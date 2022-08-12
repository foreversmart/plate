package config

import (
	"encoding/json"
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

		s, e := json.Marshal(vv)
		if e != nil {
			panic(e)
		}

		e = json.Unmarshal(s, &v)
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

	filename := "config.ini"
	if configName != "" {
		filename = configName
	}

	switch mode {
	case Development:
		// try application.development.json
		filename = "config.development.ini"

	case Test:
		// try application.test.json
		filename = "config.test.ini"

	case Production:
		// skip

	}

	file := path.Join(srcPath, "config", filename)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		file = path.Join(srcPath, "config", "config.ini")
	}

	return file
}
