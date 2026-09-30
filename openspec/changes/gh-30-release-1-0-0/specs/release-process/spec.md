## ADDED Requirements

### Requirement: The first release is published
Cutting `v1.0.0` SHALL publish release binaries and checksums through the release workflow.

#### Scenario: Release artifacts exist
- **WHEN** the `v1.0.0` tag is pushed
- **THEN** the release has binaries for linux and darwin on amd64 and arm64, plus a checksum file.

### Requirement: The documented install path works
`curl | sh` and `go install .../octospec/cli/cmd/octospec@latest` SHALL install a runnable binary that reports version `1.0.0`.

#### Scenario: Curl install
- **WHEN** the installer runs on a clean machine
- **THEN** `octospec version` prints `octospec 1.0.0`.

#### Scenario: Go install
- **WHEN** `go install .../cli/cmd/octospec@latest` runs
- **THEN** the installed binary reports `octospec 1.0.0`.
