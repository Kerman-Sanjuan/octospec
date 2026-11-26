# archive-lock Specification

## Purpose
Serialises writes to `openspec/specs/`, so two archives never fold into the same capability at once.
## Requirements
### Requirement: The archive takes a per-repository lock

The `/archive` stage SHALL acquire a lock over `openspec/specs/` before it reads or writes any
capability spec, and SHALL release it only after the archive commit lands.

#### Scenario: First archive acquires the lock

- **WHEN** `/archive` starts for change `gh-122-archive-lock` and no lock is held
- **THEN** it writes `.octospec/locks/specs.lock` naming that change, proceeds, and removes the
  lock after committing the folded specs

#### Scenario: Lock is released on failure

- **WHEN** the archive fails after acquiring the lock
- **THEN** it releases the lock before returning, so the next run is not blocked by a dead holder

### Requirement: A second archive parks instead of racing

The stage SHALL NOT modify `openspec/specs/` when it cannot acquire the lock, and SHALL report a
specific reason instead of a generic failure.

#### Scenario: Contention within the timeout

- **WHEN** a second archive starts while another change holds the lock and the timeout expires
- **THEN** it exits with `outcome: stalled`, `reason: needs_human`, leaves `openspec/specs/`
  untouched, and reports which change holds the lock

#### Scenario: Waiting is bounded

- **WHEN** the holder releases within the configured timeout
- **THEN** the waiter acquires the lock and completes normally

### Requirement: The lock is inspectable

`octospec lock status` and `octospec doctor` SHALL report the lock state without modifying it.

#### Scenario: Doctor reports the holder

- **WHEN** a lock is held by `gh-122-archive-lock`
- **THEN** `octospec doctor` prints that the lock is held, by that change, and since when

#### Scenario: No lock

- **WHEN** no lock file exists
- **THEN** `octospec lock status` exits zero and reports that the lock is free

### Requirement: A stale lock is detected, never stolen

A lock whose holder is gone SHALL be reported as stale and SHALL be cleared only by an explicit
human command.

#### Scenario: Holder process is gone

- **WHEN** the lock names a process id that no longer exists on this host
- **THEN** `octospec lock status` reports it as stale and `acquire` refuses with a message naming
  the command to clear it

#### Scenario: Clearing is explicit

- **WHEN** a human runs `octospec lock release --force --change <name>`
- **THEN** the lock is removed and the next archive proceeds

### Requirement: The archive diff is scoped to its own change

The stage SHALL fail when the archive changes a capability that the change did not declare.

#### Scenario: Unexpected capability touched

- **WHEN** the archive diff modifies a capability spec that the change's deltas do not list
- **THEN** the stage fails with the capability named and does not commit

#### Scenario: Declared capabilities only

- **WHEN** the diff touches exactly the capabilities declared by the change
- **THEN** the archive commits

