## ADDED Requirements

### Requirement: A repository declares its limits

`.octospec/budget.json` SHALL accept `max_minutes`, `max_tokens`, `max_files`, `max_diff_lines` and
`max_tasks`, globally and per stage, with conservative defaults when the file or a field is absent.

#### Scenario: No file

- **WHEN** a repository has no budget file
- **THEN** the defaults apply and `doctor` prints them as the budgets in force

#### Scenario: A stage overrides one limit

- **WHEN** the file sets `stages.apply.max_files` and nothing else
- **THEN** `apply` uses that number and every other limit takes its default

### Requirement: A run is told when it has no budget left

`octospec budget --stage <s>` SHALL compare what the run consumed with the limits and SHALL exit 6,
naming the limit that was hit, when one is exhausted. It SHALL exit 0 when there is budget left.

#### Scenario: Files exhausted

- **WHEN** a run has touched 26 files with `max_files` 25
- **THEN** the command names `max_files` and exits 6

#### Scenario: Under budget

- **WHEN** every measure is below its limit
- **THEN** the command exits 0 and prints what is left

#### Scenario: The runner sets a timeout before starting

- **WHEN** `--plan` is asked for
- **THEN** the command prints the limits without consuming anything

### Requirement: What cannot be measured is reported as unknown

Token usage SHALL come from the caller. When it cannot be reported, the check SHALL say `unknown` and
SHALL NOT block.

#### Scenario: An agent that cannot report usage

- **WHEN** the caller passes no token count
- **THEN** the answer says tokens are unknown and the verdict is still computed from the rest

### Requirement: The stop happens between tasks

A stage that runs out of budget SHALL finish the task in hand, commit it, and report
`stalled / budget_exhausted` naming the limit.

#### Scenario: A resumable stop

- **WHEN** `/apply` hits `max_files` after a task
- **THEN** it commits that task, writes the result with `reason: budget_exhausted`, and leaves the
  branch where the next run can continue
