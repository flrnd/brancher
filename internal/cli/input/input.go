// Package input
package input

import (
	"bufio"
	"os"
	"strings"

	"github.com/flrnd/brancher/internal/cli/output"
	"github.com/spf13/cobra"
)

var AskFn = func(cmd *cobra.Command, question string) string {
	output.Prompt(cmd, question)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

func Ask(cmd *cobra.Command, question string) string {
	return AskFn(cmd, question)
}
