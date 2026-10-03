package session

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// initRepo creates a git repo with one commit on main and returns its path.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"HOME="+dir,
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	return dir
}

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Run independent agents in parallel": "run-independent-agents-in-parallel",
		"Fix: the `foo` bar!":                "fix-the-foo-bar",
		"  spaced  out  ":                    "spaced-out",
		"":                                   "change",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBranchName(t *testing.T) {
	if got := branchName("type:bug", 12, "thing"); got != "fix/12-thing" {
		t.Errorf("bug branch = %q", got)
	}
	if got := branchName("type:feature", 12, "thing"); got != "feat/12-thing" {
		t.Errorf("feature branch = %q", got)
	}
}

func TestRegistryRoundTrip(t *testing.T) {
	repo := initRepo(t)
	common, err := CommonDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	path := registryPath(common)
	s := Session{Issue: 7, Slug: "seven", Branch: "feat/7-seven", Worktree: "/tmp/wt", PID: 1, Host: "h"}
	if err := saveTo(path, &Registry{Sessions: []Session{s}}); err != nil {
		t.Fatal(err)
	}
	got, err := loadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Sessions) != 1 || got.Sessions[0].Issue != 7 {
		t.Fatalf("round trip = %+v", got)
	}
}

func TestConcurrentRegistryWrites(t *testing.T) {
	repo := initRepo(t)
	common, err := CommonDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	path := registryPath(common)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = withRegistryLock(common, func() error {
				reg, err := loadFrom(path)
				if err != nil {
					return err
				}
				reg.upsert(Session{Issue: n, Branch: "b", Worktree: "w", Host: "h"})
				return saveTo(path, reg)
			})
		}(i)
	}
	wg.Wait()
	reg, err := loadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Sessions) != 10 {
		t.Fatalf("want 10 sessions, got %d", len(reg.Sessions))
	}
	b, _ := json.Marshal(reg)
	if !json.Valid(b) {
		t.Fatal("registry is not valid JSON")
	}
}

func TestRegistrySharedAcrossWorktrees(t *testing.T) {
	repo := initRepo(t)
	common, err := CommonDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	wt := filepath.Join(repo, worktreesDir, "9-nine")
	if _, err := git(repo, "worktree", "add", "-b", "feat/9-nine", wt, "main"); err != nil {
		t.Fatal(err)
	}
	if err := saveTo(registryPath(common), &Registry{Sessions: []Session{{Issue: 9, Worktree: wt}}}); err != nil {
		t.Fatal(err)
	}
	// The worktree's own common dir resolves to the same registry.
	gotCommon, err := CommonDir(wt)
	if err != nil {
		t.Fatal(err)
	}
	if gotCommon != common {
		t.Fatalf("common dirs differ: %s vs %s", gotCommon, common)
	}
	got, err := loadFrom(registryPath(gotCommon))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Sessions) != 1 || got.Sessions[0].Issue != 9 {
		t.Fatalf("worktree cannot see the session: %+v", got)
	}
}

func TestAddWorktreeReusesBranch(t *testing.T) {
	repo := initRepo(t)
	if _, err := git(repo, "branch", "feat/3-three"); err != nil {
		t.Fatal(err)
	}
	wt := filepath.Join(repo, worktreesDir, "3-three")
	if err := addWorktree(repo, wt, "feat/3-three"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Fatalf("worktree not created: %v", err)
	}
	// Idempotent: a second call is a no-op, not an error.
	if err := addWorktree(repo, wt, "feat/3-three"); err != nil {
		t.Fatalf("second addWorktree: %v", err)
	}
}

func TestFindAndRemove(t *testing.T) {
	repo := initRepo(t)
	common, err := CommonDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	path := registryPath(common)
	if err := saveTo(path, &Registry{Sessions: []Session{{Issue: 1}, {Issue: 2}}}); err != nil {
		t.Fatal(err)
	}
	reg, _ := loadFrom(path)
	if s := reg.find(2); s == nil || s.Issue != 2 {
		t.Fatal("find(2) failed")
	}
	reg.remove(1)
	if len(reg.Sessions) != 1 || reg.Sessions[0].Issue != 2 {
		t.Fatalf("remove left %+v", reg.Sessions)
	}
}

func TestDuplicateSession(t *testing.T) {
	prev := &Session{Issue: 5, Worktree: "/tmp/wt", Branch: "feat/5-five",
		PID: 1, Host: hostname(), StartedAt: time.Now()}
	err := duplicateError(prev)
	if !strings.Contains(err.Error(), "already has a session") {
		t.Fatalf("duplicate error = %v", err)
	}
	// Liveness is advisory: a huge pid is not alive on this host, a session on
	// another host is assumed live.
	if alive(1<<30, hostname()) {
		t.Fatal("a huge pid should not be alive")
	}
	if !alive(1<<30, "some-other-host") {
		t.Fatal("a session on another host is assumed live")
	}
}

func TestPresentUsesWorktree(t *testing.T) {
	dir := t.TempDir()
	if present(&Session{Worktree: ""}) {
		t.Fatal("a session with no worktree is not present")
	}
	if present(&Session{Worktree: filepath.Join(dir, "nope")}) {
		t.Fatal("a session whose worktree is gone is not present")
	}
	if !present(&Session{Worktree: dir}) {
		t.Fatal("a session whose worktree exists is present")
	}
}
