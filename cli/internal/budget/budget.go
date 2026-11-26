// Package budget answers one question: has this run any allowance left. It is not a supervisor — the
// thing that runs the agent kills it — but it is the only place that knows what "too much" means,
// because the numbers belong to the repository, not to the tool.
package budget

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// PolicyPath is where a repository keeps its limits.
const PolicyPath = ".octospec/budget.json"

// Defaults are conservative on purpose. They are the numbers I would defend in review, not measured
// optima: the first calibration comes from the runs that hit them.
const (
	DefaultMinutes   = 45
	DefaultTokens    = 250_000
	DefaultFiles     = 25
	DefaultDiffLines = 800
	DefaultTasks     = 12
)

// Limits is what a run may spend. A zero field means "not declared here", which falls back to the
// default rather than to unlimited: a missing number must never mean no ceiling.
type Limits struct {
	MaxMinutes   int `json:"max_minutes"`
	MaxTokens    int `json:"max_tokens"`
	MaxFiles     int `json:"max_files"`
	MaxDiffLines int `json:"max_diff_lines"`
	MaxTasks     int `json:"max_tasks"`
}

// Default is what a repository gets when it declares nothing.
func Default() Limits {
	return Limits{
		MaxMinutes:   DefaultMinutes,
		MaxTokens:    DefaultTokens,
		MaxFiles:     DefaultFiles,
		MaxDiffLines: DefaultDiffLines,
		MaxTasks:     DefaultTasks,
	}
}

// Policy is the file: global limits plus per-stage overrides.
type Policy struct {
	Limits `json:",inline"`
	Stages map[string]Limits `json:"stages,omitempty"`
}

// With fills every unset field from the defaults, then applies the stage overrides on top.
func (p Policy) With(stage string) Limits {
	eff := p.Limits.or(Default())
	if s, ok := p.Stages[stage]; ok {
		eff = s.or(eff)
	}
	return eff
}

func (l Limits) or(fallback Limits) Limits {
	out := l
	if out.MaxMinutes <= 0 {
		out.MaxMinutes = fallback.MaxMinutes
	}
	if out.MaxTokens <= 0 {
		out.MaxTokens = fallback.MaxTokens
	}
	if out.MaxFiles <= 0 {
		out.MaxFiles = fallback.MaxFiles
	}
	if out.MaxDiffLines <= 0 {
		out.MaxDiffLines = fallback.MaxDiffLines
	}
	if out.MaxTasks <= 0 {
		out.MaxTasks = fallback.MaxTasks
	}
	return out
}

// Load reads the repository's policy. A missing file is the defaults; a broken one is an error,
// because a budget nobody can parse is a budget nobody is enforcing.
func Load(repo string) (Policy, error) {
	p := Policy{Limits: Default()}
	b, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(PolicyPath)))
	if errors.Is(err, os.ErrNotExist) {
		return p, nil
	}
	if err != nil {
		return p, err
	}
	var fromFile Policy
	if err := json.Unmarshal(b, &fromFile); err != nil {
		return p, fmt.Errorf("%s: %w", PolicyPath, err)
	}
	return fromFile, nil
}

// Usage is what the run has spent so far. Tokens below zero mean "the caller cannot say".
type Usage struct {
	Minutes   float64
	Tokens    int
	Files     int
	DiffLines int
	Tasks     int
}

// Verdict is the answer: whether there is allowance left, which limit was hit, and what remains.
type Verdict struct {
	OK         bool               `json:"ok"`
	Hit        string             `json:"hit,omitempty"`
	TokensSeen string             `json:"tokens,omitempty"` // reported | unknown
	Remaining  map[string]float64 `json:"remaining"`
}

// Check compares usage with the limits. The first limit that is exhausted wins, in the order a
// person would care about them: time, tokens, files, diff, tasks.
func Check(l Limits, u Usage) Verdict {
	v := Verdict{OK: true, Remaining: map[string]float64{}}
	if u.Tokens < 0 {
		v.TokensSeen = "unknown"
	} else {
		v.TokensSeen = "reported"
	}
	checks := []struct {
		name string
		used float64
		max  int
	}{
		{"max_minutes", u.Minutes, l.MaxMinutes},
		{"max_files", float64(u.Files), l.MaxFiles},
		{"max_diff_lines", float64(u.DiffLines), l.MaxDiffLines},
		{"max_tasks", float64(u.Tasks), l.MaxTasks},
	}
	if u.Tokens >= 0 {
		checks = append(checks, struct {
			name string
			used float64
			max  int
		}{"max_tokens", float64(u.Tokens), l.MaxTokens})
	}
	for _, c := range checks {
		left := float64(c.max) - c.used
		v.Remaining[c.name] = left
		if left <= 0 && v.OK {
			v.OK = false
			v.Hit = c.name
		}
	}
	return v
}

// Render is the line a person reads.
func (v Verdict) Render() string {
	if v.OK {
		return fmt.Sprintf("budget left: %.0f min, %.0f files, %.0f diff lines, %.0f tasks (tokens %s)",
			v.Remaining["max_minutes"], v.Remaining["max_files"], v.Remaining["max_diff_lines"],
			v.Remaining["max_tasks"], v.TokensSeen)
	}
	return fmt.Sprintf("budget exhausted: %s (tokens %s)", v.Hit, v.TokensSeen)
}
