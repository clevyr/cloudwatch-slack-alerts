package main

import (
	"os"

	"github.com/clevyr/cloudwatch-slack-alerts/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
