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
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stkali/garden/util"
	"github.com/stkali/log"
)

const (
	PROGRAM_NAME = "garden"
	CONFIG_TYPE = "yaml"
	CONFIG_NAME = ".garden"
	VERSION = "0.1.0"	
)
var cfgFile string
var setting = util.GetSetting()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   PROGRAM_NAME,
	Short: "go language backend playground",
	Version: VERSION,
	// disable default command eg: completion
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "config"{ return }
	cobra.OnInitialize(initConfig, initLogger)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.garden.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".garden" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType(CONFIG_TYPE)
		viper.SetConfigName(CONFIG_NAME)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "not found config file\nrun 'garden config' have a good idea")
		os.Exit(util.START_FAILED)
	}
	viper.Unmarshal(setting)
}

func initLogger(){
	log.SetLevel(setting.LogLevel)
	var (
		logger io.Writer
		err    error
	)
	if setting.LogFile == "" {
		logger = os.Stdout
	} else {
		logger, err = os.OpenFile(setting.LogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		util.CheckError("failed to open log file", err)
	}
	log.SetOutput(logger)
	info := fmt.Sprintf(infoTemplate, PROGRAM_NAME, VERSION, viper.ConfigFileUsed(), setting.LogFile, setting.LogLevel)
	log.Info("successfully load config and logger\n", info)
}

var infoTemplate = `
-------------- %s: %s --------------
    config file  : %s
    log file     : %s
    log level    : %s
`