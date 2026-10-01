// Package uninstall removes the files octospec manages, using the recorded
// state, and preserves hand edits.
package uninstall

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kerman-sanjuan/octospec/cli/internal/config"
)

// Run removes the unmodified managed files for repo and keeps the ones edited
// by hand. It removes the state file when it finishes.
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
	if len(cfg.Files) == 0 {
		fmt.Println("nothing to uninstall")
		return nil
	}
	removed, kept := 0, 0
	for path, hash := range cfg.Files {
		b, err := os.ReadFile(path)
		if err != nil {
			continue // already gone
		}
		if config.Hash(b) != hash {
			fmt.Printf("kept (local edit): %s\n", path)
			kept++
			continue
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		pruneEmpty(filepath.Dir(path))
		removed++
	}
	if err := os.Remove(config.Path(repo)); err != nil && !os.IsNotExist(err) {
		return err
	}
	pruneEmpty(filepath.Join(repo, config.Dir))
	fmt.Printf("uninstalled %d file(s), kept %d\n", removed, kept)
	return nil
}

// pruneEmpty removes dir and its parents while they are empty.
func pruneEmpty(dir string) {
	for {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = filepath.Dir(dir)
	}
}
