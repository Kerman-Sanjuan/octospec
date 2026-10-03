package models

import (
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

func TestCatalogCoversEveryTool(t *testing.T) {
	for _, tp := range targets.Targets {
		models := Catalog([]string{tp.Name})
		if len(models) == 0 {
			t.Errorf("tool %q has no catalog entry", tp.Name)
		}
	}
}

func TestCatalogSkipsUnknownTool(t *testing.T) {
	if got := Catalog([]string{"nope"}); len(got) != 0 {
		t.Fatalf("unknown tool should yield no models, got %v", got)
	}
}

func TestCatalogDeduplicates(t *testing.T) {
	// pi and opencode share the same identifiers; asking for both must not
	// repeat them.
	got := Catalog([]string{"pi", "opencode"})
	seen := map[string]bool{}
	for _, m := range got {
		if seen[m] {
			t.Fatalf("duplicate model %q in %v", m, got)
		}
		seen[m] = true
	}
}

func TestCatalogEmptyToolsCoversAll(t *testing.T) {
	empty := Catalog(nil)
	if len(empty) == 0 {
		t.Fatal("an empty tool list should still offer the catalog")
	}
}
