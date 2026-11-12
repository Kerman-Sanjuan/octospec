package auto

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSelectIssuePicksOldestActive(t *testing.T) {
	data := []byte(`[
		{"number": 41, "title": "Backlog bug", "labels": [{"name": "status:backlog"}]},
		{"number": 57, "title": "In progress", "labels": [{"name": "type:bug"}, {"name": "status:in-progress"}]},
		{"number": 33, "title": "Spec ready", "labels": [{"name": "type:feature"}, {"name": "status:spec-ready"}]}
	]`)
	got, err := SelectIssue(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Number != 33 {
		t.Fatalf("number = %d, want the lowest-numbered active issue 33", got.Number)
	}
}

func TestSelectIssueSkipsBacklogOnly(t *testing.T) {
	data := []byte(`[
		{"number": 41, "title": "Backlog bug", "labels": [{"name": "status:backlog"}]},
		{"number": 12, "title": "Feature", "labels": [{"name": "type:feature"}, {"name": "status:backlog"}]}
	]`)
	if _, err := SelectIssue(data); !errors.Is(err, ErrNoActiveIssue) {
		t.Fatalf("err = %v, want ErrNoActiveIssue", err)
	}
}

func TestSelectIssueEmptyList(t *testing.T) {
	if _, err := SelectIssue([]byte(`[]`)); !errors.Is(err, ErrNoActiveIssue) {
		t.Fatalf("err = %v, want ErrNoActiveIssue", err)
	}
}

func TestSelectIssueWithoutLabels(t *testing.T) {
	// gh omits the labels key when --json drops it; the issue is skipped.
	data := []byte(`[{"number": 5, "title": "No labels"}]`)
	if _, err := SelectIssue(data); !errors.Is(err, ErrNoActiveIssue) {
		t.Fatalf("err = %v, want ErrNoActiveIssue", err)
	}
}

func TestSelectIssueParsesGhShape(t *testing.T) {
	// Shape produced by: gh issue list --state open --json number,title,labels
	data := []byte(`[
		{
			"number": 9,
			"title": "Real gh shape",
			"labels": [
				{"id": 1, "name": "type:feature", "color": "a2eeef", "description": "New capability"},
				{"id": 2, "name": "status:in-review", "color": "0e8a16", "description": "PR open"}
			]
		}
	]`)
	got, err := SelectIssue(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Number != 9 || got.Title != "Real gh shape" {
		t.Fatalf("issue = %+v", got)
	}
	if len(got.Labels) != 2 || got.Labels[1] != "status:in-review" {
		t.Fatalf("labels = %v", got.Labels)
	}
}

func TestSelectIssueBadJSON(t *testing.T) {
	if _, err := SelectIssue([]byte(`{`)); err == nil {
		t.Fatal("expected a parse error on malformed JSON")
	}
}

func TestHasActiveStatus(t *testing.T) {
	cases := []struct {
		labels []string
		want   bool
	}{
		{[]string{StatusSpecReady}, true},
		{[]string{StatusInProgress}, true},
		{[]string{StatusInReview}, true},
		{[]string{"status:backlog"}, false},
		{[]string{"type:feature"}, false},
		{nil, false},
	}
	for _, c := range cases {
		if got := (Issue{Labels: c.labels}).HasActiveStatus(); got != c.want {
			t.Errorf("HasActiveStatus(%v) = %v, want %v", c.labels, got, c.want)
		}
	}
}

func TestRemainingStagesResume(t *testing.T) {
	cases := []struct {
		labels []string
		want   []string
	}{
		{[]string{StatusInProgress}, []string{"apply", "ship", "archive"}},
		{[]string{StatusSpecReady}, []string{"apply", "ship", "archive"}},
		{[]string{StatusInProgress, StatusInReview}, []string{"ship", "archive"}},
		{[]string{StatusInReview}, []string{"ship", "archive"}},
		{[]string{"status:backlog"}, []string{"spec", "apply", "ship", "archive"}},
		{nil, []string{"spec", "apply", "ship", "archive"}},
	}
	for _, c := range cases {
		got := RemainingStages(Issue{Number: 1, Labels: c.labels})
		var names []string
		for _, s := range got {
			names = append(names, s.Name)
		}
		if len(names) != len(c.want) {
			t.Fatalf("RemainingStages(%v) = %v, want %v", c.labels, names, c.want)
		}
		for i := range names {
			if names[i] != c.want[i] {
				t.Errorf("RemainingStages(%v)[%d] = %q, want %q", c.labels, i, names[i], c.want[i])
			}
		}
	}
}

func TestNextStep(t *testing.T) {
	cases := []struct {
		name    string
		i       int
		attempt int
		max     int
		passed  bool
		want    Step
	}{
		{"pass advances", 0, 1, DefaultMaxRetries, true, StepAdvance},
		{"first failure retries", 0, 1, DefaultMaxRetries, false, StepRetry},
		{"mid failure retries", 1, 2, DefaultMaxRetries, false, StepRetry},
		{"bound reached gives up", 0, DefaultMaxRetries, DefaultMaxRetries, false, StepGiveUp},
		{"pass on last allowed attempt still advances", 2, DefaultMaxRetries, DefaultMaxRetries, true, StepAdvance},
		{"bound 1 gives up first try", 1, 1, 1, false, StepGiveUp},
	}
	for _, c := range cases {
		got, err := NextStep(c.i, c.attempt, c.max, c.passed)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: NextStep = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestNextStepRejectsBadInput(t *testing.T) {
	if _, err := NextStep(-1, 1, 3, true); err == nil {
		t.Error("expected an error for a negative stage index")
	}
	if _, err := NextStep(0, 0, 3, true); err == nil {
		t.Error("expected an error for a zero attempt number")
	}
}

func TestGiveUpComment(t *testing.T) {
	body := GiveUpComment("apply", 3, DefaultMaxRetries, "go test failed: TestX")
	for _, want := range []string{
		"## Autonomous run failed",
		"Stage `apply` failed 3 attempt(s), the retry bound is 3.",
		"go test failed: TestX",
		"The loop stopped working this issue.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("give-up comment missing %q:\n%s", want, body)
		}
	}
}

func TestPlanForDefaults(t *testing.T) {
	issue := Issue{Number: 109, Title: "Title", Labels: []string{StatusInProgress, "type:feature"}}
	p := PlanFor(issue, 0)
	if p.MaxRetries != DefaultMaxRetries {
		t.Fatalf("MaxRetries = %d, want the default %d", p.MaxRetries, DefaultMaxRetries)
	}
	p = PlanFor(issue, 5)
	if p.MaxRetries != 5 {
		t.Fatalf("MaxRetries = %d, want 5", p.MaxRetries)
	}
}

func TestPlanString(t *testing.T) {
	issue := Issue{Number: 109, Title: "Autonomous mode", Labels: []string{"type:feature", StatusInProgress}}
	got := PlanFor(issue, 0).String()
	for _, want := range []string{
		"issue:   #109 Autonomous mode",
		"status:  status:in-progress",
		"resume:  /apply",
		"stages:  apply -> ship -> archive",
		fmt.Sprintf("retries: %d per stage", DefaultMaxRetries),
		"/auto command",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("plan missing %q:\n%s", want, got)
		}
	}
}

func TestPlanStringInReview(t *testing.T) {
	issue := Issue{Number: 7, Title: "T", Labels: []string{StatusInReview}}
	got := PlanFor(issue, 0).String()
	if !strings.Contains(got, "resume:  /ship\n") {
		t.Fatalf("resume should be /ship for an issue in review:\n%s", got)
	}
}

func TestPlanStringNoStatus(t *testing.T) {
	got := PlanFor(Issue{Number: 1, Title: "T"}, 0).String()
	if !strings.Contains(got, "status:  (none)\n") {
		t.Fatalf("empty status should print (none):\n%s", got)
	}
}
