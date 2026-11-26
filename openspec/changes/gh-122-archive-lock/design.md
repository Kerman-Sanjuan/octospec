# Design — archive lock

## Context

`/archive` folds a change's delta specs into `openspec/specs/<capability>/spec.md` and is
therefore the only stage with shared mutable state. Today the serialisation is accidental: a
human runs one archive at a time. The queue ([`dispatch`](https://github.com/Kerman-Sanjuan/dispatch))
removes that accident, so the serialisation has to become explicit, and it has to be visible to
somebody who is not in the session.

## Goals / Non-Goals

**Goals**

- One archive writing to `openspec/specs/` at a time, per repository.
- A loser that parks with a reason, not a loser that half-writes.
- A human can see who holds the lock and clear a dead one deliberately.

**Non-Goals**

- Distributed consensus. One repository, one filesystem, one lock file.
- Locking branches or pull requests: git already serialises those.
- Automatic recovery of stale locks: stealing a lock when the holder is merely slow is worse
  than waiting for a human.

## Decisions

**Decision: a file created with `O_CREATE|O_EXCL`, not a lock daemon.**
The cheapest primitive that is atomic on every filesystem we care about, has no process to keep
alive, and survives a `kill -9` in a way a human can see and reason about. Rejected: `flock`
alone (invisible to `doctor`, released silently on process death with no record of who held it);
a Kubernetes `Lease` (only exists inside the cluster, and the same repository can be archived
from a laptop).

**Decision: the lock names a pid and a host.**
Stale means "the pid no longer exists on the host that took it", which is checkable without a
clock. For a lock taken on another host we fall back to age, with a generous TTL, because
`kill -0` across hosts is not a thing.

**Decision: never steal, always report.**
An automatic steal turns a slow archive into a corrupted spec. The cost of a false stale is a
human typing one command; the cost of a wrong steal is a lost requirement.

**Decision: the wait is bounded and declared.**
`acquire --wait <duration>` defaults to zero in interactive use and to a few minutes in the
`/archive` payload, so the queue's behaviour is a property of the change, not of luck.

**Decision: the diff check runs after the fold, before the commit.**
Archive cannot know which capabilities it will touch until it has folded them; checking before
means parsing the deltas twice. Failing before the commit keeps the tree recoverable.

## Risks / Trade-offs

- [A human archives from a second machine while the first is mid-run] → the second waits or
  parks; the TTL is generous (15 min) and the message names the host and pid.
- [The lock file is committed by accident] → `.octospec/locks/` is written to `.gitignore` by
  `octospec seed`, and `doctor` reports a tracked lock file.
- [The pid is reused by another process] → the lock also records the host and the acquisition
  time; a reused pid within the TTL is indistinguishable in practice from a slow holder, which
  is exactly why we do not steal.

## Migration Plan

Additive. Repositories without a lock file behave exactly as before; the first `/archive` after
the upgrade creates it. Rollback is deleting the subcommand and the payload step.

## Open Questions

- Should `acquire` retry with jitter by default in the queue, or leave retries to the scheduler?
  Today: leave them to the scheduler, which knows the retry vocabulary.
