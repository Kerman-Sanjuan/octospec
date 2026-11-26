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

// The bug this change fixes: the writer is a one-shot process, so its pid is always gone by the
// time anybody looks. A dead pid inside a valid lease must not make the lock stale.
func TestADeadWriterDoesNotExpireTheLease(t *testing.T) {
	repo := t.TempDir()
	write(t, repo, Record{
		Change:     "gh-1-demo",
		PID:        999999,
		Host:       host(t),
		AcquiredAt: time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(10 * time.Minute),
	})
	st, err := Status(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Held || st.Stale {
		t.Fatalf("a valid lease with a dead pid must be held, got %+v", st)
	}
}

func TestAcquireThenStatus(t *testing.T) {
	repo := t.TempDir()
	rec, err := Acquire(repo, "gh-1-demo", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Change != "gh-1-demo" || rec.PID != os.Getpid() {
		t.Fatalf("record = %+v", rec)
	}
	if got := time.Until(rec.Expiry()); got <= 0 || got > DefaultTTL {
		t.Fatalf("lease = %s, want (0, %s]", got, DefaultTTL)
	}
	// Status re-reads the file, which is what a second process does.
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
	if _, err := Acquire(repo, "gh-1-demo", 0, 0); err != nil {
		t.Fatal(err)
	}
	_, err := Acquire(repo, "gh-2-other", 0, 0)
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
	if _, err := Acquire(repo, "gh-1-demo", 0, 0); err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(150 * time.Millisecond)
		_ = Release(repo, "gh-1-demo", false)
	}()
	rec, err := Acquire(repo, "gh-2-other", 3*time.Second, 0)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Change != "gh-2-other" {
		t.Fatalf("record = %+v", rec)
	}
}

func TestAcquireRefusesAStaleLock(t *testing.T) {
	repo := t.TempDir()
	write(t, repo, Record{
		Change:     "gh-1-demo",
		Host:       host(t),
		AcquiredAt: time.Now().UTC().Add(-40 * time.Minute),
		ExpiresAt:  time.Now().UTC().Add(-10 * time.Minute),
	})
	_, err := Acquire(repo, "gh-2-other", 0, 0)
	if !errors.Is(err, ErrStale) {
		t.Fatalf("want ErrStale, got %v", err)
	}
	st, _ := Status(repo)
	if !st.Stale || st.Record.Change != "gh-1-demo" {
		t.Fatalf("the lock must be reported stale and never stolen: %+v", st)
	}
}

func TestALockWithoutExpiryIsReadByAge(t *testing.T) {
	repo := t.TempDir()
	// A record written by #122: no expires_at, taken on another host.
	write(t, repo, Record{Change: "gh-1-demo", PID: 1, Host: "another-host", AcquiredAt: time.Now().UTC()})
	if st, _ := Status(repo); st.Stale {
		t.Fatalf("a fresh record without an expiry is not stale: %+v", st)
	}
	write(t, repo, Record{
		Change:     "gh-1-demo",
		PID:        1,
		Host:       "another-host",
		AcquiredAt: time.Now().Add(-DefaultTTL - time.Minute).UTC(),
	})
	if st, _ := Status(repo); !st.Stale {
		t.Fatalf("a record older than the default TTL is stale: %+v", st)
	}
}

func TestHeartbeatExtendsTheLease(t *testing.T) {
	repo := t.TempDir()
	rec, err := Acquire(repo, "gh-1-demo", 0, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	extended, err := Heartbeat(repo, "gh-1-demo", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if !extended.ExpiresAt.After(rec.ExpiresAt) {
		t.Fatalf("expiry = %s, want after %s", extended.ExpiresAt, rec.ExpiresAt)
	}
	if !extended.AcquiredAt.Equal(rec.AcquiredAt) {
		t.Fatalf("acquired_at moved: %s -> %s", rec.AcquiredAt, extended.AcquiredAt)
	}
	if st, _ := Status(repo); st.Stale {
		t.Fatalf("heartbeat must leave the lock valid: %+v", st)
	}
}

func TestHeartbeatDoesNotReviveAnExpiredLease(t *testing.T) {
	repo := t.TempDir()
	write(t, repo, Record{
		Change:     "gh-1-demo",
		Host:       host(t),
		AcquiredAt: time.Now().UTC().Add(-time.Hour),
		ExpiresAt:  time.Now().UTC().Add(-time.Minute),
	})
	if _, err := Heartbeat(repo, "gh-1-demo", time.Hour); !errors.Is(err, ErrStale) {
		t.Fatalf("want ErrStale, got %v", err)
	}
	if st, _ := Status(repo); !st.Stale {
		t.Fatalf("the lease must stay expired: %+v", st)
	}
}

func TestHeartbeatIsGuardedByChange(t *testing.T) {
	repo := t.TempDir()
	if _, err := Acquire(repo, "gh-1-demo", 0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := Heartbeat(repo, "gh-2-other", 0); !errors.Is(err, ErrHeld) {
		t.Fatalf("want ErrHeld, got %v", err)
	}
}

func TestHeartbeatWithoutALock(t *testing.T) {
	repo := t.TempDir()
	if _, err := Heartbeat(repo, "gh-1-demo", 0); !errors.Is(err, ErrFree) {
		t.Fatalf("want ErrFree, got %v", err)
	}
}

// A heartbeat that lost the race against a force-clear plus a fresh acquire must not overwrite the
// new holder's record with a lease for work that is no longer happening.
func TestHeartbeatDoesNotOverwriteANewHolder(t *testing.T) {
	repo := t.TempDir()
	before := Record{
		Change:     "gh-1-demo",
		Host:       host(t),
		AcquiredAt: time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(time.Minute),
	}
	next := before
	next.Change = "gh-2-other"
	next.AcquiredAt = before.AcquiredAt.Add(time.Second)
	write(t, repo, next)

	err := replaceUnchanged(Path(repo), before, before)
	if !errors.Is(err, ErrHeld) {
		t.Fatalf("want ErrHeld, got %v", err)
	}
	if st, _ := Status(repo); st.Record.Change != "gh-2-other" {
		t.Fatalf("the new holder was overwritten: %+v", st.Record)
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
	if _, err := Acquire(repo, "gh-1-demo", 0, 0); err != nil {
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

func host(t *testing.T) string {
	t.Helper()
	h, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	return h
}
