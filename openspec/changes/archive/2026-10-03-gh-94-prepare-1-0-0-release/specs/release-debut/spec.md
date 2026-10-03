## ADDED Requirements

### Requirement: A single polished debut release
The repository SHALL present exactly one public release, v1.0.0, and SHALL NOT carry earlier releases or tags (v1.1.0 and the pre-polish v1.0.0 are removed first).

#### Scenario: Exactly one release
- **WHEN** a visitor opens the releases page after the debut
- **THEN** exactly one release, v1.0.0, is listed, with binaries, a checksums file, and release notes from the changelog.

#### Scenario: No stray tags
- **WHEN** the tags are listed
- **THEN** no `v1.1.0` tag and no duplicate `v1.0.0` tag remain.

### Requirement: One changelog section for the debut
`CHANGELOG.md` SHALL present a single dated 1.0.0 section for the debut, generated from the archived changes.

#### Scenario: The changelog has one 1.0.0 section
- **WHEN** a reader opens `CHANGELOG.md`
- **THEN** there is one `1.0.0` section and the earlier `1.1.0` section is folded into it.

### Requirement: End-to-end verification of the published binary
Before the debut is called done, the maintainer SHALL verify the published binary end to end: a fresh clone, the `curl | sh` install, `octospec seed`, `octospec install` into a tool, one full change through the loop, and `octospec version` matching the tag.

#### Scenario: Install and version
- **WHEN** the `curl | sh` installer runs on a clean machine after the tag
- **THEN** `octospec version` prints `octospec 1.0.0`.

#### Scenario: One change end to end
- **WHEN** a newcomer follows the README
- **THEN** they seed a repository, install into a tool, and run one change from `/idea` to `/archive`.

### Requirement: The repository About is filled in
The repository About SHALL have a description, a set of topics, and the project link.

#### Scenario: About is complete
- **WHEN** a visitor lands on the repository
- **THEN** the About panel shows a description, topics, and the project link.
