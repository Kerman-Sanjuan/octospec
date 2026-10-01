# issue-intake Specification

## Purpose
How `/idea` and `/bug` create the issue from the body and apply the labels.
## Requirements
### Requirement: The intake commands create the issue from the body
`/idea` and `/bug` SHALL create the issue by passing the sections as the body, and SHALL NOT use `--template` together with `--body`.

#### Scenario: File a feature
- **WHEN** `/idea` files an issue
- **THEN** it runs `gh issue create` with the sections as the body and passes no `--template`.

#### Scenario: File a bug
- **WHEN** `/bug` files an issue
- **THEN** it runs `gh issue create` with the sections as the body and passes no `--template`.

### Requirement: The intake commands apply the labels
`/idea` SHALL add `type:feature` and `status:backlog`; `/bug` SHALL add `type:bug` and `status:backlog`.

#### Scenario: Feature labels
- **WHEN** `/idea` creates the issue
- **THEN** it carries `type:feature` and `status:backlog`.

#### Scenario: Bug labels
- **WHEN** `/bug` creates the issue
- **THEN** it carries `type:bug` and `status:backlog`.

### Requirement: The YAML forms remain the human contract
The repository SHALL keep `feature.yml` and `bug.yml` as the human web-form contract.

#### Scenario: A human opens the web form
- **WHEN** a human opens the new-issue form in the browser
- **THEN** `feature.yml` or `bug.yml` renders with its sections and its labels.

