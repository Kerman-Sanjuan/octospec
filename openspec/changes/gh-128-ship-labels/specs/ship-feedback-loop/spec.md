## ADDED Requirements

### Requirement: The pull request carries the issue's kind

`/ship` SHALL copy every `type:*` label from the issue onto the pull request, and SHALL stop, naming
the reason, when the issue carries none.

#### Scenario: The kind is copied

- **WHEN** `/ship` opens or updates the pull request for an issue labelled `type:bug`
- **THEN** the pull request carries `type:bug`

#### Scenario: No kind to copy

- **WHEN** the issue has no `type:*` label
- **THEN** `/ship` stops with a reason naming the issue, and does not open the pull request

### Requirement: The pull request says whether a machine may own it

`/ship` SHALL copy `mode:autonomous` onto the pull request when the issue carries it, and SHALL NOT
copy `status:*` labels there.

#### Scenario: An autonomous issue

- **WHEN** the issue carries `mode:autonomous`
- **THEN** the pull request carries `mode:autonomous` too

#### Scenario: Status stays on the issue

- **WHEN** the gates pass and `/ship` adds `status:in-review`
- **THEN** the label lands on the issue, and the pull request's `status:*` labels are unchanged

#### Scenario: Saying it twice changes nothing

- **WHEN** `/ship` runs again on a pull request that already carries both labels
- **THEN** the label set is unchanged
