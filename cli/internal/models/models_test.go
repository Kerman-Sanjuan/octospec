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
