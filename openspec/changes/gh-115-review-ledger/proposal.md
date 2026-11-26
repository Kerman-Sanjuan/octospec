# Proposal — review ledger

## Why

Everything this tool learns about a pull request dies with it. The `## Review` comment says "10.0,
risk high, needs a human" and then the pull request closes and the fact goes with it. Labels are no
better: they are current state, they stack up, and they vanish.

The question Kerman asked is the right one — *what risk have we been merging?* — and it cannot be
answered from artefacts that evaporate.

## What Changes

- `octospec review --ledger <dir>` writes one file per pull request, `<pr>.json`, with the score, the
  risk, **why** the risk is what it is, the verdict, the findings count and the commit.
- Re-running rewrites the same file: the key is the pull request, so a re-review never duplicates.
- `octospec ledger <dir>` prints the series: a table, counts by risk and by verdict, and the mean
  score. Files only, no service.
- A high-risk review puts `risk:high` on the pull request, so the list itself says where to look; a
  later review that finds it lower removes it. The label and the file never disagree.

One file per pull request instead of one append-only file on purpose: each branch writes only its
own record, so two open pull requests never conflict on the same line.

## Capabilities

- `review-ledger`: added — the record, its fields and the summary.
- `label-provisioning`: added — the `risk:high` label is provisioned with the rest.

## Impact

- `cli/internal/review/` — the record and the summary.
- `cli/cmd/octospec/main.go` — `review --ledger`, `ledger` subcommand.
- `cli/internal/seed/` — the `risk:high` label.
- `cli/internal/payload/commands/review.md` — write the record, move the label.
- `docs/commands.md`.

## Changelog

### Added

- `octospec review --ledger <dir>` writes one record per pull request, and `octospec ledger <dir>`
  summarises them: counts by risk and verdict, and the mean score.
- `/review` moves the `risk:high` label so a pull request says, at a glance, whether its review asked
  for a human.
