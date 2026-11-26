## ADDED Requirements

### Requirement: The score is arithmetic over findings, not an opinion

The score SHALL be computed from structured findings by a deterministic function, and the same
findings SHALL always produce the same score.

#### Scenario: Same findings, same score

- **WHEN** the same findings are scored twice, on any host
- **THEN** both runs produce the same score, the same risk level and the same verdict

#### Scenario: A clean change

- **WHEN** a change produces no findings
- **THEN** the score is 10.0 and the verdict is to pass

#### Scenario: Each finding costs what its type declares

- **WHEN** a finding of a known type is scored
- **THEN** the score drops by that type's weight from `.octospec/review.json`, or by the built-in
  default when the repository does not override it

### Requirement: Five explicit axes

Every finding SHALL be attributed to one of the axes: correctness, architecture, domain fit, tests,
or UI and style. UI and style findings MAY be absent when the change has no UI surface.

#### Scenario: A finding with an unknown axis is rejected

- **WHEN** a finding names an axis the command does not know
- **THEN** the command fails with the finding id and the axis named, and scores nothing

#### Scenario: The comment shows the axes

- **WHEN** `/review` publishes its comment
- **THEN** the comment groups the costed findings by axis, so the reader sees where the score was
  lost rather than only how much

### Requirement: Risk is computed from facts the script can read

Risk SHALL be derived from at least three facts: whether the change touches a declared sensitive
path, the size of the diff, and whether the author has touched those paths before. It SHALL NOT be
the model's impression of danger.

#### Scenario: A sensitive path always asks for a human

- **WHEN** the change touches a file inside a declared sensitive path
- **THEN** the risk is `high` and the verdict asks for a human, whatever the score

#### Scenario: A large diff from an unfamiliar author

- **WHEN** the diff exceeds the configured size and the author has never touched those paths
- **THEN** the risk is `high`

#### Scenario: Neither of the two

- **WHEN** the change is small and inside paths the author has touched
- **THEN** the risk is `low` and the verdict depends only on the score

### Requirement: The policy lives in the repository

The threshold, the weights, the sensitive paths and the large-diff size SHALL be read from
`.octospec/review.json` in the adopting repository.

#### Scenario: Changing the threshold needs no release

- **WHEN** a maintainer raises the threshold in `.octospec/review.json`
- **THEN** the next run enforces the new number with no change to octospec

#### Scenario: Missing file means defaults

- **WHEN** the repository has no `.octospec/review.json`
- **THEN** the command uses the built-in defaults and says so in its output

### Requirement: Report first, enforce on request

The command SHALL always report the score, and SHALL block only when asked, so that a threshold can
be set from a real distribution.

#### Scenario: Report-only by default

- **WHEN** the command runs without `--enforce` and the score is below the threshold
- **THEN** it prints the score and the verdict and exits zero

#### Scenario: Enforcing

- **WHEN** the command runs with `--enforce` and the score is at or below the threshold, or the risk
  is high
- **THEN** it exits non-zero with a stable code, so a CI job blocks without parsing its output

### Requirement: One comment per pull request

`/review` SHALL publish exactly one comment per pull request and update it in place on every
re-run, and the comment SHALL state that its content comes from an agentic process.

#### Scenario: Re-run updates instead of adding

- **WHEN** `/review` runs twice on the same pull request
- **THEN** the second run updates the existing comment and the pull request still carries one

#### Scenario: The comment is labelled as agentic

- **WHEN** the comment is published
- **THEN** it states that a process, not a person, produced the score
