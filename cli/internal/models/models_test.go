package models

import (
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/config"
)

func TestSetWritesModels(t *testing.T) {
	dir := t.TempDir()
	if err := Set(Options{Repo: dir, Pairs: map[string]string{config.RoleReviewer: "opus"}}); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Models[config.RoleReviewer] != "opus" {
		t.Fatalf("models = %v", cfg.Models)
	}
}

func TestSetRejectsUnknownRole(t *testing.T) {
	if err := Set(Options{Repo: t.TempDir(), Pairs: map[string]string{"nope": "x"}}); err == nil {
		t.Fatal("expected an error for an unknown role")
	}
}

func TestSeedChoice(t *testing.T) {
	catalog := []string{"sonnet", "opus"}

	if choice, custom := seedChoice("", catalog); choice != "" || custom != "" {
		t.Fatalf("empty model should select the default, got %q %q", choice, custom)
	}
	if choice, custom := seedChoice("sonnet", catalog); choice != "sonnet" || custom != "" {
		t.Fatalf("a known model should select itself, got %q %q", choice, custom)
	}
	if choice, custom := seedChoice("gpt-4o", catalog); choice != customChoice || custom != "gpt-4o" {
		t.Fatalf("an unknown model should seed the custom input, got %q %q", choice, custom)
	}
}
