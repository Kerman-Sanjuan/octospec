# parallel-sessions Specification

## Purpose
How several agents run at once, one per issue: each in its own git worktree on its own branch, tracked in a shared registry, with stages that refuse to hijack another issue's checkout.
## Requirements
### Requirement: One worktree per issue isolates parallel agents
octospec SHALL give each active issue its own git worktree on its own change
branch, so several agents run at once with no shared working tree, HEAD, or
index.

#### Scenario: Two issues, two checkouts
- **WHEN** sessions exist for issue 85 and issue 91
- **THEN** each has its own worktree and branch, and neither sees the other's
  uncommitted work.

#### Scenario: No session, no worktree
- **WHEN** no session exists for an issue
- **THEN** the stages run in the current checkout as before.

### Requirement: `session start` creates the worktree and registers the session
`octospec session start <issue> [--tool <tool>]` SHALL create a worktree at
`.worktrees/<issue>-<slug>` on the branch `feat|fix/<issue>-<slug>`, record the
session, and, with `--tool`, launch that tool's agent inside the worktree.
Without `--tool`, it SHALL print the worktree path.

#### Scenario: Start with a tool
- **WHEN** `octospec session start 85 --tool opencode` runs
- **THEN** a worktree exists at `.worktrees/85-<slug>` on
  `feat/85-<slug>` and the opencode agent is launched in it.

#### Scenario: Start without a tool
- **WHEN** `octospec session start 85` runs with no `--tool`
- **THEN** the worktree is created and its path is printed.

### Requirement: The registry lives in the shared git dir
The session registry SHALL live under the shared git directory
(`$(git rev-parse --git-common-dir)/octospec/`), SHALL NOT be committed, SHALL be
shared by every worktree, and SHALL serialize its writes with a file lock.

#### Scenario: Every worktree sees one list
- **WHEN** a session starts in worktree A and `octospec session list` runs in
  worktree B
- **THEN** the session started in A appears in the list.

#### Scenario: Concurrent writes do not corrupt the registry
- **WHEN** two sessions start at the same moment
- **THEN** both are recorded and the registry stays valid.

### Requirement: A session record names its owner
The registry SHALL record, per issue: the branch, the worktree path, the pid,
the host, the stage, and the start time.

#### Scenario: Owner is recorded
- **WHEN** a session starts
- **THEN** its record carries branch, worktree, pid, host, stage, and start
  time.

### Requirement: Same-issue protection is structural, not a lock
octospec SHALL rely on git refusing a second checkout of a held branch and SHALL
warn when the issue already has a live session. It SHALL NOT keep a lock file and
SHALL NOT implement stale-lock recovery.

#### Scenario: Second session warns
- **WHEN** a session already exists for issue 85 and `session start 85` runs
  again
- **THEN** octospec warns that the issue already has a live session and does not
  create a second worktree.

### Requirement: Stages are session-aware and sessions are optional
`/spec`, `/apply`, `/ship`, and `/archive` SHALL run in the session's worktree
when one exists for the issue, and SHALL behave as today when none exists.

#### Scenario: Stage runs in the session worktree
- **WHEN** a session exists for issue 85 and `/apply` runs
- **THEN** it works in the session's worktree.

#### Scenario: Stage runs without a session
- **WHEN** no session exists and `/apply` runs
- **THEN** it works in the current checkout, unchanged.

### Requirement: A stage refuses to hijack another issue's checkout
A stage SHALL detect when the current checkout is on a change branch that belongs
to a different issue, and SHALL refuse to switch, commit, or push on it. It SHALL
direct the user to start a session for the issue or move to the right checkout,
instead of switching the branch out from under the other task.

#### Scenario: Stage on the wrong issue's checkout
- **WHEN** the checkout is on `feat/84-<slug>` and `/spec 85` runs with no
  session for 85
- **THEN** `/spec` stops, does not run `git checkout`, and tells the user to
  start a session for 85 or move to the right checkout.

#### Scenario: The other task's branch is untouched
- **WHEN** a stage refuses because the checkout belongs to a different issue
- **THEN** the other issue's branch and working tree are unchanged.

### Requirement: `session list` shows the live sessions
`octospec session list` SHALL show every active session with its issue, branch,
worktree, pid, and age.

#### Scenario: List active sessions
- **WHEN** two sessions are live and `octospec session list` runs
- **THEN** both appear with issue, branch, worktree, pid, and age.

### Requirement: `session end` tears the session down
`octospec session end <issue>` SHALL remove the worktree and unregister the
session.

#### Scenario: End a session
- **WHEN** `octospec session end 85` runs
- **THEN** the worktree is removed and the session no longer appears in the
  registry.

