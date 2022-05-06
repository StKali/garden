/*
Copyright © 2022 <st·kali clarkmonkey@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stkali/garden/util"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			cfgFile = fmt.Sprintf("%s.%s", CONFIG_NAME, CONFIG_TYPE)
		}
		err := CreateConfigFile(cfgFile)
		util.CheckError("failed to create default config", err)
		fmt.Println("successfully create default config at:", cfgFile)
	},
}

func CreateConfigFile(config string) error {

	viper.SetConfigType(CONFIG_TYPE)
	viper.SetDefault("log_file", "")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("driver_name", "postgres")
	viper.SetDefault("database_driver_string", "postgresql://root:password@localhost:5432/garden?sslmode=disable")
	viper.SetDefault("migration_directory", "file://db/migration")
	viper.SetDefault("gin_server_address", "0.0.0.0:8000")
	viper.SetDefault("token_type", "paseto")
	viper.SetDefault("token_duration", time.Hour*24)
	viper.SetDefault("refresh_token_duration", time.Hour*24*14)
	return viper.WriteConfigAs(config)
}

func init() {
	rootCmd.AddCommand(configCmd)
}
