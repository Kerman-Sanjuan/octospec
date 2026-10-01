## ADDED Requirements

### Requirement: Each change carries a changelog entry
Every change SHALL carry a `## Changelog` section in its `proposal.md`. The section holds the entry the release will publish: an `Added`, `Changed`, `Fixed`, or `Removed` bullet list, or `None` when the change has no user-visible effect.

#### Scenario: A change with a user-visible effect
- **WHEN** a change adds a command
- **THEN** its `proposal.md` has a `## Changelog` section with the bullet that describes the command.

#### Scenario: A change with no user-visible effect
- **WHEN** a change only updates tests
- **THEN** its `## Changelog` section reads `None`.

### Requirement: Archive folds the entry into the changelog
`/archive` SHALL copy the change's `## Changelog` entry into `CHANGELOG.md`, under the current unreleased version heading, before it commits the archive.

#### Scenario: Archive updates the changelog
- **WHEN** `/archive` runs on a change whose `## Changelog` is not `None`
- **THEN** `CHANGELOG.md` gains the entry under the unreleased version heading, in the same commit as the archive.

#### Scenario: Archive skips an empty entry
- **WHEN** `/archive` runs on a change whose `## Changelog` reads `None`
- **THEN** `CHANGELOG.md` is unchanged.

### Requirement: The changelog is generated, not hand-edited
`CHANGELOG.md` SHALL be maintained through the archive step. A release SHALL only rename the unreleased heading to the version and date.

#### Scenario: Cut a release
- **WHEN** a maintainer cuts a release
- **THEN** the entries under the unreleased heading move under the new `X.Y.Z - YYYY-MM-DD` heading, and no entry is written by hand.
