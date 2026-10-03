# release-pipeline Specification

## Purpose
The CLI's Go tests, the pull-request CI gates (format, vet, test, plus the octospec gates), and the tagged release binaries.
## Requirements
### Requirement: Go tests
The CLI SHALL have Go tests (unit and integration) runnable as `go test ./...`.

#### Scenario: Tests run
- **WHEN** `go test ./...` runs
- **THEN** the CLI's unit and integration tests execute.

### Requirement: Pull request CI
CI SHALL run on every pull request: build, `gofmt` check, `go vet`, `go test ./...`, plus the octospec gates, and the e2e step SHALL run against a real OpenSpec install.

#### Scenario: CI on a PR
- **WHEN** a pull request is opened
- **THEN** the build, format, vet, test, and octospec gate checks run.

#### Scenario: The e2e runs
- **WHEN** the `cli` job runs
- **THEN** it installs the OpenSpec CLI before the e2e step, so the step runs and does not skip.

#### Scenario: The e2e fails without OpenSpec
- **WHEN** `scripts/e2e-test.sh` runs and `openspec` is not on PATH
- **THEN** it exits non-zero with a clear message instead of printing a skip and passing.

### Requirement: Enforced gates
A failing gate SHALL block the pull request; gates are enforced controls, not prose.

#### Scenario: Red CI blocks
- **WHEN** any gate fails
- **THEN** the pull request cannot merge.

### Requirement: Release binaries
CD SHALL publish versioned release binaries on tag, installable by the `curl | sh` installer, and the binary SHALL report the tag it was built from, injected at build time rather than written by hand.

#### Scenario: Tag a release
- **WHEN** a version tag is pushed
- **THEN** release binaries are published and the `curl | sh` installer can install them.

#### Scenario: The binary reports the tag
- **WHEN** the release build injects the tag through `ldflags`
- **THEN** `octospec version` prints that version without the leading `v`.

#### Scenario: A local build
- **WHEN** a developer runs `go build` without the injected value
- **THEN** the binary still runs and `octospec version` prints the `dev` fallback.

### Requirement: The release runs an acceptance test
A release workflow SHALL run the gates, then run an acceptance test against the published binary on linux and macOS, and SHALL fail when either fails.

#### Scenario: The gates run at release
- **WHEN** the release acceptance workflow starts
- **THEN** it runs `scripts/check-gates.sh` and fails the release when a gate fails.

#### Scenario: The binary works
- **WHEN** the acceptance test downloads the published binary and runs `octospec version`, `octospec install`, and `octospec seed --no-labels` in a temporary directory
- **THEN** the commands, the 5 agents, and the schema are present and the test passes.

### Requirement: Generated release notes
`goreleaser` SHALL generate the release notes from the commits.

#### Scenario: A release has notes
- **WHEN** a version is released
- **THEN** the GitHub release carries notes generated from the commit history.

### Requirement: The gate script runs on a non-PR ref
`scripts/check-gates.sh` SHALL skip the PR-only branch gate when the ref is a long-lived branch, and SHALL report the skip instead of failing on it.

#### Scenario: Run on a long-lived branch
- **WHEN** `scripts/check-gates.sh` runs with `HEAD_REF=main`
- **THEN** it reports the branch gate as skipped and does not fail on it.

#### Scenario: Run on a PR branch
- **WHEN** `scripts/check-gates.sh` runs with a `feat|fix/<issue>-<slug>` ref
- **THEN** it checks the branch name and passes.

### Requirement: Lint, race, and coverage in CI
The `cli` job SHALL run a Go linter, `go test -race ./...`, and report coverage.

#### Scenario: Lint runs
- **WHEN** a pull request is opened
- **THEN** the linter runs and a finding fails the job.

#### Scenario: Race and coverage
- **WHEN** a pull request is opened
- **THEN** the tests run with the race detector and coverage is reported.

