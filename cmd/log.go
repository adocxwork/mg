package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mgit/pkg/core"
)

// logCmd displays the commit history starting from HEAD
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := core.GetRepoPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Fatal:", err)
			os.Exit(1)
		}

		if err := core.Log(repoRoot); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
