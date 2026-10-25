package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"mgit/pkg/core"
)

var UIAssets fs.FS

type API struct {
	RepoRoot string
}

func StartServer(repoRoot, port string) error {
	a := &API{RepoRoot: repoRoot}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		staged, _, _, _ := core.GetStatus(a.RepoRoot)
		branch, _ := core.GetCurrentBranch(a.RepoRoot)
		head, _ := core.GetHEADCommit(a.RepoRoot)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"staged": staged, "currentBranch": branch, "headCommit": head,
		})
	})

	mux.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Paths []string `json:"paths"` }
		_ = json.NewDecoder(r.Body).Decode(&req)
		for _, p := range req.Paths { _ = core.Add(a.RepoRoot, p) }
	})

	mux.HandleFunc("/api/commit", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Message string `json:"message"` }
		_ = json.NewDecoder(r.Body).Decode(&req)
		head, _ := core.GetHEADCommit(a.RepoRoot)
		parents := []string{}
		if head != "" { parents = append(parents, head) }
		_ = core.CreateCommit(a.RepoRoot, req.Message, parents)
	})

	mux.HandleFunc("/api/log", func(w http.ResponseWriter, r *http.Request) {
		hash, _ := core.GetHEADCommit(a.RepoRoot)
		var list []core.Commit
		for hash != "" {
			c, err := core.ReadCommit(a.RepoRoot, hash)
			if err != nil { break }
			list = append(list, *c)
			if len(c.Parents) > 0 { hash = c.Parents[0] } else { hash = "" }
		}
		json.NewEncoder(w).Encode(list)
	})

	mux.HandleFunc("/api/branch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			b, _ := core.ListBranches(a.RepoRoot)
			curr, _ := core.GetCurrentBranch(a.RepoRoot)
			json.NewEncoder(w).Encode(map[string]interface{}{"branches": b, "current": curr})
		} else {
			var req struct{ Name string `json:"name"` }
			_ = json.NewDecoder(r.Body).Decode(&req)
			head, _ := core.GetHEADCommit(a.RepoRoot)
			_ = core.CreateBranch(a.RepoRoot, req.Name, head)
		}
	})

	mux.HandleFunc("/api/checkout", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Target string `json:"target"` }
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = core.Checkout(a.RepoRoot, req.Target)
	})

	mux.HandleFunc("/api/merge", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ TargetBranch string `json:"targetBranch"` }
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = core.Merge(a.RepoRoot, req.TargetBranch)
	})

	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Mode, Target string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = core.Reset(a.RepoRoot, req.Mode, req.Target)
	})

	if UIAssets != nil {
		if fsys, err := fs.Sub(UIAssets, "web/dist"); err == nil {
			mux.Handle("/", http.FileServer(http.FS(fsys)))
		}
	}

	fmt.Printf("UI server running on http://localhost:%s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
