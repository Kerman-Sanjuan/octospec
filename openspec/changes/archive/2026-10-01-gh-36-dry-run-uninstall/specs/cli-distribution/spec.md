## ADDED Requirements

### Requirement: Dry-run for install and seed
`octospec install` and `octospec seed` SHALL support `--dry-run`, which writes nothing and prints what they would do.

#### Scenario: Install dry-run writes nothing
- **WHEN** `octospec install --dry-run` runs
- **THEN** it prints the files it would write and writes nothing.

#### Scenario: Seed dry-run writes nothing
- **WHEN** `octospec seed --dry-run` runs
- **THEN** it prints what it would install and writes nothing.

### Requirement: Uninstall removes the managed files
`octospec uninstall` SHALL remove the files octospec manages, using the recorded state, and SHALL leave hand-edited files alone.

#### Scenario: Remove managed files
- **WHEN** `octospec uninstall` runs
- **THEN** the unmodified managed files are removed and reported.

#### Scenario: Preserve a hand edit
- **WHEN** a managed file was edited by hand
- **THEN** `uninstall` keeps it and reports it.
