// Package auto implements the autonomous mode: selecting the issue to work
// from the GitHub backlog, deciding the next step when a stage fails, and
// giving up visibly when a stage exhausts its retries.
//
// The CLI cannot dispatch agents, so the package computes and prints the run
// plan; the loop itself is driven by the autonomous-mode command
// (cli/internal/payload/commands/auto.md) executing in an agent session.
package auto

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Status labels an issue can be resumed from. Issues in status:backlog wait
// for a human decision and are never selected.
const (
	StatusSpecReady  = "status:spec-ready"
	StatusInProgress = "status:in-progress"
	StatusInReview   = "status:in-review"
)

// activeStatuses lists the status labels the autonomous mode may work on.
var activeStatuses = []string{StatusSpecReady, StatusInProgress, StatusInReview}

// ErrNoActiveIssue reports a backlog with nothing to resume.
var ErrNoActiveIssue = errors.New("no open issue carries an active status (status:spec-ready, status:in-progress, status:in-review)")

// Issue is one open GitHub issue considered for work.
type Issue struct {
	Number int
	Title  string
	Labels []string
}

// ghIssue mirrors gh issue list --json number,title,labels output; gh
// renders each label as an object, so the wire shape needs its own type.
type ghIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

// HasActiveStatus reports whether any label marks the issue resumable.
func (i Issue) HasActiveStatus() bool {
	for _, want := range activeStatuses {
		for _, l := range i.Labels {
			if l == want {
				return true
			}
		}
	}
	return false
}

// SelectIssue parses `gh issue list --state open --json number,title,labels`
// output and returns the issue to work: the lowest-numbered open issue that
// carries an active status. It returns ErrNoActiveIssue when the backlog has
// nothing resumable.
func SelectIssue(data []byte) (Issue, error) {
	var wire []ghIssue
	if err := json.Unmarshal(data, &wire); err != nil {
		return Issue{}, fmt.Errorf("parse gh issue list: %w", err)
	}
	var best Issue
	for _, g := range wire {
		issue := Issue{Number: g.Number, Title: g.Title}
		for _, l := range g.Labels {
			issue.Labels = append(issue.Labels, l.Name)
		}
		if !issue.HasActiveStatus() {
			continue
		}
		if best.Number == 0 || issue.Number < best.Number {
			best = issue
		}
	}
	if best.Number == 0 {
		return Issue{}, ErrNoActiveIssue
	}
	return best, nil
}

// Stage is one command in the autonomous loop, in order.
type Stage struct {
	Name    string // the payload command that runs it, e.g. "apply"
	Command string // the rendered slash command, e.g. "/apply"
}

// StagesAfterPlan returns the stages the loop runs for an issue whose plan
// is drafted. Idea and spec produce the issue and its plan, so a resumed
// loop always starts at apply.
func StagesAfterPlan() []Stage {
	return []Stage{
		{Name: "apply", Command: "/apply"},
		{Name: "ship", Command: "/ship"},
		{Name: "archive", Command: "/archive"},
	}
}

// RemainingStages returns the stage sequence left for the issue, resuming
// from its status labels: an issue in review already has its pull request
// and resumes at ship; anything else active resumes at apply.
func RemainingStages(i Issue) []Stage {
	stages := StagesAfterPlan()
	if hasLabel(i, StatusInReview) {
		return stages[1:]
	}
	return stages
}

func hasLabel(i Issue, want string) bool {
	for _, l := range i.Labels {
		if l == want {
			return true
		}
	}
	return false
}

// DefaultMaxRetries is the per-stage attempt bound before the loop gives up
// on the issue and moves to the next one.
const DefaultMaxRetries = 3

// Step tells the driver what to do after one stage attempt.
type Step int

const (
	// StepAdvance moves to the next stage: this one passed.
	StepAdvance Step = iota
	// StepRetry re-runs the same stage: it failed within the bound.
	StepRetry
	// StepGiveUp stops working the issue, visibly: retries exhausted.
	StepGiveUp
)

// NextStep decides the move after the stage at index i was attempted for the
// attempt-th time (1-based). passed reports whether the stage's gates were
// green. A failed stage retries the same stage while attempt <= maxRetries,
// then gives up.
func NextStep(i, attempt, maxRetries int, passed bool) (Step, error) {
	if i < 0 {
		return 0, fmt.Errorf("stage index %d is negative", i)
	}
	if attempt < 1 {
		return 0, fmt.Errorf("attempt %d is not 1-based", attempt)
	}
	if passed {
		return StepAdvance, nil
	}
	if attempt < maxRetries {
		return StepRetry, nil
	}
	return StepGiveUp, nil
}

// GiveUpComment returns the body to post on the issue when a stage exhausts
// its retry bound. It is the visible give-up contract: name the stage, the
// attempts, and the last error, so the morning reader can resume by hand.
func GiveUpComment(stage string, attempts int, maxRetries int, lastErr string) string {
	body := fmt.Sprintf("## Autonomous run failed\n\nStage `%s` failed %d attempt(s), the retry bound is %d.\n\nLast error:\n\n```\n%s\n```\n\nThe loop stopped working this issue. Re-run `/auto` after fixing, or work the stage manually.",
		stage, attempts, maxRetries, lastErr)
	return body
}
