package cmd

import (
	"github.com/bo-er/corrupt-mysql/backend"
	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/spf13/cobra"
)

// bigTransaction creates big transactions in a mysql instance.
var bigTransaction = &cobra.Command{
	Use:   "bt",
	Short: "creating big transactions in your database",
	Long:  `./corrupt-mysql bt -H10.186.62.63 -P25690 -uuniverse_udb -p123`,
	Run: func(cmd *cobra.Command, args []string) {
		connect := pkg.Connect{
			User:     user,
			Host:     host,
			Password: password,
			DBName:   "mysql",
			Port:     port,
		}
		if len(args) == 0 || args[0] == "" {
			panic("please enter a big transaction size(must be greater than 10kb)")
		}

		err := backend.CreatesBigTransactions(connect, args[0])
		if err != nil {
			panic(err.Error())
		}
	},
}
