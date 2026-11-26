# Proposal — archive lock

## Why

`/archive` is the only stage that writes to shared state. Today it is safe only because one
person runs it at a time; with a queue, two jobs can fold their deltas into the same
`openspec/specs/<capability>/spec.md` within seconds, and the loser's work disappears without
anyone noticing.

## What Changes

- A per-repository lock file guards `openspec/specs/` for the duration of the archive.
- A new `octospec lock` subcommand: `acquire`, `release`, `status`, so the lock is inspectable
  and clearable by a human.
- `/archive` acquires before reading and releases after committing; failing to acquire parks
  the change as `needs_human` with nothing written.
- `octospec doctor` reports whether the lock is held, by which change, and whether it is stale.
- The archive checks that its diff touches only the capabilities of its own change.

## Capabilities

**New capabilities**

- `archive-lock` — `specs/archive-lock/spec.md`

**Modified capabilities**

None.

## Impact

- Affected code: new `cli/internal/lock` package, a `lock` subcommand in `cli/cmd/octospec`,
  `cli/internal/doctor`, and the `/archive` command payload.
- Affected files on disk: `.octospec/locks/specs.lock` (git-ignored).
- No change to the OpenSpec schema, to the gates, or to the issue-side workflow.

## Changelog

### Added

- `octospec lock` to inspect and manage the lock that guards `openspec/specs/` during archive.
- `/archive` now parks instead of racing when another archive holds the lock.
