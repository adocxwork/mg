package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mgit/pkg/core"
)

var hard bool
var soft bool
var mixed bool

// resetCmd moves the branch pointer backwards to undo commits
// Modes:
// --soft : only moves branch pointer
// --mixed: moves branch pointer and resets staging area (default)
// --hard : moves branch pointer, resets staging area, and discards file changes
var resetCmd = &cobra.Command{
	Use:   "reset <commit-or-branch>",
	Short: "Reset current HEAD to a specified commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := core.GetRepoPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Fatal:", err)
			os.Exit(1)
		}

		target := args[0]
		mode := "mixed"
		if hard {
			mode = "hard"
		} else if soft {
			mode = "soft"
		} else if mixed {
			mode = "mixed"
		}

		if err := core.Reset(repoRoot, mode, target); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	resetCmd.Flags().BoolVar(&hard, "hard", false, "Reset staging area and working files (discards changes)")
	resetCmd.Flags().BoolVar(&soft, "soft", false, "Only move HEAD pointer (keep files and staging area)")
	resetCmd.Flags().BoolVar(&mixed, "mixed", false, "Reset staging area but keep working files (default)")
	rootCmd.AddCommand(resetCmd)
}
