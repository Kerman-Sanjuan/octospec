package review

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func rec(pr int, score float64, risk, verdict string) Record {
	return Record{PR: pr, Score: score, Risk: risk, Verdict: verdict, Threshold: 7,
		RiskWhy: []string{"burns the payload"}, At: time.Now().UTC()}
}

func TestWriteRecordCreatesOneFilePerPullRequest(t *testing.T) {
	dir := t.TempDir()
	p, err := WriteRecord(dir, rec(129, 10, RiskHigh, VerdictHuman))
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "129.json"); p != want {
		t.Fatalf("path = %s, want %s", p, want)
	}
	rs, skipped, err := ReadAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 || skipped != 0 {
		t.Fatalf("records = %d skipped = %d, want 1 and 0", len(rs), skipped)
	}
	if rs[0].Schema != LedgerSchema || rs[0].Risk != RiskHigh || len(rs[0].RiskWhy) != 1 {
		t.Fatalf("record = %+v", rs[0])
	}
	if rs[0].At.IsZero() {
		t.Fatal("a record without a timestamp cannot be ordered")
	}
}

func TestWritingTwiceRewritesTheSameRecord(t *testing.T) {
	dir := t.TempDir()
	if _, err := WriteRecord(dir, rec(129, 10, RiskHigh, VerdictHuman)); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteRecord(dir, rec(129, 4.5, RiskMedium, VerdictHuman)); err != nil {
		t.Fatal(err)
	}
	rs, _, err := ReadAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 {
		t.Fatalf("want one record, got %d", len(rs))
	}
	if rs[0].Score != 4.5 || rs[0].Risk != RiskMedium {
		t.Fatalf("the second review must replace the first: %+v", rs[0])
	}
}

func TestARecordNeedsAPullRequest(t *testing.T) {
	if _, err := WriteRecord(t.TempDir(), rec(0, 10, RiskLow, VerdictPass)); err == nil {
		t.Fatal("a record with no key must be refused")
	}
}

func TestReadAllSkipsUnreadableFiles(t *testing.T) {
	dir := t.TempDir()
	if _, err := WriteRecord(dir, rec(1, 9, RiskLow, VerdictPass)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "2.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), []byte("ignore me"), 0o644); err != nil {
		t.Fatal(err)
	}
	rs, skipped, err := ReadAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 || skipped != 1 {
		t.Fatalf("records = %d skipped = %d, want 1 and 1", len(rs), skipped)
	}
}

func TestReadAllOnAMissingDirectoryIsEmpty(t *testing.T) {
	rs, skipped, err := ReadAll(filepath.Join(t.TempDir(), "nope"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 0 || skipped != 0 {
		t.Fatalf("records = %d skipped = %d", len(rs), skipped)
	}
}

func TestSummariseCountsAndAverages(t *testing.T) {
	rs := []Record{
		rec(1, 10, RiskLow, VerdictPass),
		rec(2, 7, RiskMedium, VerdictPass),
		rec(3, 4, RiskHigh, VerdictHuman),
	}
	s := Summarise(rs)
	if s.Total != 3 || s.ByRisk[RiskLow] != 1 || s.ByRisk[RiskMedium] != 1 || s.ByRisk[RiskHigh] != 1 {
		t.Fatalf("summary = %+v", s)
	}
	if s.ByVerdict[VerdictPass] != 2 || s.ByVerdict[VerdictHuman] != 1 {
		t.Fatalf("verdicts = %+v", s.ByVerdict)
	}
	if got := s.MeanScore; got < 6.99 || got > 7.01 {
		t.Fatalf("mean = %.2f, want 7", got)
	}
	if e := Summarise(nil); e.Total != 0 || e.MeanScore != 0 {
		t.Fatalf("empty summary = %+v", e)
	}
}

func TestRenderSaysEmptyInsteadOfFailing(t *testing.T) {
	if out := Render(nil, 0); out == "" || out[:2] != "no" {
		t.Fatalf("render = %q", out)
	}
	out := Render([]Record{rec(129, 10, RiskHigh, VerdictHuman), rec(127, 6, RiskLow, VerdictPass)}, 1)
	for _, want := range []string{"PR", "129", "high", "human", "mean score", "skipped"} {
		if !contains(out, want) {
			t.Fatalf("render misses %q:\n%s", want, out)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
