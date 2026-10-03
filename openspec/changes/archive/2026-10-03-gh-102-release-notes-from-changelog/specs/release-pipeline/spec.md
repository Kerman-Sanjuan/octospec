## MODIFIED Requirements

### Requirement: Generated release notes
The release notes SHALL be the `CHANGELOG.md` section for the released version, grouped by impact, and SHALL NOT be the raw git commit log. When a version has no section, the release SHALL use a short fallback and SHALL NOT fail.

#### Scenario: A release has notes
- **WHEN** a version is released
- **THEN** the GitHub release carries the `CHANGELOG.md` section for that version, with its `### Added / Changed / Fixed / Removed` structure.

#### Scenario: A version without a section
- **WHEN** no `CHANGELOG.md` section matches the released version
- **THEN** the release still publishes, with a short fallback note.
