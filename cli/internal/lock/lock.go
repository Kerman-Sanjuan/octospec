// Package lock serialises the only stage that writes shared state: /archive, which folds a
// change's deltas into openspec/specs/. One repository, one lock file, one holder.
//
// The lock is a file created with O_CREATE|O_EXCL, not a daemon and not a system lock, because
// the person who has to clear a dead lock is not in the session that took it. The file is a
// lease: it names the change, the host, the time it was taken and when it expires. Liveness is
// the expiry, never the process, because `octospec lock acquire` is a one-shot call whose pid is
// gone before the caller does any work. The pid stays in the record as a hint for whoever reads
// `lock status`.
package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// Dir is the lock directory, relative to the repository root. It is git-ignored.
	Dir = ".octospec/locks"
	// FileName is the lock that guards openspec/specs/.
	FileName = "specs.lock"
	// DefaultTTL is how long a lease lasts without a heartbeat. An archive is a fold, a changelog
	// entry, a commit and a push: minutes. Thirty leaves room for a slow run and bounds the damage
	// of a crash to half an hour of blocked work.
	DefaultTTL = 30 * time.Minute

	pollEvery = 200 * time.Millisecond
)

var (
	// ErrHeld means somebody else holds the lock and the wait expired.
	ErrHeld = errors.New("archive lock is held")
	// ErrStale means the lease expired. The lock is never stolen; it is reported.
	ErrStale = errors.New("archive lock is stale")
	// ErrFree means no lock was held when one was expected.
	ErrFree = errors.New("no archive lock is held")
)

// Record is the lock file's content.
type Record struct {
	Change     string    `json:"change"`
	PID        int       `json:"pid,omitempty"`
	Host       string    `json:"host"`
	AcquiredAt time.Time `json:"acquired_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// Expiry is when the lease ends. A record written before the lease existed carries no
// ExpiresAt; it is read as acquired_at plus the default TTL, which is what the old code did for
// holders on another host.
func (r Record) Expiry() time.Time {
	if !r.ExpiresAt.IsZero() {
		return r.ExpiresAt
	}
	return r.AcquiredAt.Add(DefaultTTL)
}

// State is what Status reports: whether the lock exists, who holds it, and whether its lease has
// run out.
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
	if r.Host == "" {
		return State{Held: true, Record: r, Stale: true, Why: "the lock names no host"}, nil
	}
	now := time.Now()
	if exp := r.Expiry(); !now.Before(exp) {
		return State{Held: true, Record: r, Stale: true,
			Why: fmt.Sprintf("lease expired %s ago", now.Sub(exp).Round(time.Second))}, nil
	}
	return State{Held: true, Record: r}, nil
}

// Acquire takes the lock for change, waiting at most wait for a free one, with a lease of ttl. A
// stale lock is reported, never stolen: clearing it is a human decision.
func Acquire(repo, change string, wait, ttl time.Duration) (Record, error) {
	if change == "" {
		return Record{}, errors.New("a change name is required to take the lock")
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	p := Path(repo)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return Record{}, err
	}
	deadline := time.Now().Add(wait)
	for {
		rec, created, err := tryCreate(p, change, ttl)
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
			return Record{}, fmt.Errorf("%w: %s; held by %q; clear it with `octospec lock release --force --change %s`",
				ErrStale, st.Why, st.Record.Change, st.Record.Change)
		}
		if !time.Now().Before(deadline) {
			return Record{}, fmt.Errorf("%w: held by %q since %s on %s, lease expires in %s",
				ErrHeld, st.Record.Change, st.Record.AcquiredAt.Format(time.RFC3339), st.Record.Host,
				time.Until(st.Record.Expiry()).Round(time.Second))
		}
		time.Sleep(pollEvery)
	}
}

// Heartbeat extends the lease of the holder that already owns the lock. It refuses to revive an
// expired lease: by then somebody else may hold it, and reviving would be the same silent
// double-holder in a new costume.
func Heartbeat(repo, change string, ttl time.Duration) (Record, error) {
	if change == "" {
		return Record{}, errors.New("a change name is required to extend the lease")
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	p := Path(repo)
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return Record{}, fmt.Errorf("%w: take it with `octospec lock acquire --change %s`", ErrFree, change)
	}
	if err != nil {
		return Record{}, err
	}
	var r Record
	if err := json.Unmarshal(b, &r); err != nil {
		return Record{}, fmt.Errorf("%w: the lock file is unreadable", ErrStale)
	}
	before := r
	if r.Change != change {
		return Record{}, fmt.Errorf("%w: held by %q, not %q", ErrHeld, r.Change, change)
	}
	now := time.Now()
	if exp := r.Expiry(); !now.Before(exp) {
		return Record{}, fmt.Errorf("%w: lease expired %s ago; take it again", ErrStale,
			now.Sub(exp).Round(time.Second))
	}
	r.ExpiresAt = now.UTC().Add(ttl)
	if err := replaceUnchanged(p, r, before); err != nil {
		return Record{}, err
	}
	return r, nil
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

func tryCreate(p, change string, ttl time.Duration) (Record, bool, error) {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if errors.Is(err, os.ErrExist) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, err
	}
	defer f.Close()
	host, _ := os.Hostname()
	now := time.Now().UTC()
	rec := Record{Change: change, PID: os.Getpid(), Host: host, AcquiredAt: now, ExpiresAt: now.Add(ttl)}
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

// replaceUnchanged rewrites the lease only if the lock it read is still the lock on disk: same
// change and same acquisition instant. A heartbeat that lost the race against a force-clear plus a
// fresh acquire would otherwise overwrite that holder's record with a lease extending work that is
// no longer happening. The window shrinks to microseconds rather than closing, which is the honest
// limit of a lock built on a file.
func replaceUnchanged(p string, r, before Record) error {
	cur, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: the lock disappeared while extending it", ErrHeld)
	}
	if err != nil {
		return err
	}
	var onDisk Record
	if err := json.Unmarshal(cur, &onDisk); err != nil || onDisk.Change != before.Change ||
		!onDisk.AcquiredAt.Equal(before.AcquiredAt) {
		return fmt.Errorf("%w: the lock changed while extending it, now held by %q", ErrHeld, onDisk.Change)
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), "specs.lock.*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, p); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
