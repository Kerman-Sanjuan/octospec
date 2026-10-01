## ADDED Requirements

### Requirement: Lint, race, and coverage in CI
The `cli` job SHALL run a Go linter, `go test -race ./...`, and report coverage.

#### Scenario: Lint runs
- **WHEN** a pull request is opened
- **THEN** the linter runs and a finding fails the job.

#### Scenario: Race and coverage
- **WHEN** a pull request is opened
- **THEN** the tests run with the race detector and coverage is reported.
