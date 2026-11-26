# Proposal — budgets

## Why

An unattended run has no floor and no ceiling. In a session somebody notices at minute ninety; in a
Job nobody notices until the invoice, and by then the pull request touches three hundred files.

`#117` already taught a stage to say `budget_exhausted`. Nothing has ever produced that reason,
because nothing measures anything.

## What Changes

- `.octospec/budget.json` holds per-stage limits: `max_minutes`, `max_tokens`, `max_files`,
  `max_diff_lines`, `max_tasks`. Defaults are conservative and written by `seed`.
- `octospec budget --stage <s>` takes what the run has consumed and answers one question: is there
  budget left. It exits 6 when there is not, so a wrapper can stop without parsing English.
- `octospec budget --plan` prints the limits so a runner can set a timeout before it starts.
- Tokens are the agent's own report. `--tokens -1` means "the agent cannot say", and the answer says
  `unknown` instead of quietly passing.
- `doctor` prints the effective budgets, so a repository cannot pretend it has them.
- `/apply` and `/ship` check between tasks, commit what is complete and write
  `stalled / budget_exhausted`.

## Capabilities

- `budgets`: added — the file, the defaults, the check and the exit code.

## Impact

- `cli/internal/budget/` — limits, usage, the check.
- `cli/cmd/octospec/main.go` — `budget`.
- `cli/internal/seed/` — the default file.
- `cli/internal/doctor/` — the effective budgets.
- `cli/internal/payload/commands/{apply,ship}.md`, `docs/commands.md`.

## Changelog

### Added

- `octospec budget`: per-stage limits for time, tokens, files, diff lines and tasks, checked between
  tasks and reported as `stalled / budget_exhausted` when they run out.
