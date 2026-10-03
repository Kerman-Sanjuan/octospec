## MODIFIED Requirements

### Requirement: G2 validate and completeness
A gate SHALL fail when `openspec validate --all --strict` fails, or when an unarchived change is incomplete. When `openspec` is unavailable, both the validate check (G2a) and the completeness check (G2b) SHALL report a SKIP or a FAIL, consistently, and SHALL NOT pass silently.

#### Scenario: A change is incomplete
- **WHEN** an unarchived change lacks a required artifact
- **THEN** G2 fails and names the change.

#### Scenario: OpenSpec is unavailable
- **WHEN** `openspec` is not on PATH
- **THEN** G2a and G2b report the same way, so neither passes silently.

### Requirement: G6 spec delta
A gate SHALL fail when a pull request touches `openspec/changes/<name>/` without a `specs/*.md` delta. The gate SHALL skip when no change is touched. The failure message SHALL name the real control, `skip_specs: true` in the change's `.openspec.yaml`.

#### Scenario: A change without a delta
- **WHEN** the pull request edits a change but adds no `specs/*.md`
- **THEN** G6 fails and the message names `skip_specs: true`.

### Requirement: G7 tasks checked
A gate SHALL warn, without failing, when an unarchived change has unchecked tasks. When `openspec` is unavailable, G7 SHALL report a SKIP or a FAIL, consistently with G2, and SHALL NOT pass silently.

#### Scenario: An unarchived change with open tasks
- **WHEN** an unarchived change has an unchecked task
- **THEN** G7 warns and the pull request can still merge.

#### Scenario: OpenSpec is unavailable
- **WHEN** `openspec` is not on PATH
- **THEN** G7 reports the same way as G2, so it does not pass silently.

### Requirement: Archived changes carry no unchecked tasks
Every archived change SHALL have all of its tasks checked, so an archived change never contradicts gate G7.

#### Scenario: Archived changes have no open tasks
- **WHEN** an archived change is inspected
- **THEN** every checkbox in its `tasks.md` is checked.
