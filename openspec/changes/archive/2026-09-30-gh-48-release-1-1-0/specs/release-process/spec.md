## MODIFIED Requirements

### Requirement: The documented install path works
`curl | sh` and `go install .../octospec/cli/cmd/octospec@latest` SHALL install a runnable binary that reports the released version.

#### Scenario: Curl install
- **WHEN** the installer runs on a clean machine
- **THEN** `octospec version` prints `octospec` and the released version.

#### Scenario: Go install
- **WHEN** `go install .../cli/cmd/octospec@latest` runs
- **THEN** the installed binary reports the released version.
