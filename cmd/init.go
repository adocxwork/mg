package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mgit/pkg/core"
)

// initCmd creates a new empty mgit repository in the current folder
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new mgit repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.InitRepo(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
