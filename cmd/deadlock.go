package cmd

import (
	"github.com/bo-er/corrupt-mysql/backend"
	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/spf13/cobra"
)

// deadlockCmd creates a deadlock in a MySQL instance.
var deadlockCmd = &cobra.Command{
	Use:   "deadlock",
	Short: "creating a deadlock in a MySQL instance.",
	Long:  `./corrupt-mysql deadlock -H10.186.62.63 -P25690 -uuniverse_udb -p123`,
	Run: func(cmd *cobra.Command, args []string) {
		connect := pkg.Connect{
			User:     user,
			Host:     host,
			Password: password,
			DBName:   "mysql",
			Port:     port,
		}
		err := backend.CreateDeadlock(connect)
		if err != nil {
			panic(err.Error())
		}
	},
}
