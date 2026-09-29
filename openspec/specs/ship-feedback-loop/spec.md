# ship-feedback-loop Specification

## Purpose
How `/ship` runs or observes CI and hands a failing check back to `/apply`, only opening or updating the pull request once the gates are green.
## Requirements
### Requirement: `/ship` observes CI
`/ship` SHALL run or observe CI before opening or updating the pull request.

#### Scenario: Ship checks CI
- **WHEN** `/ship` runs
- **THEN** it runs or observes the gates and records the result.

### Requirement: Failures loop back to `/apply`
When CI fails, `/ship` SHALL hand the specific failure back to `/apply` to fix and iterate, instead of only reporting a failure.

#### Scenario: Iterate until green
- **WHEN** CI fails
- **THEN** `/ship` returns the specific failing check to `/apply`, which fixes it, and `/ship` re-checks.

### Requirement: Proceed only on green
`/ship` SHALL open or update the pull request only once the gates are green.

#### Scenario: Green CI
- **WHEN** all gates pass
- **THEN** `/ship` opens or updates the pull request.

