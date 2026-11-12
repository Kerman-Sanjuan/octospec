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
