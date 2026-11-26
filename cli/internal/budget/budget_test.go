package budget

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultsAreConservative(t *testing.T) {
	if DefaultMinutes != 45 || DefaultTokens != 250_000 || DefaultFiles != 25 ||
		DefaultDiffLines != 800 || DefaultTasks != 12 {
		t.Fatalf("defaults changed without a note: %+v", Default())
	}
}

func TestLoadWithoutAFileIsTheDefaults(t *testing.T) {
	p, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if got := p.With("apply"); got != Default() {
		t.Fatalf("limits = %+v, want the defaults", got)
	}
}

func TestAStageOverridesOneLimitAndInheritsTheRest(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".octospec"), 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"max_minutes": 30, "stages": {"apply": {"max_files": 5}}}`
	if err := os.WriteFile(filepath.Join(repo, ".octospec", "budget.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Load(repo)
	if err != nil {
		t.Fatal(err)
	}
	apply := p.With("apply")
	if apply.MaxFiles != 5 {
		t.Fatalf("max_files = %d, want the stage override 5", apply.MaxFiles)
	}
	if apply.MaxMinutes != 30 {
		t.Fatalf("max_minutes = %d, want the global 30", apply.MaxMinutes)
	}
	if apply.MaxTokens != DefaultTokens {
		t.Fatalf("max_tokens = %d, want the default", apply.MaxTokens)
	}
}

func TestUnderBudgetIsFine(t *testing.T) {
	v := Check(Default(), Usage{Minutes: 3, Tokens: 1000, Files: 2, DiffLines: 40, Tasks: 1})
	if !v.OK || v.Hit != "" {
		t.Fatalf("verdict = %+v", v)
	}
	if v.Remaining["max_files"] != float64(DefaultFiles-2) {
		t.Fatalf("remaining files = %v", v.Remaining["max_files"])
	}
}

func TestEachLimitCanStopARun(t *testing.T) {
	cases := []struct {
		name string
		use  Usage
		hit  string
	}{
		{"time", Usage{Minutes: DefaultMinutes}, "max_minutes"},
		{"files", Usage{Files: DefaultFiles + 1}, "max_files"},
		{"diff", Usage{DiffLines: DefaultDiffLines + 1}, "max_diff_lines"},
		{"tasks", Usage{Tasks: DefaultTasks + 1}, "max_tasks"},
		{"tokens", Usage{Tokens: DefaultTokens + 1}, "max_tokens"},
	}
	for _, c := range cases {
		v := Check(Default(), c.use)
		if v.OK || v.Hit != c.hit {
			t.Fatalf("%s: verdict = %+v, want hit %s", c.name, v, c.hit)
		}
	}
}

func TestTimeExactlyAtTheLimitIsExhausted(t *testing.T) {
	// "Left <= 0" is deliberate: a run at its limit has no allowance, and rounding down to zero
	// because there are four seconds left is how a budget stops being a budget.
	v := Check(Default(), Usage{Minutes: DefaultMinutes})
	if v.OK {
		t.Fatal("a run exactly at its limit was allowed to continue")
	}
}

func TestUnknownTokensAreReportedAndDoNotBlock(t *testing.T) {
	for _, tokens := range []int{-1, -100} {
		v := Check(Default(), Usage{Minutes: 1, Tokens: tokens, Files: 1})
		if !v.OK {
			t.Fatalf("unknown tokens blocked the run: %+v", v)
		}
		if v.TokensSeen != "unknown" {
			t.Fatalf("tokens = %q, want unknown", v.TokensSeen)
		}
		if _, ok := v.Remaining["max_tokens"]; ok {
			t.Fatal("a limit we cannot measure must not report a remaining allowance")
		}
	}
}

func TestAnUnreadablePolicyIsAnError(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".octospec"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".octospec", "budget.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(repo); err == nil {
		t.Fatal("a budget nobody can parse was accepted")
	}
}

func TestAVerdictReadsAsASentence(t *testing.T) {
	if s := Check(Default(), Usage{Minutes: 1, Files: 1}).Render(); s == "" {
		t.Fatal("empty render")
	}
	if s := Check(Default(), Usage{Files: 99}).Render(); s == "" {
		t.Fatal("empty render")
	}
}
