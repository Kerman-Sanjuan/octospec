// Package update re-applies an install from the saved state, preserving local
// edits to managed files.
package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/cleanup"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/plan"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// Run re-applies the recorded install for repo.
func Run(repo string) error {
	if repo == "" {
		var err error
		if repo, err = os.Getwd(); err != nil {
			return err
		}
	}
	cfg, err := config.Load(repo)
	if err != nil {
		return err
	}
	if len(cfg.Tools) == 0 {
		return errors.New("nothing installed here; run `octospec install` first")
	}
	cfg.Scopes = inferScopes(cfg)
	files, err := plan.Files(cfg.Tools, repo, cfg.Models, cfg.Scopes)
	if err != nil {
		return err
	}
	preserved := 0
	next := map[string]string{}
	keep := make(map[string]bool, len(files))
	for _, f := range files {
		keep[f.Path] = true
		if prev, ok := cfg.Files[f.Path]; ok {
			if cur, err := os.ReadFile(f.Path); err == nil && config.Hash(cur) != prev {
				fmt.Printf("preserved (local edit): %s\n", f.Path)
				next[f.Path] = prev // keep the old hash so the edit stays protected
				preserved++
				continue
			}
		}
		if err := os.MkdirAll(filepath.Dir(f.Path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(f.Path, []byte(f.Content), 0o644); err != nil {
			return err
		}
		next[f.Path] = config.Hash([]byte(f.Content))
	}
	cleanup.Global(&cfg, keep, false)
	cfg.Files = next
	if err := config.Save(repo, cfg); err != nil {
		return err
	}
	fmt.Printf("updated %d files (%d preserved)\n", len(files), preserved)
	return nil
}

// inferScopes returns the recorded scope per tool, and infers a missing one
// from the recorded files: global when a recorded file sits under the tool's
// global directory, repo-local otherwise.
func inferScopes(cfg config.Config) map[string]string {
	scopes := map[string]string{}
	for _, tool := range cfg.Tools {
		if s := cfg.Scopes[tool]; s != "" {
			scopes[tool] = s
			continue
		}
		scope := targets.ScopeLocal
		if t, ok := targets.Get(tool); ok {
			prefix := t.Expand("", targets.ScopeGlobal) + string(filepath.Separator)
			for path := range cfg.Files {
				if strings.HasPrefix(path, prefix) {
					scope = targets.ScopeGlobal
					break
				}
			}
		}
		scopes[tool] = scope
	}
	return scopes
}
