## MODIFIED Requirements

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
