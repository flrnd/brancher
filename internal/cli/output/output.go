// Package output
package output

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func write(cmd *cobra.Command, format string, args ...any) {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), format, args...); err != nil {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "output error:", err)
		os.Exit(1)
	}
}

func Println(cmd *cobra.Command, format string, args ...any) {
	write(cmd, format+"\n", args...)
}

func Prompt(cmd *cobra.Command, text string) {
	write(cmd, "%s ", text)
}

func Task(cmd *cobra.Command, id, title string) {
	write(cmd, "#%s - %s\n", id, title)
}

func BranchCreated(cmd *cobra.Command, name string) {
	write(cmd, "Created branch: %s", name)
}
