# Proposal — lock lease

## Why

The lock shipped in #122 decides whether a holder is alive by asking the kernel about a PID. That
signal is correct only if the process that took the lock lives as long as the critical section, and
`octospec lock acquire` is a one-shot call: its PID is gone before the shell prompt returns. So the
first `lock status` after a legitimate `acquire` reports `stale`, and the `/archive` payload tells
the reader how to force-clear it. The guard recommends its own bypass.

That is worse than no lock: without one, two archives race and one of them fails on a merge
conflict. With a lying one, an agent follows documented advice and folds specs under another
archive's feet.

## What Changes

- The record carries `expires_at`. Liveness is the lease, not the process.
- `octospec lock heartbeat --change <change>` extends a valid lease and refuses to revive an
  expired one.
- `acquire` takes `--ttl`, default 30 minutes: long enough that a normal archive never needs the
  heartbeat, short enough that a crash does not block the repository for a day.
- The PID stays in the record, as a hint. It no longer decides anything.
- `status` prints holder, age and expiry, and only mentions `--force` once the lease has expired.
- `/archive` heartbeats between phases instead of trusting a timer it cannot see.

## Capabilities

- `archive-lock`: modified — inspectability and staleness, plus a new requirement for the lease.

## Impact

- `cli/internal/lock/` — record, staleness, heartbeat.
- `cli/cmd/octospec/main.go` — `lock heartbeat`, `--ttl`.
- `cli/internal/payload/commands/archive.md` — heartbeat between phases, no force advice.
- `docs/commands.md`, `docs/troubleshooting.md`.
- Records written by #122 have no `expires_at`; they are read as `acquired_at + default TTL`.

## Changelog

### Fixed

- The archive lock no longer reports a live holder as stale: liveness is a lease with a TTL, and
  the PID is only a hint.
- `octospec lock heartbeat` extends a lease, and `/archive` calls it between phases.
