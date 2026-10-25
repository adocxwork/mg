package main

import (
	"embed"
	"mgit/cmd"
	"mgit/pkg/api"
)

// Embed the compiled React frontend into the Go binary
//
//go:embed web/dist
var uiAssets embed.FS

// main is the entry point that initializes the UI assets and runs the CLI commands
func main() {
	api.UIAssets = uiAssets
	cmd.Execute()
}
