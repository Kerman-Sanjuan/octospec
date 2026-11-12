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
	"os/exec"
	"strings"
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

// RemainingStages returns the stage sequence left for the issue, resuming
// from its status labels: an issue in review already has its pull request
// and resumes at ship; spec-ready and in-progress resume at apply; anything
// else (an explicitly named backlog issue) starts at spec.
func RemainingStages(i Issue) []Stage {
	stages := []Stage{
		{Name: "spec", Command: "/spec"},
		{Name: "apply", Command: "/apply"},
		{Name: "ship", Command: "/ship"},
		{Name: "archive", Command: "/archive"},
	}
	var from int
	switch {
	case hasLabel(i, StatusInReview):
		from = 2
	case hasLabel(i, StatusSpecReady), hasLabel(i, StatusInProgress):
		from = 1
	}
	return stages[from:]
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

// Options configures one terminal invocation of the autonomous mode.
type Options struct {
	Repo       string // repository root for gh (default: cwd)
	Issue      int    // explicit issue number; 0 selects automatically
	MaxRetries int    // per-stage bound; <= 0 becomes DefaultMaxRetries
}

// Plan is one run's computed state: what to work and where to resume.
type Plan struct {
	Issue      Issue
	MaxRetries int
}

// PlanFor computes the run plan for an issue.
func PlanFor(issue Issue, maxRetries int) Plan {
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}
	return Plan{Issue: issue, MaxRetries: maxRetries}
}

// StatusLabels returns the issue's status labels, for display.
func (i Issue) StatusLabels() []string {
	var out []string
	for _, l := range i.Labels {
		if strings.HasPrefix(l, "status:") {
			out = append(out, l)
		}
	}
	return out
}

// String renders the plan for the terminal.
func (p Plan) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "issue:   #%d %s\n", p.Issue.Number, p.Issue.Title)
	status := strings.Join(p.Issue.StatusLabels(), ", ")
	if status == "" {
		status = "(none)"
	}
	fmt.Fprintf(&b, "status:  %s\n", status)
	stages := RemainingStages(p.Issue)
	fmt.Fprintf(&b, "resume:  %s\n", stages[0].Command)
	names := make([]string, len(stages))
	for i, s := range stages {
		names[i] = s.Name
	}
	fmt.Fprintf(&b, "stages:  %s\n", strings.Join(names, " -> "))
	fmt.Fprintf(&b, "retries: %d per stage\n", p.MaxRetries)
	b.WriteString("\nThe loop runs in an agent session through the /auto command\n(cli/internal/payload/commands/auto.md).\n")
	return b.String()
}

// Run selects the issue and prints the run plan. The CLI cannot dispatch
// agents, so the terminal invocation computes and prints the plan; the loop
// itself executes in an agent session following the autonomous-mode command.
// With no --issue it selects the oldest open issue carrying an active
// status, and stops cleanly when the backlog has nothing resumable.
func Run(opts Options) error {
	gh, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh not on PATH")
	}
	var issue Issue
	if opts.Issue > 0 {
		out, err := ghRun(gh, opts.Repo, "issue", "view", fmt.Sprint(opts.Issue), "--json", "number,title,labels")
		if err != nil {
			return fmt.Errorf("gh issue view %d (is this a GitHub repo?): %w", opts.Issue, err)
		}
		if issue, err = parseIssue(out); err != nil {
			return err
		}
	} else {
		out, err := ghRun(gh, opts.Repo, "issue", "list", "--state", "open", "--json", "number,title,labels")
		if err != nil {
			return fmt.Errorf("gh issue list (is this a GitHub repo?): %w", err)
		}
		issue, err = SelectIssue(out)
		if errors.Is(err, ErrNoActiveIssue) {
			fmt.Println(ErrNoActiveIssue)
			return nil
		}
		if err != nil {
			return err
		}
	}
	fmt.Print(PlanFor(issue, opts.MaxRetries))
	return nil
}

// ghRun runs gh with args in repo (cwd when empty) and returns stdout.
func ghRun(gh, repo string, args ...string) ([]byte, error) {
	cmd := exec.Command(gh, args...)
	if repo != "" {
		cmd.Dir = repo
	}
	return cmd.Output()
}

// parseIssue parses `gh issue view --json number,title,labels` output.
func parseIssue(data []byte) (Issue, error) {
	var w ghIssue
	if err := json.Unmarshal(data, &w); err != nil {
		return Issue{}, fmt.Errorf("parse gh issue view: %w", err)
	}
	i := Issue{Number: w.Number, Title: w.Title}
	for _, l := range w.Labels {
		i.Labels = append(i.Labels, l.Name)
	}
	if i.Number == 0 {
		return Issue{}, fmt.Errorf("issue has no number: %s", data)
	}
	return i, nil
}
