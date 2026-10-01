## MODIFIED Requirements

### Requirement: The release runs an acceptance test
A release workflow SHALL run the gates, then run an acceptance test against the published binary on linux and macOS, and SHALL fail when either fails.

#### Scenario: The gates run at release
- **WHEN** the release acceptance workflow starts
- **THEN** it runs `scripts/check-gates.sh` and fails the release when a gate fails.

#### Scenario: The binary works
- **WHEN** the acceptance test downloads the published binary and runs `octospec version`, `octospec install`, and `octospec seed --no-labels` in a temporary directory
- **THEN** the commands, the 5 agents, and the schema are present and the test passes.
