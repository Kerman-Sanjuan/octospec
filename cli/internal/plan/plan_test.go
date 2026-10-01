package plan

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

func TestFilesDefaultLocal(t *testing.T) {
	repo := filepath.FromSlash("/repo")
	fs, err := Files([]string{"pi", "opencode"}, repo, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if !strings.HasPrefix(f.Path, repo) {
			t.Fatalf("default scope should be repo-local: %s", f.Path)
		}
		if f.Tool == "pi" && f.Kind == KindAgent {
			t.Fatalf("pi should have no agents: %s", f.Path)
		}
	}
}

func TestFilesGlobalScope(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	fs, err := Files([]string{"opencode"}, filepath.FromSlash("/repo"), nil, map[string]string{"opencode": targets.ScopeGlobal})
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) == 0 {
		t.Fatal("expected files")
	}
	for _, f := range fs {
		if !strings.HasPrefix(f.Path, home) {
			t.Fatalf("global scope should be under home: %s", f.Path)
		}
	}
}
