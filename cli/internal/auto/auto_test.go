package auto

import (
	"errors"
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
