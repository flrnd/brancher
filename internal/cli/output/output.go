// Package output
package output

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func println(cmd *cobra.Command, format string, args ...any) {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), format+"\n", args...); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "output error: %v\n", err)
		os.Exit(1)
	}
}

func Task(cmd *cobra.Command, id, title string) {
	println(cmd, "%s  %s", id, title)
}

func BranchCreated(cmd *cobra.Command, name string) {
	println(cmd, "Created branch: %s", name)
}
