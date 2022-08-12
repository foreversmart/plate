package mongo

// Config db配置
type Config struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Passwd   string `toml:"passwd"`
	Database string `toml:"database"`
	Mode     string `toml:"mode"`
	Pool     int    `toml:"pool"`
	Timeout  int    `toml:"timeout"`
}

// Copy 拷贝配置
func (c *Config) Copy() *Config {
	config := *c
	return &config
}
