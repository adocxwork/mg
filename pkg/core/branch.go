package core

import (
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentBranch reads current branch name from .mgit/HEAD.
func GetCurrentBranch(repoRoot string) (string, error) {
	data, _ := os.ReadFile(filepath.Join(repoRoot, MgitDir, "HEAD"))
	content := strings.TrimSpace(string(data))
	return strings.TrimPrefix(content, "ref: refs/heads/"), nil
}

// GetBranchCommit reads the commit hash from .mgit/refs/heads/<branch>.
func GetBranchCommit(repoRoot, branch string) (string, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, MgitDir, "refs", "heads", branch))
	return strings.TrimSpace(string(data)), err
}

// UpdateBranch updates .mgit/refs/heads/<branch> with a new commit hash.
func UpdateBranch(repoRoot, branch, commitHash string) error {
	branchFile := filepath.Join(repoRoot, MgitDir, "refs", "heads", branch)
	return os.WriteFile(branchFile, []byte(commitHash+"\n"), 0644)
}

// ListBranches returns all branch names.
func ListBranches(repoRoot string) ([]string, error) {
	entries, _ := os.ReadDir(filepath.Join(repoRoot, MgitDir, "refs", "heads"))
	var list []string
	for _, e := range entries {
		list = append(list, e.Name())
	}
	return list, nil
}

// CreateBranch creates a new branch pointing to commitHash.
func CreateBranch(repoRoot, branch, commitHash string) error {
	return UpdateBranch(repoRoot, branch, commitHash)
}

// UpdateHEAD points .mgit/HEAD to a new branch.
func UpdateHEAD(repoRoot, branch string) error {
	return os.WriteFile(filepath.Join(repoRoot, MgitDir, "HEAD"), []byte("ref: refs/heads/"+branch+"\n"), 0644)
}

// GetHEADCommit returns the commit hash for the current branch.
func GetHEADCommit(repoRoot string) (string, error) {
	branch, _ := GetCurrentBranch(repoRoot)
	return GetBranchCommit(repoRoot, branch)
}
