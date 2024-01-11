package cmd

import (
	"github.com/bo-er/corrupt-mysql/backend"
	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/spf13/cobra"
)

// highCPUCmd increases mysql's cpu usage dramatically.
var highCPUCmd = &cobra.Command{
	Use:   "hcpu",
	Short: "increasing mysql's cpu usage dramatically",
	Long:  `./corrupt-mysql hcpu -H10.186.62.63 -P25690 -uuniverse_udb -p123`,
	Run: func(cmd *cobra.Command, args []string) {
		connect := pkg.Connect{
			User:     user,
			Host:     host,
			Password: password,
			DBName:   "mysql",
			Port:     port,
		}
		err := backend.MakeCPUUsageHigh(connect)
		if err != nil {
			panic(err.Error())
		}
	},
}
