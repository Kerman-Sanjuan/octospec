# release-process Specification

## Purpose
How a version is cut and recorded: a changelog grouped by impact, a documented release procedure, and semantic versioning.
## Requirements
### Requirement: A changelog grouped by impact
The repository SHALL keep a `CHANGELOG.md` whose entries are grouped by impact (Added, Changed, Fixed, Removed) and written for a reader.

#### Scenario: A release has an entry
- **WHEN** a version is released
- **THEN** `CHANGELOG.md` has a dated entry for it, grouped by impact.

### Requirement: A documented release process
The repository SHALL document the release steps and the tag convention.

#### Scenario: A maintainer cuts a release
- **WHEN** a maintainer follows the release document
- **THEN** they tag a version, the release workflow publishes the binaries, and the changelog is updated.

### Requirement: Semantic versioning
Release tags SHALL follow semantic versioning.

#### Scenario: A breaking change
- **WHEN** a change breaks behaviour consumers rely on
- **THEN** the release is a major version.

### Requirement: The first release is published
Cutting `v1.0.0` SHALL publish release binaries and checksums through the release workflow.

#### Scenario: Release artifacts exist
- **WHEN** the `v1.0.0` tag is pushed
- **THEN** the release has binaries for linux and darwin on amd64 and arm64, plus a checksum file.

### Requirement: The documented install path works
`curl | sh` and `go install .../octospec/cli/cmd/octospec@latest` SHALL install a runnable binary that reports the released version.

#### Scenario: Curl install
- **WHEN** the installer runs on a clean machine
- **THEN** `octospec version` prints `octospec` and the released version.

#### Scenario: Go install
- **WHEN** `go install .../cli/cmd/octospec@latest` runs
- **THEN** the installed binary reports the released version.

