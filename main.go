package main

import (
	"os"

	"github.com/bo-er/corrupt-mysql/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetOutput(os.Stdout)
	cmd.Execute()
}
