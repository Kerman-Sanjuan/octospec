# Design — budgets

## Context

The result contract (#117) has a reason nobody can produce yet: `budget_exhausted`. Producing it needs
a policy, a measurement and a place to check both.

## Goals / Non-Goals

**Goals**

- One file per repository, five numbers per stage, one question answered.
- A stop that leaves the work resumable: commit what is done, stop between tasks.
- An honest `unknown` for what the tool cannot measure.

**Non-Goals**

- A hard kill. `octospec` is not the process supervisor: the wrapper that runs the agent kills it,
  and reads the limits from `--plan` to do it.
- Cost in currency.
- Counting anything from the filesystem at run time. The caller reports files and diff lines, because
  only the caller knows what the run has touched.

## Decisions

### `.octospec/budget.json`, not `.octospec.yaml`

The issue said `.octospec.yaml`. I am keeping every policy in one directory, JSON, next to
`.octospec/review.json`, because the tool already reads and seeds that shape and a second file format
for four numbers is a second thing to get wrong. The deviation is recorded here rather than silently.

### Defaults that are conservative on purpose

| limit | default | why |
|---|---|---|
| `max_minutes` | 45 | an interactive session notices long before; a Job does not |
| `max_tokens` | 250 000 | roughly two thirds of a strong model's context spent twice over |
| `max_files` | 25 | past that, a change is not one change |
| `max_diff_lines` | 800 | the review already calls 400 large |
| `max_tasks` | 12 | a change with more tasks is a change that should be split |

They are the numbers I would defend in review, not measured optima. First calibration comes from the
runs that hit them.

### Between tasks, never mid-task

Checking mid-task leaves a half-written file and a branch nobody can resume. The contract is: finish
the task, commit it, check, stop. A run that stops always stops at a boundary.

### `unknown`, not `fine`

`--tokens -1` is the agent saying it cannot report usage. The check then reports `tokens: unknown`
and does not block on it: blocking on a number we do not have would stop every run whose agent is
quiet, and pretending it is zero is how a budget becomes decoration.

### Exit 6

3 and 4 belong to the lock, 5 to the review verdict. 6 says "the run has no budget left", which a
wrapper acts on by stopping the job, not by retrying.

## Risks

- Self-reported usage can be understated. The wall clock is the backstop: `max_minutes` is measured
  against a start time the caller supplies and cannot be faked by the agent's optimism.
- A repository with no budget file gets the defaults, which may be too small for a big change. The
  fix is a file, and `doctor` says which budgets are in force.
