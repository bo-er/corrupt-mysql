package cmd

import (
	"github.com/bo-er/corrupt-mysql/backend"
	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/spf13/cobra"
)

// slowlogCmd creates a slowlog record in a MySQL instance.
var slowlogCmd = &cobra.Command{
	Use:   "slowlog",
	Short: "creating a slowlog record in a MySQL instance.",
	Long:  `./corrupt-mysql slowlog -H10.186.62.63 -P25690 -uuniverse_udb -p123`,
	Run: func(cmd *cobra.Command, args []string) {
		connect := pkg.Connect{
			User:     user,
			Host:     host,
			Password: password,
			DBName:   "mysql",
			Port:     port,
		}
		err := backend.CreateSlowlog(connect)
		if err != nil {
			panic(err.Error())
		}
	},
}
