// Package lock serialises the only stage that writes shared state: /archive, which folds a
// change's deltas into openspec/specs/. One repository, one lock file, one holder.
//
// The lock is a file created with O_CREATE|O_EXCL, not a daemon and not a system lock, because
// the person who has to clear a dead lock is not in the session that took it. The file names the
// pid, the host and the acquisition time, so `octospec lock status` can say who holds it and
// whether that holder is gone.
package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	// Dir is the lock directory, relative to the repository root. It is git-ignored.
	Dir = ".octospec/locks"
	// FileName is the lock that guards openspec/specs/.
	FileName = "specs.lock"
	// RemoteTTL is how old a lock taken on another host may be before it is reported stale.
	// There is no `kill -0` across hosts, so age is the only signal available.
	RemoteTTL = 15 * time.Minute

	pollEvery = 200 * time.Millisecond
)

var (
	// ErrHeld means somebody else holds the lock and the wait expired.
	ErrHeld = errors.New("archive lock is held")
	// ErrStale means the holder is gone. The lock is never stolen; it is reported.
	ErrStale = errors.New("archive lock is stale")
)

// Record is the lock file's content.
type Record struct {
	Change     string    `json:"change"`
	PID        int       `json:"pid"`
	Host       string    `json:"host"`
	AcquiredAt time.Time `json:"acquired_at"`
}

// State is what Status reports: whether the lock exists, who holds it, and whether that holder
// looks gone.
type State struct {
	Held   bool
	Record Record
	Stale  bool
	Why    string
}

// Path returns the lock file for a repository.
func Path(repo string) string { return filepath.Join(repo, filepath.FromSlash(Dir), FileName) }

// Status reads the lock without modifying it. A missing lock is a free lock, not an error.
func Status(repo string) (State, error) {
	b, err := os.ReadFile(Path(repo))
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	var r Record
	if err := json.Unmarshal(b, &r); err != nil {
		// A lock we cannot read is a lock we cannot trust: report it as stale rather than
		// pretending the repository is free.
		return State{Held: true, Stale: true, Why: "lock file is unreadable"}, nil
	}
	stale, why := isStale(r, time.Now())
	return State{Held: true, Record: r, Stale: stale, Why: why}, nil
}

// Acquire takes the lock for change, waiting at most wait for a free one. A stale lock is
// reported, never stolen: clearing it is a human decision.
func Acquire(repo, change string, wait time.Duration) (Record, error) {
	if change == "" {
		return Record{}, errors.New("a change name is required to take the lock")
	}
	p := Path(repo)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return Record{}, err
	}
	deadline := time.Now().Add(wait)
	for {
		rec, created, err := tryCreate(p, change)
		if err != nil {
			return Record{}, err
		}
		if created {
			return rec, nil
		}
		st, err := Status(repo)
		if err != nil {
			return Record{}, err
		}
		if st.Stale {
			return Record{}, fmt.Errorf("%w: held by %q (%s); clear it with `octospec lock release --force --change %s`",
				ErrStale, st.Record.Change, st.Why, st.Record.Change)
		}
		if !time.Now().Before(deadline) {
			return Record{}, fmt.Errorf("%w: held by %q since %s on %s",
				ErrHeld, st.Record.Change, st.Record.AcquiredAt.Format(time.RFC3339), st.Record.Host)
		}
		time.Sleep(pollEvery)
	}
}

// Release removes the lock. Without force it only removes a lock taken by change.
func Release(repo, change string, force bool) error {
	p := Path(repo)
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var r Record
	_ = json.Unmarshal(b, &r)
	if change == "" && !force {
		return fmt.Errorf("the lock is held by %q; name it with --change, or clear it with --force", r.Change)
	}
	if r.Change != change && !force {
		return fmt.Errorf("the lock is held by %q, not %q; use --force to clear it deliberately", r.Change, change)
	}
	return os.Remove(p)
}

func tryCreate(p, change string) (Record, bool, error) {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if errors.Is(err, os.ErrExist) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, err
	}
	defer f.Close()
	host, _ := os.Hostname()
	rec := Record{Change: change, PID: os.Getpid(), Host: host, AcquiredAt: time.Now().UTC()}
	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		os.Remove(p)
		return Record{}, false, err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		os.Remove(p)
		return Record{}, false, err
	}
	return rec, true, nil
}

func isStale(r Record, now time.Time) (bool, string) {
	host, _ := os.Hostname()
	if r.Host == "" {
		return true, "the lock names no host"
	}
	if r.Host != host {
		age := now.Sub(r.AcquiredAt)
		if age > RemoteTTL {
			return true, fmt.Sprintf("taken on %s %s ago", r.Host, age.Round(time.Second))
		}
		return false, ""
	}
	if r.PID <= 0 {
		return true, "the lock names no process"
	}
	if processAlive(r.PID) {
		return false, ""
	}
	return true, fmt.Sprintf("process %d is gone", r.PID)
}

func processAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = p.Signal(syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}
