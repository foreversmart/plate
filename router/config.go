package router

var Config = &ConfigType{}

// ServerConfig server配置
type ConfigType struct {
	Host string `toml:"host"`
	Url  string `toml:"url"`
	Port int    `toml:"port"`
}
