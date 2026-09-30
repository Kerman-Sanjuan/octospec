// Package install writes the canonical commands into each selected target and
// records the result so `update` can re-apply it.
package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/plan"
	"github.com/kerman-sanjuan/octospec/cli/internal/skills"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// Options controls an install run.
type Options struct {
	Tools []string // empty means all targets
	Repo  string   // repo root for repo-local targets; defaults to the cwd
}

// Run installs the canonical commands into every selected target.
func Run(opts Options) error {
	repo := opts.Repo
	if repo == "" {
		var err error
		if repo, err = os.Getwd(); err != nil {
			return err
		}
	}
	tools := opts.Tools
	if len(tools) == 0 {
		for _, t := range targets.Targets {
			tools = append(tools, t.Name)
		}
	}
	files, err := plan.Files(tools, repo, nil)
	if err != nil {
		return err
	}
	cfg, err := config.Load(repo)
	if err != nil {
		return err
	}
	if cfg.Files == nil {
		cfg.Files = map[string]string{}
	}
	cfg.Tools = tools
	for _, f := range files {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(f.Path, []byte(f.Content), 0o644); err != nil {
			return err
		}
		cfg.Files[f.Path] = config.Hash([]byte(f.Content))
	}
	if err := config.Save(repo, cfg); err != nil {
		return err
	}
	if err := skills.Install(repo, tools); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("installed %d files for %s\n", len(files), strings.Join(tools, ", "))
	return nil
}
