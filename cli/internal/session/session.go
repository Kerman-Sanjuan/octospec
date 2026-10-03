// Package session isolates parallel agents in one git worktree per issue. Each
// worktree shares a single registry under the shared git directory, so an agent
// on one issue never touches another issue's checkout.
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Session is one active issue: its branch, its worktree, and its owner.
type Session struct {
	Issue     int       `json:"issue"`
	Slug      string    `json:"slug"`
	Branch    string    `json:"branch"`
	Worktree  string    `json:"worktree"`
	PID       int       `json:"pid"`
	Host      string    `json:"host"`
	Stage     string    `json:"stage,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

// Registry is the whole session list.
type Registry struct {
	Sessions []Session `json:"sessions"`
}

// dirName is the registry directory inside the shared git dir.
const dirName = "octospec"

// CommonDir returns the absolute shared git directory for repo, the one place
// every worktree points at.
func CommonDir(repo string) (string, error) {
	out, err := git(repo, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}
	dir := strings.TrimSpace(out)
	if !filepath.IsAbs(dir) {
		if repo == "" {
			repo, _ = os.Getwd()
		}
		dir = filepath.Join(repo, dir)
	}
	dir = filepath.Clean(dir)
	// Resolve symlinks so every worktree reports the same path (macOS /var vs
	// /private/var, for example).
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}
	return dir, nil
}

// registryDir is the octospec state directory inside the shared git dir.
func registryDir(commonDir string) string { return filepath.Join(commonDir, dirName) }

// registryPath is the sessions file.
func registryPath(commonDir string) string {
	return filepath.Join(registryDir(commonDir), "sessions.json")
}

// Load reads the registry for repo. A missing file is an empty registry.
func Load(repo string) (*Registry, error) {
	common, err := CommonDir(repo)
	if err != nil {
		return nil, err
	}
	return loadFrom(registryPath(common))
}

func loadFrom(path string) (*Registry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{}, nil
		}
		return nil, err
	}
	var r Registry
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func saveTo(path string, r *Registry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// withRegistryLock runs fn while holding an exclusive lock file, so two
// starting sessions do not corrupt the registry. A lock older than staleAfter
// is treated as left by a crashed process and reclaimed.
func withRegistryLock(commonDir string, fn func() error) error {
	dir := registryDir(commonDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "sessions.lock")
	const staleAfter = 5 * time.Second
	for i := 0; i < 100; i++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_ = f.Close()
			defer os.Remove(path)
			return fn()
		}
		if !os.IsExist(err) {
			return err
		}
		if fi, e := os.Stat(path); e == nil && time.Since(fi.ModTime()) > staleAfter {
			_ = os.Remove(path)
			continue
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("session registry is busy (%s)", path)
}

// git runs a git command in repo and returns its combined output.
func git(repo string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if repo != "" {
		cmd.Dir = repo
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
