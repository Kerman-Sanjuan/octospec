package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// worktreesDir is where sessions keep their isolated checkouts.
const worktreesDir = ".worktrees"

// Start creates the worktree for issue, records the session, and returns it.
// When tool is set it launches that tool inside the worktree; otherwise the
// caller prints the path. It warns and does not create a second worktree when
// the issue already has a live session.
func Start(repo string, issue int, tool string) (*Session, error) {
	if repo == "" {
		repo, _ = os.Getwd()
	}
	title, typ, err := issueMeta(repo, issue)
	if err != nil {
		return nil, err
	}
	slug := Slug(title)
	branch := branchName(typ, issue, slug)

	common, err := CommonDir(repo)
	if err != nil {
		return nil, err
	}
	reg, err := loadFrom(registryPath(common))
	if err != nil {
		return nil, err
	}
	if prev := reg.find(issue); prev != nil && alive(prev.PID, prev.Host) {
		return nil, fmt.Errorf("issue #%d already has a live session (pid %d on %s, started %s); run `octospec session list`",
			issue, prev.PID, prev.Host, prev.StartedAt.Format(time.RFC3339))
	}

	root, err := toplevel(repo)
	if err != nil {
		return nil, err
	}
	wt := filepath.Join(root, worktreesDir, fmt.Sprintf("%d-%s", issue, slug))
	if err := addWorktree(root, wt, branch); err != nil {
		return nil, err
	}

	s := Session{
		Issue:     issue,
		Slug:      slug,
		Branch:    branch,
		Worktree:  wt,
		PID:       os.Getpid(),
		Host:      hostname(),
		StartedAt: time.Now().UTC(),
	}
	err = withRegistryLock(common, func() error {
		reg, err := loadFrom(registryPath(common))
		if err != nil {
			return err
		}
		reg.upsert(s)
		return saveTo(registryPath(common), reg)
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// List returns the recorded sessions, newest first is not guaranteed; the
// caller formats them.
func List(repo string) ([]Session, error) {
	reg, err := Load(repo)
	if err != nil {
		return nil, err
	}
	return reg.Sessions, nil
}

// End removes the worktree for issue and unregisters the session.
func End(repo string, issue int) error {
	if repo == "" {
		repo, _ = os.Getwd()
	}
	common, err := CommonDir(repo)
	if err != nil {
		return err
	}
	reg, err := loadFrom(registryPath(common))
	if err != nil {
		return err
	}
	prev := reg.find(issue)
	if prev == nil {
		return fmt.Errorf("no session for issue #%d", issue)
	}
	root, err := toplevel(repo)
	if err != nil {
		return err
	}
	if _, statErr := os.Stat(prev.Worktree); statErr == nil {
		if out, err := git(root, "worktree", "remove", "--force", prev.Worktree); err != nil {
			return fmt.Errorf("remove worktree: %v: %s", err, strings.TrimSpace(out))
		}
	}
	return withRegistryLock(common, func() error {
		reg, err := loadFrom(registryPath(common))
		if err != nil {
			return err
		}
		reg.remove(issue)
		return saveTo(registryPath(common), reg)
	})
}

// Launch runs the tool's agent with its working directory set to the session's
// worktree, so the tab opens rooted in the right checkout. It reports whether
// it launched; a tool with no known command is not launched and the caller
// prints the path instead.
func Launch(s *Session, tool string) (bool, error) {
	bin := toolBin(tool)
	if bin == "" {
		return false, nil
	}
	if _, err := exec.LookPath(bin); err != nil {
		return false, nil
	}
	cmd := exec.Command(bin)
	cmd.Dir = s.Worktree
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return true, err
	}
	return true, nil
}

// toolBin maps a tool name to the command that opens its agent.
func toolBin(tool string) string {
	switch tool {
	case "opencode":
		return "opencode"
	case "claude":
		return "claude"
	case "copilot":
		return "copilot"
	case "pi":
		return "pi"
	default:
		return ""
	}
}

// Find returns the live session for issue, or nil.
func Find(repo string, issue int) (*Session, error) {
	reg, err := Load(repo)
	if err != nil {
		return nil, err
	}
	return reg.find(issue), nil
}

func (r *Registry) find(issue int) *Session {
	for i := range r.Sessions {
		if r.Sessions[i].Issue == issue {
			return &r.Sessions[i]
		}
	}
	return nil
}

func (r *Registry) upsert(s Session) {
	for i := range r.Sessions {
		if r.Sessions[i].Issue == s.Issue {
			r.Sessions[i] = s
			return
		}
	}
	r.Sessions = append(r.Sessions, s)
}

func (r *Registry) remove(issue int) {
	out := r.Sessions[:0]
	for _, s := range r.Sessions {
		if s.Issue != issue {
			out = append(out, s)
		}
	}
	r.Sessions = out
}

// addWorktree creates wt on branch. It reuses the worktree or the branch when
// they already exist, so a session whose process died does not block a restart.
func addWorktree(root, wt, branch string) error {
	if err := os.MkdirAll(filepath.Dir(wt), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(wt); err == nil {
		out, err := git(root, "worktree", "list", "--porcelain")
		if err == nil && strings.Contains(out, "worktree "+wt+"\n") {
			return nil
		}
		return fmt.Errorf("path already exists and is not a worktree for this repo: %s", wt)
	}
	if _, err := git(root, "rev-parse", "--verify", "refs/heads/"+branch); err == nil {
		_, err := git(root, "worktree", "add", wt, branch)
		return err
	}
	if _, err := git(root, "rev-parse", "--verify", "refs/remotes/origin/"+branch); err == nil {
		_, err := git(root, "worktree", "add", "--track", "-b", branch, wt, "origin/"+branch)
		return err
	}
	base := "HEAD"
	if _, err := git(root, "rev-parse", "--verify", "origin/main"); err == nil {
		base = "origin/main"
	}
	_, err := git(root, "worktree", "add", "-b", branch, wt, base)
	return err
}

// issueMeta fetches the issue title and its type label through gh.
func issueMeta(repo string, issue int) (title, typ string, err error) {
	cmd := exec.Command("gh", "issue", "view", fmt.Sprint(issue),
		"--json", "title,labels", "--jq", `.title, (.labels[].name | select(startswith("type:")))`)
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("gh issue view %d: %w", issue, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", "", fmt.Errorf("issue #%d has no title", issue)
	}
	title = strings.TrimSpace(lines[0])
	for _, l := range lines[1:] {
		if strings.HasPrefix(strings.TrimSpace(l), "type:") {
			typ = strings.TrimSpace(l)
		}
	}
	return title, typ, nil
}

// branchName builds feat|fix/<issue>-<slug> from the type label.
func branchName(typ string, issue int, slug string) string {
	prefix := "feat"
	if typ == "type:bug" {
		prefix = "fix"
	}
	return fmt.Sprintf("%s/%d-%s", prefix, issue, slug)
}

// Slug turns a title into a kebab-case slug for the change and branch name.
func Slug(title string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	s := strings.Trim(b.String(), "-")
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-")
	}
	if s == "" {
		s = "change"
	}
	return s
}

func toplevel(repo string) (string, error) {
	out, err := git(repo, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

// alive reports whether pid is a live process on host. A session on another
// host is assumed live, since we cannot check it here.
func alive(pid int, host string) bool {
	if host != hostname() {
		return true
	}
	if pid <= 0 {
		return false
	}
	return processAlive(pid)
}
