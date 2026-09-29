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

// Config is the saved install state: which tools, and the hash of every file
// octospec wrote (so `update` can tell a managed file from a hand edit).
type Config struct {
	Tools []string          `json:"tools"`
	Files map[string]string `json:"files"`
}

// Path returns the state file path for a repo.
func Path(repo string) string { return filepath.Join(repo, Dir, File) }

// Load reads the state, returning an empty Config when there is none.
func Load(repo string) (Config, error) {
	b, err := os.ReadFile(Path(repo))
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
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
