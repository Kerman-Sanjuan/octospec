package seed

import (
	"os"
	"path/filepath"
)

// SchemaDir returns the global OpenSpec schema directory for octospec.
func SchemaDir() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "openspec", "schemas", "octospec")
}

// InstallSchema writes the embedded octospec schema to the global directory.
func InstallSchema() (int, error) {
	n := 0
	if err := copyTree(schemaFS, "schema", SchemaDir(), &n); err != nil {
		return n, err
	}
	return n, nil
}
