package cmd

import (
	"github.com/bo-er/corrupt-mysql/backend"
	"github.com/bo-er/corrupt-mysql/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// bigDDLCmd creates a big DDL in a MySQL instance.
var bigDDLCmd = &cobra.Command{
	Use:   "bigddl",
	Short: "creating a big DDL in a MySQL instance.",
	Long:  `./corrupt-mysql bigddl -H10.186.62.63 -P25690 -uuniverse_udb -p123`,
	Run: func(cmd *cobra.Command, args []string) {
		connect := pkg.Connect{
			User:     user,
			Host:     host,
			Password: password,
			DBName:   "mysql",
			Port:     port,
		}
		err := backend.CreateBigDDL(connect)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("operation is done")
	},
}
