## ADDED Requirements

### Requirement: The gate script runs on a non-PR ref
`scripts/check-gates.sh` SHALL skip the PR-only branch gate when the ref is a long-lived branch, and SHALL report the skip instead of failing on it.

#### Scenario: Run on a long-lived branch
- **WHEN** `scripts/check-gates.sh` runs with `HEAD_REF=main`
- **THEN** it reports the branch gate as skipped and does not fail on it.

#### Scenario: Run on a PR branch
- **WHEN** `scripts/check-gates.sh` runs with a `feat|fix/<issue>-<slug>` ref
- **THEN** it checks the branch name and passes.
