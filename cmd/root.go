package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for mgit
var rootCmd = &cobra.Command{
	Use:   "mgit",
	Short: "mgit is a simple version control system",
	Long:  `mgit is a lightweight Git-like version control system built for learning VCS fundamentals.`,
}

// Execute starts the command line tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
