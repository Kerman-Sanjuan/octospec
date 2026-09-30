## ADDED Requirements

### Requirement: Install the OpenSpec skills
`octospec install` SHALL make the OpenSpec skills available to each selected tool, so that commands that load them run.

#### Scenario: Skills present after install
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs in a repository
- **THEN** the OpenSpec skills exist in that tool's skill directory.

#### Scenario: Missing OpenSpec is reported
- **WHEN** the `openspec` CLI is not on the PATH
- **THEN** `install` reports that the skills were skipped, and the commands still install.
