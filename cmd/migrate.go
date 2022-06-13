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
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/stkali/garden/util"
)

var step int

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database changed(init)",
	PreRun: func(cmd *cobra.Command, args []string) {
		//log.SetLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(
			setting.MigrationDirectory,
			setting.DatabaseDriverString,
		)
		util.CheckError("failed to connect database", err)
		// migrate all revision
		if step == 0 {
			err = m.Up()
		// migrate specify revision (up or down)
		} else {
			err = m.Steps(step)
		}
		switch err {
		case migrate.ErrNoChange:
			fmt.Println("warning: migrate does nothing")
		case os.ErrNotExist:
			fmt.Printf("warning: specify step overflow: %d\n", step)
		case nil:
			fmt.Println("successfully to migrate database")
		default:
			util.CheckError("failed to migrate database changed", err)	
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().IntVarP(&step, "step", "s", 0, "")
}
