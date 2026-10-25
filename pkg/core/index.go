package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Index holds staged files (filename -> hash).
type Index struct {
	Entries map[string]string `json:"entries"`
}

// ReadIndex loads the index from .mgit/index.
func ReadIndex(repoRoot string) (*Index, error) {
	index := &Index{Entries: make(map[string]string)}
	data, err := os.ReadFile(filepath.Join(repoRoot, MgitDir, "index"))
	if err == nil {
		_ = json.Unmarshal(data, index)
	}
	return index, nil
}

// WriteIndex saves the index to .mgit/index.
func WriteIndex(repoRoot string, index *Index) error {
	data, _ := json.Marshal(index)
	return os.WriteFile(filepath.Join(repoRoot, MgitDir, "index"), data, 0644)
}
