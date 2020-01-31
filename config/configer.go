package config

type Configer interface {
	Init(mode ModeType, path, configName, host, meta string) error
	Register(key string, config interface{})
}
