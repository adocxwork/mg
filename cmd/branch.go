package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mgit/pkg/core"
)

// branchCmd lists all branches or creates a new branch
var branchCmd = &cobra.Command{
	Use:   "branch [branch-name]",
	Short: "List or create branches",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := core.GetRepoPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Fatal:", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Case 1: List all branches
			branches, err := core.ListBranches(repoRoot)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			currentBranch, _ := core.GetCurrentBranch(repoRoot)

			for _, b := range branches {
				if b == currentBranch {
					fmt.Printf("* %s\n", b)
				} else {
					fmt.Printf("  %s\n", b)
				}
			}
		} else {
			// Case 2: Create a new branch pointing to current HEAD
			branchName := args[0]
			headCommit, err := core.GetHEADCommit(repoRoot)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			if headCommit == "" {
				fmt.Fprintln(os.Stderr, "Error: cannot create branch before first commit")
				os.Exit(1)
			}

			if err := core.CreateBranch(repoRoot, branchName, headCommit); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			fmt.Printf("Created branch '%s'\n", branchName)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
