package main

import (
	"os"

	"github.com/flrnd/brancher/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
