package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	host     string
	user     string
	port     int
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  `corrupt mysql for checking purposes.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deadlockCmd)
	rootCmd.AddCommand(slowlogCmd)
	rootCmd.AddCommand(bigTransaction)

	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "host of the mysql server")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 3306, "port of the mysql server")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "user of the mysql server")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "123", "password of the mysql server")
}
