## ADDED Requirements

### Requirement: `/archive` commits and pushes the archive
`/archive` SHALL commit the archive (the moved change, the synced specs, and the filled Purposes) and push it, so the archive reaches the remote.

#### Scenario: Archive lands
- **WHEN** `/archive <change>` runs after the PR is merged
- **THEN** the archive is committed and pushed, and no manual commit is needed.

### Requirement: `/archive` closes the issue and clears status labels
`/archive` SHALL close the issue if it is still open and SHALL remove its `status:*` labels.

#### Scenario: Issue closed and labelled
- **WHEN** `/archive` finishes
- **THEN** the issue is closed and carries no `status:*` label.

### Requirement: `/archive` deletes the merged branch
`/archive` SHALL delete the merged local branch and MAY delete the remote branch.

#### Scenario: Branch removed
- **WHEN** `/archive` finishes
- **THEN** the merged branch no longer exists locally.

### Requirement: The archive is not blocked by branch protection
The branch-protection configuration SHALL allow the owner to push a maintenance commit such as an archive, while still requiring the `cli` and `gates` checks on pull requests.

#### Scenario: Owner pushes the archive
- **WHEN** `/archive` pushes to `main`
- **THEN** the push succeeds because the owner is not blocked by `enforce_admins`.
