// Package result is the one thing a stage leaves behind that a machine can read. A person reads the
// comment; the scheduler reads this. The format is uninteresting; the closed vocabulary of reasons
// is the point, because a timeout, a crash, an empty diff and a deliberate park need different
// retries and "failed" for all four is how a control plane starts lying.
package result

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Schema versions the object so a consumer can tell an old record from a broken one.
const Schema = 1

// Stages that write a result. The interactive ones (/idea, /bug) do not: a person is in the session
// and their answer is the result.
var Stages = []string{"apply", "ship", "review", "archive"}

// Outcomes answer "did it move". The reason answers "what now".
const (
	OutcomeOK      = "ok"
	OutcomeStalled = "stalled"
)

var Outcomes = []string{OutcomeOK, OutcomeStalled}

// Reasons are the closed list a retry policy switches on.
const (
	ReasonTimeout         = "timeout"
	ReasonCrashed         = "crashed"
	ReasonEmptyDiff       = "empty_diff"
	ReasonNoPRRequested   = "no_pr_requested"
	ReasonNeedsHuman      = "needs_human"
	ReasonBudgetExhausted = "budget_exhausted"
)

var Reasons = []string{
	ReasonTimeout, ReasonCrashed, ReasonEmptyDiff,
	ReasonNoPRRequested, ReasonNeedsHuman, ReasonBudgetExhausted,
}

// StagesRequiringArtifacts cannot honestly report success without pointing at something: an `ok`
// from these two with nothing attached is the shape of a stage that returned without working.
var StagesRequiringArtifacts = []string{"ship", "archive"}

// Result is one stage execution.
type Result struct {
	Schema     int       `json:"schema"`
	Stage      string    `json:"stage"`
	Issue      int       `json:"issue"`
	Change     string    `json:"change,omitempty"`
	Outcome    string    `json:"outcome"`
	Reason     string    `json:"reason,omitempty"`
	Artifacts  []string  `json:"artifacts"`
	DurationMS int64     `json:"duration_ms"`
	At         time.Time `json:"at"`
}

// Validate checks the object against the vocabulary and the two rules that keep it honest. It is
// deliberately strict: a result nobody can act on is worse than no result.
func (r *Result) Validate() error {
	if !contains(Stages, r.Stage) {
		return fmt.Errorf("unknown stage %q: expected one of %s", r.Stage, strings.Join(Stages, ", "))
	}
	if !contains(Outcomes, r.Outcome) {
		return fmt.Errorf("unknown outcome %q: expected one of %s", r.Outcome, strings.Join(Outcomes, ", "))
	}
	switch r.Outcome {
	case OutcomeStalled:
		if r.Reason == "" {
			return errors.New("a stalled result needs a reason")
		}
		if !contains(Reasons, r.Reason) {
			return fmt.Errorf("unknown reason %q: expected one of %s", r.Reason, strings.Join(Reasons, ", "))
		}
	case OutcomeOK:
		if r.Reason != "" {
			return fmt.Errorf("an ok result carries no reason, got %q: a finished stage that also "+
				"claims to have stopped is two answers to one question", r.Reason)
		}
		if contains(StagesRequiringArtifacts, r.Stage) && len(r.Artifacts) == 0 {
			return fmt.Errorf("an ok %s result must point at an artefact; an ok with nothing attached "+
				"is what a stage writes when it returned without working", r.Stage)
		}
	}
	if r.DurationMS < 0 {
		return errors.New("duration_ms cannot be negative")
	}
	if r.Issue <= 0 {
		return errors.New("a result needs the issue it belongs to")
	}
	return nil
}

// FileName is the record for a stage on a commit. A re-run on the same commit rewrites it; a new
// commit gets a new file, which is the question a scheduler asks: what happened last time here.
func FileName(stage, sha string) string {
	if sha == "" {
		return stage + ".json"
	}
	return stage + "-" + sha + ".json"
}

// Normalize fills what the writer owns, so the object on disk and the object printed are the same
// one. A caller that prints its own struct before writing it prints schema 0 and a zero timestamp.
func (r *Result) Normalize() {
	if r.Schema == 0 {
		r.Schema = Schema
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	if r.Artifacts == nil {
		r.Artifacts = []string{}
	}
}

// Write validates and stores r as dir/<stage>-<sha>.json, or dir/<stage>.json when the caller has no
// commit to key on.
func Write(dir string, r Result, sha string) (string, error) {
	r.Normalize()
	if err := r.Validate(); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, FileName(r.Stage, sha))
	return p, writeAtomic(p, b)
}

// Read parses and validates one record.
func Read(path string) (Result, error) {
	var r Result
	b, err := os.ReadFile(path)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return r, fmt.Errorf("%s: %w", path, err)
	}
	if err := r.Validate(); err != nil {
		return r, fmt.Errorf("%s: %w", path, err)
	}
	return r, nil
}

func writeAtomic(p string, b []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(p), ".result-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Chmod(name, 0o644); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, p); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
