package main

import (
	"os"

	"github.com/flrnd/brancher/internal/cli"
	_ "github.com/flrnd/brancher/internal/provider/github"
)

func main() {
	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
