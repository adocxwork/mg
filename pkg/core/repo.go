package core

import (
	"fmt"
	"os"
	"path/filepath"
)

const MgitDir = ".mgit"

// InitRepo creates .mgit directory with objects and refs subfolders.
func InitRepo() error {
	_ = os.MkdirAll(filepath.Join(MgitDir, "objects"), 0755)
	_ = os.MkdirAll(filepath.Join(MgitDir, "refs", "heads"), 0755)
	_ = os.WriteFile(filepath.Join(MgitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0644)
	fmt.Println("Initialized empty mgit repository in .mgit/")
	return nil
}

// GetRepoPath returns current working directory.
func GetRepoPath() (string, error) {
	return os.Getwd()
}
