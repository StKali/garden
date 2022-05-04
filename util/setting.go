package util

import (
	"time"
)

type Setting struct {
	LogFile              string `mapstructure:"log_file"`
	LogLevel             string `mapstructure:"log_level"`
	DriverName           string `mapstructure:"driver_name"`
	DatabaseDriverString string `mapstructure:"database_driver_string"`
	MigrationDirectory   string `mapstructure:"migration_directory"`
	GinServerAddress     string `mapstructure:"gin_server_address"`
	TokenType            string `mapstructure:"token_type"`
	TokenDuration        time.Duration `mapstructure:"token_duration`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

var defaultSetting Setting

func GetSetting() *Setting {
	return &defaultSetting
}
