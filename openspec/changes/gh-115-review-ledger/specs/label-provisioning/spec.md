## ADDED Requirements

### Requirement: The risk label is provisioned

`seed` SHALL provision the `risk:high` label, and `/review` SHALL add it to a pull request whose risk
is high and remove it when the risk is not high, so the label and the record cannot disagree.

#### Scenario: Seeding the label

- **WHEN** `octospec seed` runs in a repository without `risk:high`
- **THEN** the label exists afterwards, and running it again changes nothing

#### Scenario: A high-risk review

- **WHEN** `/review` finds risk high
- **THEN** the pull request carries `risk:high`

#### Scenario: A later review lowers it

- **WHEN** a later review on the same pull request finds risk medium
- **THEN** `risk:high` is removed from the pull request
