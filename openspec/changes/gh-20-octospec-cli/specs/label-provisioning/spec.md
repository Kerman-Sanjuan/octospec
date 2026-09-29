## MODIFIED Requirements

### Requirement: Seeding provisions the workflow labels
`octospec seed` SHALL provision the six workflow labels: `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, `status:in-review`.

#### Scenario: Fresh repo gets the labels
- **WHEN** `octospec seed` runs against a repository that does not have the workflow labels
- **THEN** all six labels exist in that repository.

#### Scenario: Label names match the commands
- **WHEN** `/idea`, `/spec`, `/apply`, `/ship`, or `/archive` set a `type:*` or `status:*` label after seeding
- **THEN** the label exists and the command does not fail with "label not found".

### Requirement: Label provisioning is idempotent
Provisioning the labels SHALL be idempotent: re-running the seed SHALL NOT fail when a label already exists.

#### Scenario: Re-running the seed succeeds
- **WHEN** `octospec seed` runs a second time against an already-seeded repository
- **THEN** it exits successfully and the six labels remain present.
