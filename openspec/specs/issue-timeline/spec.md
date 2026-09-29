# issue-timeline Specification

## Purpose
The fixed, three-layer comment structure an issue carries - the issue body, `## Plan`, and `## Implementation` - and how each layer is kept in sync in place.
## Requirements
### Requirement: `/spec` posts a single plan comment
`/spec` SHALL publish exactly one comment on the issue (the plan), and SHALL NOT post one comment per artifact.

#### Scenario: One plan comment
- **WHEN** `/spec` finishes on a fresh issue
- **THEN** the issue has exactly one comment authored by the plan step, headed `## Plan`.

#### Scenario: Re-run refreshes the plan comment
- **WHEN** the artifacts change and `/spec` is re-run
- **THEN** the same plan comment is updated in place and no new comment is added.

### Requirement: The plan comment has a fixed structure
The plan comment SHALL consolidate the artifacts under fixed headings: `### Proposal`, `### Capabilities`, `### Design`, `### Tasks`.

#### Scenario: Fixed headings present
- **WHEN** the plan comment is rendered
- **THEN** it contains `### Proposal`, `### Capabilities`, `### Design`, and `### Tasks`, in that order.

### Requirement: `/apply` posts a single implementation comment
`/apply` SHALL publish exactly one comment on the issue (the implementation) containing the task checklist with a short insight per task, and SHALL NOT post one comment per group.

#### Scenario: One implementation comment
- **WHEN** `/apply` completes a group
- **THEN** the issue has exactly one implementation comment, headed `## Implementation`, containing the full checklist with an insight on each done task.

#### Scenario: Progress refreshes in place
- **WHEN** `/apply` finishes the next group
- **THEN** the same implementation comment is updated and no new comment is added.

### Requirement: Each layer is recognisable
Every layer SHALL carry a stable, predefined header so it can be identified: the issue body, `## Plan` (`/spec`), and `## Implementation` (`/apply`).

#### Scenario: Layer headers
- **WHEN** a reader scans the issue
- **THEN** the plan and the implementation are each identifiable by their fixed header.

### Requirement: `/spec` asks instead of inventing requirements
If the issue is ambiguous, `/spec` SHALL ask clarifying questions and route the answers to the issue body; it SHALL NOT invent requirements.

#### Scenario: Ambiguous issue
- **WHEN** `/spec` cannot derive clear requirements from the issue
- **THEN** it asks, and does not finalize the plan until the issue body is clarified.

