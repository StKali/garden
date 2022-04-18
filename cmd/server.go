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
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	"github.com/stkali/garden/api"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"github.com/stkali/log"
	"os"
	"os/signal"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start garden server",
	Run: func(cmd *cobra.Command, args []string) {
		// create database connect
		conn, err := sql.Open(setting.DriverName, setting.DatabaseDriverString)
		log.Infof("database driver: %s, source name: %s", setting.DriverName, setting.DatabaseDriverString)
		util.CheckError("cannot open database", err)

		// create store
		store := db.NewStore(conn)
		log.Infof("successfully created store instance.")

		// create token maker
		maker, err := token.NewMaker(token.GenerateSymmetricKey(), setting.TokenType)
		util.CheckError("failed to create token maker", err)
		log.Infof("successfully created %s token token", setting.TokenType)

		// launcher func of gin server
		ginServer := func(address string) {
			server := api.NewServer(store, maker)
			server.Start(address)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
		ActiveServer("gin-http", setting.GinServerAddress, ginServer, cancel)
		<-ctx.Done()
		log.Infof("garden stop!")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

type LaunchType func(address string)

func ActiveServer(name string, address string, f LaunchType, cancel context.CancelFunc) {
	go func() {
		defer func() {
			exc := recover()
			if exc == nil {
				log.Infof("%s server has stopped", name)
				cancel()
			} else {
				log.Errorf("%s server stop, err: %+s", name, exc)
			}
			cancel()
		}()
		f(address)
	}()
	log.Infof("successfully active %s server on: %s", name, address)
}
