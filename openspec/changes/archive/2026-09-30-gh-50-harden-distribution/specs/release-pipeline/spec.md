## ADDED Requirements

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
