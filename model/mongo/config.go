package config

// LogExtrConfig 日志属性


// ModelConfig db配置
type ModelConfig struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Passwd   string `toml:"passwd"`
	Database string `toml:"database"`
	Mode     string `toml:"mode"`
	Pool     int    `toml:"pool"`
	Timeout  int    `toml:"timeout"`
}

// Copy 拷贝配置
func (c *ModelConfig) Copy() *ModelConfig {
	config := *c
	return &config
}
