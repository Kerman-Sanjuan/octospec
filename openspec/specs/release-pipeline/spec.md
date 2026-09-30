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
CI SHALL run on every pull request: build, `gofmt` check, `go vet`, `go test ./...`, plus the octospec gates.

#### Scenario: CI on a PR
- **WHEN** a pull request is opened
- **THEN** the build, format, vet, test, and octospec gate checks run.

### Requirement: Enforced gates
A failing gate SHALL block the pull request; gates are enforced controls, not prose.

#### Scenario: Red CI blocks
- **WHEN** any gate fails
- **THEN** the pull request cannot merge.

### Requirement: Release binaries
CD SHALL publish versioned release binaries on tag, installable by the `curl | sh` installer.

#### Scenario: Tag a release
- **WHEN** a version tag is pushed
- **THEN** release binaries are published and the `curl | sh` installer can install them.

### Requirement: The release runs an acceptance test
A release workflow SHALL run an acceptance test against the published binary on linux and macOS, and SHALL fail when the test fails.

#### Scenario: The binary works
- **WHEN** the acceptance test downloads the published binary and runs `octospec version`, `octospec install`, and `octospec seed --no-labels` in a temporary directory
- **THEN** the commands, the 5 agents, and the schema are present and the test passes.

#### Scenario: A broken binary
- **WHEN** the published binary fails the acceptance test
- **THEN** the release workflow fails.

### Requirement: Generated release notes
`goreleaser` SHALL generate the release notes from the commits.

#### Scenario: A release has notes
- **WHEN** a version is released
- **THEN** the GitHub release carries notes generated from the commit history.

