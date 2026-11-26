## ADDED Requirements

### Requirement: Every requirement names the check that proves it

A change SHALL carry `coverage.md` mapping each requirement in its spec deltas to a symbol, in the form
`- <requirement title> -> <symbol>`, and the symbol SHALL exist in the tree outside `openspec/`.

#### Scenario: A requirement with no mapping

- **WHEN** a change adds a requirement and `coverage.md` has no line for it
- **THEN** G9 warns and names the requirement

#### Scenario: A mapping that points at nothing

- **WHEN** a coverage line names a symbol that appears nowhere outside `openspec/`
- **THEN** G9 warns and names the symbol

#### Scenario: A change without a spec delta

- **WHEN** the change touches no `specs/*.md`
- **THEN** G9 is skipped, because there is no requirement to prove

#### Scenario: Fully covered

- **WHEN** every requirement has a line and every symbol exists
- **THEN** G9 passes
