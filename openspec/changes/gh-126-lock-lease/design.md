# Design — lock lease

## Context

A lock has two questions to answer: *who holds it* and *is that holder still working*. #122 answered
the second with `kill -0` on the holder's PID, which is the correct answer when the holder is a
long-lived process and the wrong answer when it is a CLI invocation, which is our case.

The failure is not theoretical. Two archives of `main` on 2026-11-26 each acquired the lock and were
immediately told their own lock was stale, with the command to force-clear it.

## Goals / Non-Goals

**Goals**

- A lock that never lies about being free while somebody is working.
- A lock that does not block the repository forever when somebody crashes mid-archive.
- One concept for the reader: a lease with an expiry, refreshed by the holder.

**Non-Goals**

- Mutual exclusion across machines stronger than a TTL. Two hosts with skewed clocks can still
  collide; the archive commit will fail on the specs, which is recoverable.
- Detecting a *hung* holder faster than the TTL. A heartbeat is the holder's promise; if it stops
  promising, the TTL is the contract.

## Decisions

### The lease decides, the PID informs

`expires_at` is the only input to staleness. The PID remains in the record because when a human
reads `lock status` on the same host, "process 217755, running" is useful context — but a dead PID
in a valid lease means the holder finished its CLI call, not that the work stopped.

Rejected: keep the PID check but add a grace window. A five-minute archive would still be declared
stale at minute two, and the grace window would be tuned by whoever got burned last.

### The default TTL is 30 minutes

An archive is a fold, a changelog entry, a commit and a push: minutes. Thirty minutes means a
normal run never heartbeats and a crash costs half an hour of blocked work, which is short enough
to wait and long enough to not interrupt. `--ttl` overrides it for the rare long fold.

### Heartbeat never revives an expired lease

If the lease has expired, somebody may already have taken it. Refusing the heartbeat turns a silent
double-holder into an error the agent must report. Reviving would be the same bug in a new costume.

### Rewriting the lease is atomic, and refuses to clobber

`heartbeat` writes a uniquely named temporary file and renames it over the lock, so a reader never
sees half a record and two heartbeats never share a temporary path. Before the rename it re-reads
the lock and checks that the change and the acquisition instant are the ones it read; if they
moved, somebody cleared and re-acquired in between and the heartbeat fails instead of extending a
lease for work that is no longer happening. `O_EXCL` still governs creation, so acquire stays a
single atomic act.

The window between the check and the rename is microseconds and stays open: closing it needs an
advisory lock on the file, which buys little for a guard whose worst case is a repository where two
archives are minutes apart.

### Records without `expires_at` are read as `acquired_at + default TTL`

The field is new and the tool is public. A record written by #122 is expired-or-not by age, which is
exactly what the old code did for remote holders. No migration step, no version field.

## Risks

- A holder that dies mid-archive blocks the repository for 30 minutes. Acceptable: the alternative
  is a guard that is wrong immediately.
- An agent that heartbeats without working keeps a lease alive forever. Bounded in practice by the
  budgets of the `run` stage (#121), not by the lock.
