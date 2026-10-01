// Package config stores what octospec installed, so `update` can re-apply it.
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

// Dir and File are the state location, relative to the repo root.
const (
	Dir  = ".octospec"
	File = "octospec.json"
)

// Model roles. A stage agent names one of these, and the rendered agent gets
// the model configured here. An empty model means the tool default.
const (
	RoleThinking    = "thinking"
	RoleImplementer = "implementer"
	RoleReviewer    = "reviewer"
)

// Roles is the fixed set of model roles, in display order.
var Roles = []string{RoleThinking, RoleImplementer, RoleReviewer}

// ValidRole reports whether role is one of the known model roles.
func ValidRole(role string) bool {
	for _, r := range Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Config is the saved install state: which tools, the scope per tool, the
// model per role, and the hash of every file octospec wrote (so `update` can
// tell a managed file from a hand edit).
type Config struct {
	Tools  []string          `json:"tools"`
	Scopes map[string]string `json:"scopes,omitempty"`
	Models map[string]string `json:"models,omitempty"`
	Files  map[string]string `json:"files"`
}

// Path returns the state file path for a repo.
func Path(repo string) string { return filepath.Join(repo, Dir, File) }

// Load reads the state, returning an empty Config when there is none.
func Load(repo string) (Config, error) {
	b, err := os.ReadFile(Path(repo))
	if err != nil {
		if os.IsNotExist(err) {
			return Config{Files: map[string]string{}, Models: map[string]string{}, Scopes: map[string]string{}}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	if c.Files == nil {
		c.Files = map[string]string{}
	}
	if c.Models == nil {
		c.Models = map[string]string{}
	}
	if c.Scopes == nil {
		c.Scopes = map[string]string{}
	}
	return c, nil
}

// Save writes the state.
func Save(repo string, c Config) error {
	if err := os.MkdirAll(filepath.Join(repo, Dir), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(repo), append(b, '\n'), 0o644)
}

// Hash returns the content hash recorded for a managed file.
func Hash(b []byte) string {
	h := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(h[:])
}
