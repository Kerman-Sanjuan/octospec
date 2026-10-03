## Context

octospec runs one change at a time in one checkout. The goal is N independent
agents on N different issues, clash-free. Three surfaces are shared today and
can clash:

- the working tree, HEAD, and index (every stage),
- the push to `main` at `/archive`,
- any registry file, once a registry exists.

Most of an issue's footprint is already private: its
`openspec/changes/<name>/` folder, its `.openspec.yaml`, and its labels. The
branch is per issue too, by the existing naming rule. So the shared surfaces are
narrow.

## Goals / Non-Goals

**Goals:**

- Several agents at once on different issues, each isolated in its own worktree.
- A single registry every worktree can see, without committing it.
- Sessions optional: no session means today's behaviour.
- Independent archives both land, or one stops cleanly.

**Non-Goals:**

- Sharing the registry across machines or people.
- A hard lock or stale-lock recovery.
- A GitHub label or comment as a lock.
- Real stage dispatch into subagents (#79).

## Decisions

### Isolation is a worktree per issue
Each session gets `git worktree add .worktrees/<issue>-<slug>` on the change
branch. Worktrees share one `.git` but have separate HEADs and indexes, so
independent agents cannot collide. Git also refuses to check out a branch that
another worktree holds, which makes "same issue twice" a structural error.
Rejected: separate clones (heavy, duplicate objects) and a single checkout with
a lock (serializes instead of parallelizes).

### The registry lives in the shared git dir
`$(git rev-parse --git-common-dir)/octospec/sessions.json` is shared by every
worktree, is inside `.git`, and can never be committed. Writes take a file lock
(an exclusive `O_CREATE` lock file) and rewrite the JSON, so two starting
sessions do not corrupt it. Rejected: `.octospec/sessions.json` in the main
checkout (each worktree is a separate directory, so lookups differ) and a
committed file (churn and conflicts).

### Same-issue protection is structural, not a lock
Because git blocks the duplicate branch, octospec needs no lock file, no pid
liveness checks, and no stale-lock recovery. `session start` only warns when the
issue already has a live session. Rejected: a hard per-issue lock (adds
staleness handling for a window git already covers).

### Stages stay optional-session-aware
`/spec`, `/apply`, `/ship`, and `/archive` first ask the registry for a session
for the issue. If one exists, they run in its worktree; if not, they run in the
current checkout as before. `/idea` is unchanged. This keeps the feature
additive: no session, no new failure mode.

### A stage refuses to hijack another issue's checkout
Before a stage switches or commits, it checks the current branch against the
issue it is working. If the branch belongs to a different issue, the stage stops
and tells the user to start a session or move to the right checkout. This is the
bug that motivated the change: `/spec 85` launched on `feat/84` would run
`git checkout -b`, stranding task 84. Rejected: warning only (still proceeds and
tramples the other task).

### Archive pushes with a bounded rebase-retry
The only shared surface independent agents still touch is `main`. `/archive`
fetches, rebases the archive commit onto `origin/main`, and pushes, retrying a
bounded number of times. A real conflict stops and reports, so a human resolves
it. Rejected: no retry (every second archive is manual) and force-push (unsafe).

## Risks / Trade-offs

- **[Worktrees confuse tooling that assumes one checkout]** -> Document
  `.worktrees/` and gitignore it; the stages detect the session.
- **[A registry entry outlives a crashed session]** -> `session list` shows pid
  and age; `session end` cleans up. No lock, so a stale entry blocks nothing.
- **[Rebase-retry masks a genuine conflict]** -> The retry is bounded and a
  conflict stops; it never forces.
- **[Slug drift between `session start` and `/spec`]** -> Both derive
  `<issue>-<slug>` the same way, from the issue title.

## Migration Plan

Additive. Existing installs keep working with no session. `.worktrees/` is
gitignored. Rollback is removing the `session` command and the session checks;
no data migration.

## Open Questions

- Should `session start` reuse an existing change branch if `/spec` already
  created it? Leaning yes, since `/spec` may run before a session exists.
