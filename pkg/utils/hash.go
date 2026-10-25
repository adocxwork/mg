package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
)

// HashBytes computes the SHA-1 hash of raw bytes.
func HashBytes(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashFile computes the SHA-1 hash of a file's content.
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return HashBytes(data), nil
}
