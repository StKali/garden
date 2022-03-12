package util

type Setting struct {
	LogFile  string `mapstructure:"log_file"`
	LogLevel string `mapstructure:"log_level"`
}

var defaultSetting Setting

func GetSetting() *Setting {
	return &defaultSetting
}
