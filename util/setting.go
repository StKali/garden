package util

type Setting struct {
	LogFile string  `mapstructure:"log_file"`
	LogLevel string `mapstructure:"log_level"`
	DriverName string `mapstructure:"driver_name"`
	DriverSourceName string `mapstructure:"driver_source_name"`
}

var defaultSetting Setting

func GetSetting() * Setting {
	return &defaultSetting
}
