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

