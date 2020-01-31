package logger

var Config *ConfigType

type ConfigType struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
	Output string `toml:"output"`
}
