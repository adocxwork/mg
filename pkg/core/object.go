package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"mgit/pkg/utils"
)

// Commit represents a repository snapshot.
type Commit struct {
	ID        string            `json:"id"`
	Parents   []string          `json:"parents,omitempty"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Files     map[string]string `json:"files"` // filename -> hash
}

// WriteBlob stores file contents in .mgit/objects/<hash>.
func WriteBlob(repoRoot, filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := utils.HashBytes(data)
	_ = os.WriteFile(filepath.Join(repoRoot, MgitDir, "objects", hash), data, 0644)
	return hash, nil
}

// WriteCommit saves a Commit as JSON in .mgit/objects/<hash>.
func WriteCommit(repoRoot string, commit Commit) (string, error) {
	data, _ := json.Marshal(commit)
	hash := utils.HashBytes(data)
	commit.ID = hash

	data, _ = json.Marshal(commit)
	err := os.WriteFile(filepath.Join(repoRoot, MgitDir, "objects", hash), data, 0644)
	return hash, err
}

// ReadCommit loads a Commit from .mgit/objects/<hash>.
func ReadCommit(repoRoot, hash string) (*Commit, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, MgitDir, "objects", hash))
	if err != nil {
		return nil, err
	}
	var commit Commit
	_ = json.Unmarshal(data, &commit)
	return &commit, nil
}
