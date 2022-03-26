package util

type Setting struct {
	LogFile string  `mapstructure:"log_file"`
	LogLevel string `mapstructure:"log_level"`
	DriverName string `mapstructure:"driver_name"`
	DatabaseDriverString string `mapstructure:"database_driver_string"`
	MigrateionDirectory string `mapstructure:"migration_directory"`
}

var defaultSetting Setting

func GetSetting() * Setting {
	return &defaultSetting
}
