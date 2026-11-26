## ADDED Requirements

### Requirement: Every review leaves a record

`octospec review --ledger <dir> --pr <n>` SHALL write `<dir>/<n>.json` with the schema version, the
pull request, the issue, the head commit, the score, the threshold, the risk, the reasons for that
risk, the verdict, the number of findings and the instant. It SHALL rewrite the same file on a
re-run instead of adding another.

#### Scenario: A review writes its record

- **WHEN** a review scores 10 with risk high and verdict human for pull request 129
- **THEN** `ledger/129.json` holds those values and the reasons the risk was high

#### Scenario: Re-running a review

- **WHEN** the same pull request is reviewed again
- **THEN** the same file is rewritten and the directory still holds one record for it

#### Scenario: No pull request number

- **WHEN** `--ledger` is given without `--pr`
- **THEN** the command fails and says a record needs a pull request to key on

### Requirement: The series can be read without a service

`octospec ledger <dir>` SHALL print every record as a row and a summary with the counts by risk, the
counts by verdict and the mean score. It SHALL read files only, and an empty directory SHALL be
reported as empty rather than as an error.

#### Scenario: Reading the series

- **WHEN** the directory holds three records, one low, one medium and one high
- **THEN** the summary counts one of each and prints the mean of their scores

#### Scenario: A broken file

- **WHEN** one file in the directory is not valid JSON
- **THEN** the command skips it, reports how many it skipped, and still summarises the rest
