## ADDED Requirements

### Requirement: `/spec` commits the change artifacts on the change branch
`/spec` SHALL create the change branch and commit the change artifacts on it before publishing.

#### Scenario: Artifacts are committed on the branch
- **WHEN** `/spec <issue>` finishes
- **THEN** the artifacts are committed on `feat|fix/<issue>-<slug>` and the working tree is clean.

#### Scenario: Main is untouched
- **WHEN** `/spec <issue>` finishes and the PR has not merged
- **THEN** `main` does not contain the change artifacts.

### Requirement: The branch type is derived from the issue label
`/spec` SHALL derive the branch prefix from the issue's `type:*` label: `type:feature` maps to `feat`, `type:bug` maps to `fix`.

#### Scenario: Feature branch
- **WHEN** `/spec` runs on a `type:feature` issue
- **THEN** the branch is `feat/<issue>-<slug>`.

#### Scenario: Bug branch
- **WHEN** `/spec` runs on a `type:bug` issue
- **THEN** the branch is `fix/<issue>-<slug>`.

### Requirement: Published comments are pinned to the artifact commit
Each published issue comment SHALL be headed with the commit SHA of the artifacts it was generated from.

#### Scenario: Comment carries the SHA
- **WHEN** `/spec` publishes an artifact to the issue
- **THEN** the comment is headed with the SHA of the commit that contains that artifact.

### Requirement: Re-running `/spec` refreshes in place
Re-running `/spec` SHALL reuse the existing change and branch and SHALL update the existing published comments instead of appending new ones.

#### Scenario: Revision updates the same comment
- **WHEN** the artifacts are edited, committed, and `/spec` is re-run
- **THEN** the existing published comments are updated to the new SHA and no duplicate comments are added.

### Requirement: `/apply` executes on the branch `/spec` created
`/apply` SHALL work on the branch created by `/spec`; it SHALL NOT create the branch and SHALL NOT commit a baseline.

#### Scenario: Apply starts clean on the branch
- **WHEN** `/apply` runs after `/spec`
- **THEN** it is already on the change branch with a clean working tree and it does not create a branch or commit a baseline.

### Requirement: Approach approval is a human action
The "approach approved" signal SHALL be set by the human, not by `/spec`.

#### Scenario: Drafted is not approved
- **WHEN** `/spec` finishes
- **THEN** the approval signal is not set until the human acts.
