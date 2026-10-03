# Issue #85: Run independent agents in parallel with one worktree per issue

https://github.com/Kerman-Sanjuan/octospec/issues/85

## User story

As a developer, I want to spawn several agents at once, each working a different
issue, so that no two of them clash and I can run five tabs without babysitting
them.

## Context / problem

octospec assumes one session at a time. Every stage mutates the same checkout:
`/spec` creates the change branch with `git checkout -b`, `/apply` commits on
it, `/archive` pushes `main` and deletes the branch. So two agents on two
different issues fight over one working tree, one HEAD, and one index, even
though their changes are independent.

The goal is N independent agents on N different issues, clash-free. The fix is a
git worktree per issue: each agent gets its own checkout on its own branch, so
HEAD and the index never collide. Git already refuses to check out a branch that
another worktree holds, so "same issue twice" fails structurally without a lock.

Most of an issue's footprint is already private (its `openspec/changes/<name>/`
folder, its `.openspec.yaml`, its labels). The shared surfaces that still need
care are the registry file and the push to `main` at archive.

## Requirements

- octospec SHALL support running several agents at once on different issues,
  each isolated in its own git worktree, with no shared working tree, HEAD, or
  index.
- `octospec session start <issue> [--tool <tool>]` SHALL create a git worktree
  at `.worktrees/<issue>-<slug>` on the change branch
  `feat|fix/<issue>-<slug>`, register the session, and, with `--tool`, launch
  that tool's agent inside the worktree. Without `--tool`, it SHALL print the
  worktree path.
- The session registry SHALL live under the shared git dir
  (`$(git rev-parse --git-common-dir)/octospec/`), never committed, so every
  worktree shares one list. Writes SHALL be serialized with a file lock.
- The registry SHALL record per issue: branch, worktree path, pid, host, stage,
  and start time.
- Same-issue protection SHALL be structural, not a hard lock: git refuses a
  second checkout of the branch, and `session start` SHALL warn when the issue
  already has a live session. There SHALL be no lock file and no stale-lock
  handling.
- `/spec`, `/apply`, `/ship`, and `/archive` SHALL run in the session's worktree
  when one exists and behave as today when none does. Sessions SHALL be
  optional; every stage SHALL work with no session.
- `/archive` SHALL push to `main` with a bounded fetch-rebase-push retry so
  independent agents do not lose an archive to a race. If the rebase conflicts,
  it SHALL stop and hand the conflict to the human.
- `octospec session list` SHALL show every active session with issue, branch,
  worktree, pid, and age.
- `octospec session end <issue>` SHALL remove the worktree and unregister the
  session.

## Success criteria

- Five agents in five tabs, each on a different issue, run `/spec` and `/apply`
  at the same time without touching each other's branch or working tree.
- Two archives pushed close together both land, via rebase-retry, or one stops
  cleanly on a real conflict.
- octospec works with no session at all, exactly as before; every existing stage
  passes unchanged.
- `openspec validate --all --strict` and the gates stay green.

## Out of scope

- A GitHub label or comment as a lock. Labels are not atomic, carry no
  liveness, and cannot isolate a working tree; coordination stays local.
  `status:in-progress` already signals an active issue to humans.
- Sharing the registry across machines or people (for example git hooks or
  coordination through a remote git ref).
- Real dispatch of the stage agents into subagents (tracked in #79).
- Agent attribution and approval records (tracked in #38).
