package lock

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStatusFreeWhenNoFile(t *testing.T) {
	repo := t.TempDir()
	st, err := Status(repo)
	if err != nil {
		t.Fatal(err)
	}
	if st.Held || st.Stale {
		t.Fatalf("want free, got %+v", st)
	}
}

func TestAcquireThenStatus(t *testing.T) {
	repo := t.TempDir()
	rec, err := Acquire(repo, "gh-1-demo", 0)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Change != "gh-1-demo" || rec.PID != os.Getpid() {
		t.Fatalf("record = %+v", rec)
	}
	st, err := Status(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Held || st.Stale || st.Record.Change != "gh-1-demo" {
		t.Fatalf("state = %+v", st)
	}
}

func TestSecondAcquireIsHeldNotStolen(t *testing.T) {
	repo := t.TempDir()
	if _, err := Acquire(repo, "gh-1-demo", 0); err != nil {
		t.Fatal(err)
	}
	_, err := Acquire(repo, "gh-2-other", 0)
	if !errors.Is(err, ErrHeld) {
		t.Fatalf("want ErrHeld, got %v", err)
	}
	st, _ := Status(repo)
	if st.Record.Change != "gh-1-demo" {
		t.Fatalf("holder changed: %+v", st)
	}
}

func TestAcquireWaitsForRelease(t *testing.T) {
	repo := t.TempDir()
	if _, err := Acquire(repo, "gh-1-demo", 0); err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(150 * time.Millisecond)
		_ = Release(repo, "gh-1-demo", false)
	}()
	rec, err := Acquire(repo, "gh-2-other", 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Change != "gh-2-other" {
		t.Fatalf("record = %+v", rec)
	}
}

func TestAcquireRefusesAStaleLock(t *testing.T) {
	repo := t.TempDir()
	host, _ := os.Hostname()
	write(t, repo, Record{Change: "gh-1-demo", PID: 999999, Host: host, AcquiredAt: time.Now().UTC()})
	_, err := Acquire(repo, "gh-2-other", 0)
	if !errors.Is(err, ErrStale) {
		t.Fatalf("want ErrStale, got %v", err)
	}
	st, _ := Status(repo)
	if !st.Stale {
		t.Fatalf("want stale, got %+v", st)
	}
}

func TestRemoteLockIsStaleOnlyAfterTTL(t *testing.T) {
	repo := t.TempDir()
	write(t, repo, Record{Change: "gh-1-demo", PID: 1, Host: "another-host", AcquiredAt: time.Now().UTC()})
	st, _ := Status(repo)
	if st.Stale {
		t.Fatalf("a fresh remote lock is not stale: %+v", st)
	}
	write(t, repo, Record{Change: "gh-1-demo", PID: 1, Host: "another-host", AcquiredAt: time.Now().Add(-RemoteTTL - time.Minute).UTC()})
	st, _ = Status(repo)
	if !st.Stale {
		t.Fatalf("an old remote lock is stale: %+v", st)
	}
}

func TestUnreadableLockIsStale(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(Path(repo)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(repo), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, _ := Status(repo)
	if !st.Held || !st.Stale {
		t.Fatalf("want held and stale, got %+v", st)
	}
}

func TestReleaseIsGuardedByChange(t *testing.T) {
	repo := t.TempDir()
	if _, err := Acquire(repo, "gh-1-demo", 0); err != nil {
		t.Fatal(err)
	}
	if err := Release(repo, "gh-2-other", false); err == nil {
		t.Fatal("releasing somebody else's lock must fail")
	}
	if st, _ := Status(repo); !st.Held {
		t.Fatal("the lock was removed by a different change")
	}
	if err := Release(repo, "", false); err == nil {
		t.Fatal("releasing without a change and without force must fail")
	}
	if err := Release(repo, "gh-2-other", true); err != nil {
		t.Fatalf("force release: %v", err)
	}
	if st, _ := Status(repo); st.Held {
		t.Fatal("force release left the lock behind")
	}
}

func TestReleaseIsIdempotent(t *testing.T) {
	repo := t.TempDir()
	if err := Release(repo, "gh-1-demo", false); err != nil {
		t.Fatalf("releasing a free lock: %v", err)
	}
}

func write(t *testing.T, repo string, r Record) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(Path(repo)), 0o755); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(repo), b, 0o644); err != nil {
		t.Fatal(err)
	}
}
