package util

type Setting struct {
	LogFile              string `mapstructure:"log_file"`
	LogLevel             string `mapstructure:"log_level"`
	DriverName           string `mapstructure:"driver_name"`
	DatabaseDriverString string `mapstructure:"database_driver_string"`
	MigrationDirectory   string `mapstructure:"migration_directory"`
	GinServerAddress     string `mapstructure:"gin_server_address"`
	TokenType            string `mapstructure:"token_type"`
}

var defaultSetting Setting

func GetSetting() *Setting {
	return &defaultSetting
}
