# Issue #122 — Guard the archive stage against two runs touching the same specs

https://github.com/Kerman-Sanjuan/octospec/issues/122

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As a maintainer, I want `/archive` to be safe when two changes finish at once, so that the
living specs under `openspec/specs/` are never merged by two agents at the same time.

## Context

`/archive` is the only stage that writes to shared state: it folds a change's deltas into
`openspec/specs/<capability>/spec.md`. With one human that is serial by accident. With a
queue, two jobs can finish within seconds and both write the same capability.

## Requirements

- `/archive` SHALL acquire a lock before touching `openspec/specs/` and release it after the
  commit, per repository.
- A lock that cannot be acquired within a timeout SHALL return `stalled / needs_human` without
  modifying anything, not a merge attempt.
- The lock SHALL be visible: `octospec doctor` SHALL report whether it is held and by which
  change.
- A stale lock (holder gone) SHALL be detectable and clearable with an explicit command, never
  automatically.
- The stage SHALL verify that the archive produced a diff touching only the capabilities of its
  own change.

## Success criteria

- Two archives started at the same time: one wins, the other parks with a clear reason and no
  partial write.
- A killed archive leaves a detectable stale lock.
- A test simulates the contention without needing two real machines.

## Out of scope

- Locks on branches or pull requests: git already handles that.
