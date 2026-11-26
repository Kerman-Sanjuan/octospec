package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/lock"
)

// runCapture runs fn and returns everything it wrote to stdout.
func runCapture(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan string, 1)
	go func() {
		var b strings.Builder
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			b.Write(buf[:n])
			if err != nil {
				break
			}
		}
		done <- b.String()
	}()
	fn()
	w.Close()
	os.Stdout = old
	return <-done
}

func TestVersionFallback(t *testing.T) {
	old := version
	version = "dev"
	t.Cleanup(func() { version = old })

	out := runCapture(t, func() {
		cmd := newVersionCmd()
		cmd.SetArgs(nil)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
	if got := strings.TrimSpace(out); got != "octospec dev" {
		t.Fatalf("version output = %q, want %q", got, "octospec dev")
	}
}

func TestVersionOverrideMatchesLdflags(t *testing.T) {
	// The release build injects -X main.version=<tag>. Simulate it here.
	old := version
	version = "1.0.0"
	t.Cleanup(func() { version = old })

	out := runCapture(t, func() {
		cmd := newVersionCmd()
		cmd.SetArgs(nil)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
	if got := strings.TrimSpace(out); got != "octospec 1.0.0" {
		t.Fatalf("version output = %q, want %q", got, "octospec 1.0.0")
	}
}

func TestLockExitCodeIsStable(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"held", fmt.Errorf("acquire: %w", lock.ErrHeld), 3},
		{"stale", fmt.Errorf("acquire: %w", lock.ErrStale), 4},
		{"other", errors.New("boom"), 1},
	}
	for _, c := range cases {
		if got := lockExitCode(c.err); got != c.want {
			t.Errorf("%s: exit %d, want %d", c.name, got, c.want)
		}
	}
}

func TestLockStatusReportsFree(t *testing.T) {
	dir := t.TempDir()
	out := runCapture(t, func() {
		cmd := newLockCmd()
		cmd.SetArgs([]string{"status", "--repo", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
	if got := strings.TrimSpace(out); got != "free" {
		t.Fatalf("status output = %q, want %q", got, "free")
	}
}

func TestLockReleaseWithoutChangeOrForceFails(t *testing.T) {
	dir := t.TempDir()
	if _, err := lock.Acquire(dir, "gh-1-demo", 0); err != nil {
		t.Fatal(err)
	}
	cmd := newLockCmd()
	cmd.SetArgs([]string{"release", "--repo", dir})
	if err := cmd.Execute(); err == nil {
		t.Fatal("release without --change or --force must fail")
	}
}
