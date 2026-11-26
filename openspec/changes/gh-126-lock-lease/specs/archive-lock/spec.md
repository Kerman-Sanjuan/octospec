## MODIFIED Requirements

### Requirement: The lock is inspectable

`octospec lock status` and `octospec doctor` SHALL report who holds the lock, since when, and when
the lease expires, without modifying it.

#### Scenario: Doctor reports the holder

- **WHEN** a lock is held by `gh-122-archive-lock` with twenty minutes of lease left
- **THEN** `octospec doctor` prints the change, since when it is held, and when the lease expires

#### Scenario: No lock

- **WHEN** no lock file exists
- **THEN** `octospec lock status` exits zero and reports that the lock is free

#### Scenario: A fresh lock is not stale

- **WHEN** `octospec lock acquire` returns in one process and `octospec lock status` runs in another
- **THEN** the status reports `held`, never `stale`, and does not mention `--force`

### Requirement: A stale lock is detected, never stolen

A lock whose lease has expired SHALL be reported as stale and SHALL be cleared only by an explicit
command. A lock whose lease is valid SHALL NOT be reported as stale, whatever the state of the
process that wrote it.

#### Scenario: The lease expired

- **WHEN** the current time is past the lock's `expires_at`
- **THEN** `octospec lock status` reports it as stale and `acquire` refuses with a message naming
  the command to clear it

#### Scenario: The writer is gone but the lease is valid

- **WHEN** the record names a PID that no longer exists and `expires_at` is still in the future
- **THEN** the lock is held, because the PID is a hint and not the lease

#### Scenario: Clearing is explicit

- **WHEN** a human runs `octospec lock release --force --change <name>`
- **THEN** the lock is removed and the next archive proceeds

## ADDED Requirements

### Requirement: The holder keeps its lease alive

The holder SHALL be able to extend its own lease before it expires, and a heartbeat SHALL NOT
revive a lease that has already expired.

#### Scenario: Heartbeat extends

- **WHEN** the holder runs `octospec lock heartbeat --change <name>` on a valid lease
- **THEN** `expires_at` moves to now plus the TTL and the lock stays held by that change

#### Scenario: Heartbeat after expiry

- **WHEN** the holder heartbeats a lease that has already expired
- **THEN** the command fails, the lock stays stale, and nothing is revived

#### Scenario: Heartbeat by somebody else

- **WHEN** `octospec lock heartbeat --change other` runs against a lease held by `<name>`
- **THEN** it fails and reports which change holds the lock
