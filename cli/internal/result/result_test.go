package result

import (
	"encoding/json"

	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func ok(stage string) Result {
	return Result{Stage: stage, Issue: 117, Change: "gh-117-stage-result", Outcome: OutcomeOK,
		Artifacts: []string{"https://github.com/Kerman-Sanjuan/octospec/pull/131"}}
}

func TestAValidOKResult(t *testing.T) {
	r := ok("ship")
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAnUnknownStageFails(t *testing.T) {
	r := ok("deploy")
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "unknown stage") {
		t.Fatalf("err = %v", err)
	}
}

func TestAnUnknownOutcomeFails(t *testing.T) {
	r := ok("apply")
	r.Outcome = "failed"
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "unknown outcome") {
		t.Fatalf("err = %v", err)
	}
}

func TestAStalledResultNeedsAReasonFromTheList(t *testing.T) {
	r := ok("apply")
	r.Outcome = OutcomeStalled
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "needs a reason") {
		t.Fatalf("err = %v", err)
	}
	r.Reason = "timout"
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "unknown reason") {
		t.Fatalf("a typo must fail and name the list, err = %v", err)
	}
	for _, reason := range Reasons {
		r.Reason = reason
		if err := r.Validate(); err != nil {
			t.Fatalf("reason %q rejected: %v", reason, err)
		}
	}
}

func TestAnOKResultCarriesNoReason(t *testing.T) {
	r := ok("apply")
	r.Reason = ReasonTimeout
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "carries no reason") {
		t.Fatalf("err = %v", err)
	}
}

func TestOKMustPointAtSomething(t *testing.T) {
	for _, stage := range StagesRequiringArtifacts {
		r := ok(stage)
		r.Artifacts = nil
		if err := r.Validate(); err == nil {
			t.Fatalf("%s reported ok with no artefact and was accepted", stage)
		}
	}
	// A stage that produces no artefact by nature is allowed to say ok.
	r := ok("apply")
	r.Artifacts = nil
	if err := r.Validate(); err != nil {
		t.Fatalf("apply may be ok without artefacts: %v", err)
	}
}

func TestAResultNeedsAnIssue(t *testing.T) {
	r := ok("review")
	r.Issue = 0
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "issue") {
		t.Fatalf("err = %v", err)
	}
}

func TestWriteKeysOnStageAndCommit(t *testing.T) {
	dir := t.TempDir()
	p, err := Write(dir, ok("ship"), "abc1234")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "ship-abc1234.json"); p != want {
		t.Fatalf("path = %s, want %s", p, want)
	}
	back, err := Read(p)
	if err != nil {
		t.Fatal(err)
	}
	if back.Schema != Schema || back.Stage != "ship" || back.At.IsZero() || len(back.Artifacts) != 1 {
		t.Fatalf("round trip lost fields: %+v", back)
	}
}

func TestWritingTheSameRunTwiceKeepsOneFile(t *testing.T) {
	dir := t.TempDir()
	if _, err := Write(dir, ok("ship"), "abc1234"); err != nil {
		t.Fatal(err)
	}
	if _, err := Write(dir, ok("ship"), "abc1234"); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("want one file, got %d", len(entries))
	}
}

func TestWithoutACommitTheFileHasNoSuffix(t *testing.T) {
	dir := t.TempDir()
	p, err := Write(dir, ok("review"), "")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(p) != "review.json" {
		t.Fatalf("path = %s", p)
	}
}

func TestTheSameStallTwiceHasTheSameReason(t *testing.T) {
	dir := t.TempDir()
	stalled := ok("apply")
	stalled.Outcome = OutcomeStalled
	stalled.Reason = ReasonBudgetExhausted
	stalled.Artifacts = nil
	p1, err := Write(dir, stalled, "abc1234")
	if err != nil {
		t.Fatal(err)
	}
	p2, err := Write(dir, stalled, "abc1234")
	if err != nil {
		t.Fatal(err)
	}
	if p1 != p2 {
		t.Fatalf("the same run wrote two files: %s and %s", p1, p2)
	}
	if r, err := Read(p1); err != nil || r.Reason != ReasonBudgetExhausted {
		t.Fatalf("reason = %q err = %v", r.Reason, err)
	}
}

// The schema is what a consumer in another language validates against, so its enums and the code's
// vocabularies must be the same list. Adding a reason in one place only fails here.
func TestThePublishedSchemaMatchesTheCode(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "docs", "schema", "result.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var s struct {
		Properties map[string]struct {
			Enum []string `json:"enum"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}
	cases := []struct {
		field string
		want  []string
	}{
		{"stage", Stages},
		{"outcome", Outcomes},
		{"reason", Reasons},
	}
	for _, c := range cases {
		got := s.Properties[c.field].Enum
		if !sameList(got, c.want) {
			t.Fatalf("schema %q = %v, code = %v", c.field, got, c.want)
		}
	}
	for _, f := range []string{"schema", "stage", "issue", "outcome", "artifacts", "duration_ms", "at"} {
		if !contains(s.Required, f) {
			t.Fatalf("schema does not require %q", f)
		}
	}
}

func TestDurationCannotBeNegative(t *testing.T) {
	r := ok("apply")
	r.DurationMS = -1
	if err := r.Validate(); err == nil {
		t.Fatal("a negative duration was accepted")
	}
}

func TestReadRejectsAFileThatDoesNotValidate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "apply.json")
	if err := os.WriteFile(p, []byte(`{"schema":1,"stage":"apply","issue":1,"outcome":"failed"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Read(p)
	if err == nil {
		t.Fatal("an invalid record was accepted")
	}
	if !strings.Contains(err.Error(), "apply.json") {
		t.Fatalf("the error must name the file, got %v", err)
	}
}

func TestAtIsSetByTheWriter(t *testing.T) {
	r := ok("review")
	r.At = time.Time{}
	dir := t.TempDir()
	p, err := Write(dir, r, "sha1")
	if err != nil {
		t.Fatal(err)
	}
	back, err := Read(p)
	if err != nil {
		t.Fatal(err)
	}
	if back.At.IsZero() {
		t.Fatal("the writer must stamp the result")
	}
}

func sameList(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
