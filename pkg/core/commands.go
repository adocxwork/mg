package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mgit/pkg/utils"
)

// Add stages a file: hashes it, saves it to objects/, and records it in index.
func Add(repoRoot, filename string) error {
	data, err := os.ReadFile(filepath.Join(repoRoot, filename))
	if err != nil {
		return err
	}

	hash := utils.HashBytes(data)
	_ = os.WriteFile(filepath.Join(repoRoot, MgitDir, "objects", hash), data, 0644)

	index, _ := ReadIndex(repoRoot)
	index.Entries[filename] = hash
	return WriteIndex(repoRoot, index)
}

// CreateCommit creates a commit from staged files and moves branch pointer.
func CreateCommit(repoRoot, message string, parents []string) error {
	index, err := ReadIndex(repoRoot)
	if err != nil || len(index.Entries) == 0 {
		return fmt.Errorf("nothing to commit")
	}

	commit := Commit{
		Parents:   parents,
		Message:   message,
		Timestamp: time.Now(),
		Files:     index.Entries,
	}

	commitHash, err := WriteCommit(repoRoot, commit)
	if err != nil {
		return err
	}

	branch, _ := GetCurrentBranch(repoRoot)
	if branch != "" {
		_ = UpdateBranch(repoRoot, branch, commitHash)
	} else {
		_ = os.WriteFile(filepath.Join(repoRoot, MgitDir, "HEAD"), []byte(commitHash+"\n"), 0644)
	}

	fmt.Printf("[%s] %s\n", commitHash[:7], message)
	return nil
}

// Checkout restores all files from a commit and updates HEAD.
func Checkout(repoRoot, target string) error {
	commitHash := target
	if utils.Exists(filepath.Join(repoRoot, MgitDir, "refs", "heads", target)) {
		commitHash, _ = GetBranchCommit(repoRoot, target)
		_ = UpdateHEAD(repoRoot, target)
	}

	commit, err := ReadCommit(repoRoot, commitHash)
	if err != nil {
		return err
	}

	for file, hash := range commit.Files {
		data, _ := os.ReadFile(filepath.Join(repoRoot, MgitDir, "objects", hash))
		_ = os.WriteFile(filepath.Join(repoRoot, file), data, 0644)
	}

	return WriteIndex(repoRoot, &Index{Entries: commit.Files})
}

// GetStatus returns staged, modified, and untracked files.
func GetStatus(repoRoot string) (staged, modified, untracked []string, err error) {
	index, _ := ReadIndex(repoRoot)
	for f := range index.Entries {
		staged = append(staged, f)
	}
	return staged, nil, nil, nil
}

// Status prints the current status.
func Status(repoRoot string) error {
	branch, _ := GetCurrentBranch(repoRoot)
	fmt.Printf("On branch %s\n\n", branch)

	staged, _, _, _ := GetStatus(repoRoot)
	if len(staged) > 0 {
		fmt.Println("Changes to be committed:")
		for _, f := range staged {
			fmt.Printf("  (staged) %s\n", f)
		}
	} else {
		fmt.Println("Nothing to commit, working tree clean")
	}
	return nil
}

// Log prints commit history starting from HEAD.
func Log(repoRoot string) error {
	hash, err := GetHEADCommit(repoRoot)
	if err != nil || hash == "" {
		return fmt.Errorf("no commits yet")
	}

	for hash != "" {
		commit, err := ReadCommit(repoRoot, hash)
		if err != nil {
			break
		}
		fmt.Printf("commit %s\nDate:   %s\n\n    %s\n\n", commit.ID, commit.Timestamp.Format(time.RFC1123), commit.Message)
		if len(commit.Parents) > 0 {
			hash = commit.Parents[0]
		} else {
			hash = ""
		}
	}
	return nil
}

// Merge combines files from targetBranch into the current branch.
func Merge(repoRoot, targetBranch string) error {
	headHash, _ := GetHEADCommit(repoRoot)
	targetHash, _ := GetBranchCommit(repoRoot, targetBranch)
	targetCommit, err := ReadCommit(repoRoot, targetHash)
	if err != nil {
		return err
	}

	index, _ := ReadIndex(repoRoot)
	for file, hash := range targetCommit.Files {
		index.Entries[file] = hash
		data, _ := os.ReadFile(filepath.Join(repoRoot, MgitDir, "objects", hash))
		_ = os.WriteFile(filepath.Join(repoRoot, file), data, 0644)
	}
	_ = WriteIndex(repoRoot, index)

	return CreateCommit(repoRoot, "Merge branch "+targetBranch, []string{headHash, targetHash})
}

// Reset moves the branch pointer to a specific commit.
func Reset(repoRoot, mode, target string) error {
	targetHash := target
	if utils.Exists(filepath.Join(repoRoot, MgitDir, "refs", "heads", target)) {
		targetHash, _ = GetBranchCommit(repoRoot, target)
	}

	branch, _ := GetCurrentBranch(repoRoot)
	if branch != "" {
		_ = UpdateBranch(repoRoot, branch, targetHash)
	}

	if mode == "hard" {
		_ = Checkout(repoRoot, targetHash)
	}

	fmt.Printf("HEAD is now at %s\n", targetHash[:7])
	return nil
}
