# Issue #121 — Enforce hard budgets on every unattended run

https://github.com/Kerman-Sanjuan/octospec/issues/121

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As a maintainer, I want a run to stop by itself when it exceeds its budget, so that an unattended
agent cannot turn into a surprise bill or a 300-file pull request.

## Context

Nothing stops a long run today: no cap on turns, tokens, wall time, files touched or diff size. In a
session a human notices; in a Kubernetes Job nobody does until the invoice.

## Requirements

- A budget block with `max_minutes`, `max_tokens`, `max_files`, `max_diff_lines` and `max_tasks`.
- When a budget is hit, the stage stops cleanly, commits what is complete, and returns
  `stalled / budget_exhausted` naming the budget that was hit.
- Budgets are checked between tasks, never mid-task.
- Tokens come from the agent's own report; when it cannot report them, the budget reads `unknown`,
  never "fine".
- The defaults are conservative and documented.

## Success criteria

- A run that hits `max_files` stops with the reason and leaves a resumable branch.
- A run under budget is unaffected.
- `octospec doctor` prints the effective budgets.

## Out of scope

- Cost in currency.
