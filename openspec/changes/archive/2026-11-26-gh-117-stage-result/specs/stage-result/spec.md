## ADDED Requirements

### Requirement: Every run leaves a result

A non-interactive stage SHALL write `.octospec/runs/<stage>-<sha>.json` with the schema version, the
stage, the issue, the change, the outcome, the reason, the artefacts, the duration and the instant.

#### Scenario: A stage reports success

- **WHEN** `ship` opens a pull request and writes its result with `outcome: ok`
- **THEN** the file names the stage, the commit, the pull request URL as an artefact, and no reason

#### Scenario: A stage reports why it stopped

- **WHEN** `apply` runs out of its time budget
- **THEN** it writes `outcome: stalled` and `reason: budget_exhausted`, and the same run twice writes
  the same reason

#### Scenario: The same commit twice

- **WHEN** the same stage runs twice on the same commit
- **THEN** it rewrites one file instead of adding another

### Requirement: The vocabulary is closed

The outcome SHALL be `ok` or `stalled`; the reason SHALL be one of `timeout`, `crashed`,
`empty_diff`, `no_pr_requested`, `needs_human`, `budget_exhausted`; the stage SHALL be one of the six
stages. Anything else SHALL fail rather than be stored.

#### Scenario: A typo in the reason

- **WHEN** a result is written with `reason: timout`
- **THEN** the command fails and names the accepted list

#### Scenario: A reason with no stall

- **WHEN** a result says `outcome: ok` and carries a reason
- **THEN** it is refused: a finished stage that also claims to have stopped is two answers to one
  question

### Requirement: A success must point at something

An `ok` result from `ship` SHALL carry at least one artefact, and so SHALL an `ok` from `archive`.

#### Scenario: A stage that did nothing

- **WHEN** `ship` reports `ok` with no artefact
- **THEN** the result is refused, because that is what a stage writes when it returned without
  working

### Requirement: The schema is published for consumers

The repository SHALL commit `docs/schema/result.schema.json`, and a test SHALL compare its enumerations
with the code's vocabularies.

#### Scenario: A reason added in one place only

- **WHEN** a reason is added to the Go vocabulary and not to the schema
- **THEN** the test fails
