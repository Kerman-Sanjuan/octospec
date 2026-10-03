## ADDED Requirements

### Requirement: `/archive` retries the push to `main` after a rebase
`/archive` SHALL push the archive to `main` with a bounded fetch-rebase-push
retry, so two independent archives do not lose to a race. If the rebase
conflicts, it SHALL stop and hand the conflict to the human.

#### Scenario: A losing archive retries and lands
- **WHEN** the push to `main` is rejected because another archive landed first
- **THEN** `/archive` fetches `main`, rebases, and pushes again.

#### Scenario: Rebase conflict stops the archive
- **WHEN** the rebase conflicts in a shared file such as `CHANGELOG.md` or a
  synced spec
- **THEN** `/archive` stops and reports the conflict instead of forcing a push.
