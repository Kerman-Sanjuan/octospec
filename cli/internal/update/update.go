// Package update re-applies an install from the saved state, preserving local
// edits to managed files.
package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/plan"
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
	files, err := plan.Files(cfg.Tools, repo, cfg.Models, cfg.Scopes)
	if err != nil {
		return err
	}
	preserved := 0
	next := map[string]string{}
	for _, f := range files {
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
	cfg.Files = next
	if err := config.Save(repo, cfg); err != nil {
		return err
	}
	fmt.Printf("updated %d files (%d preserved)\n", len(files), preserved)
	return nil
}
