# Issue #126 — The archive lock reports a live holder as stale

https://github.com/Kerman-Sanjuan/octospec/issues/126

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As the person merging an archive, I want the lock to tell a live holder apart from a dead one, so
that two archives never fold into the same capability because the guard lied about who was working.

## Context

`octospec lock` shipped in #122. Its liveness signal is the acquiring process's PID, which assumes
that process lives as long as the critical section. It does not: `octospec lock acquire` is a
one-shot CLI call, so the PID is dead microseconds later. The guard therefore reports a live holder
as stale the moment it is taken, and `/archive` tells the agent how to clear it — the exact theft
the lock exists to prevent.

Seen on the first two real archives after the merge (2026-11-26):

```
$ octospec lock acquire --change gh-122-archive-lock --wait 30s
lock held by gh-122-archive-lock (pid 217755 on hermesagent)
$ octospec lock status
stale: held by gh-122-archive-lock on hermesagent (process 217755 is gone);
       clear it with `octospec lock release --force --change gh-122-archive-lock`
```

## Requirements

- The lock SHALL consider a holder live until its lease expires, regardless of process liveness.
- The holder SHALL refresh its lease with `octospec lock heartbeat --change <change>`.
- The PID SHALL be recorded as a hint for diagnosis and SHALL NOT decide whether the lock is stale.
- `lock status` SHALL report who holds it, since when, and when the lease expires, and SHALL NOT
  suggest `--force` unless the lease has expired.
- `--force` SHALL remain the explicit, named way to break a valid lease.

## Success criteria

- `acquire` and `status` in separate processes report `held`, never `stale`.
- An expired lease reports `stale` and names the command that clears it.
- `heartbeat` extends a valid lease and fails on an expired one, without reviving it.
- `go test ./internal/lock/` covers all of the above and `scripts/e2e-test.sh` still passes.

## Out of scope

- Distributed consensus across machines beyond a TTL.
- Changing which stage takes the lock.
