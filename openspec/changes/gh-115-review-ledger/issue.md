# Issue #115 — One machine-readable record per pull request

https://github.com/Kerman-Sanjuan/octospec/issues/115

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As a maintainer, I want each pull request to leave a record behind with its gate verdicts and the
risk its review found, so that I can tell what we have been merging without a metrics platform.

## Context

The gates print text and the review publishes a comment: enough to block a merge, useless to answer
"how much risk did we merge in November?". A label cannot answer it either — a label is current
state, and it disappears with the pull request.

## Requirements

- `octospec review --ledger <dir>` writes `<pr>.json`, rewritten on every run.
- Fields: `schema`, `pr`, `issue`, `head`, `score`, `threshold`, `risk`, `risk_why`, `verdict`,
  `findings`, `at`.
- `risk_why` is stored, not only `risk`.
- `octospec ledger <dir>` prints the table plus counts by risk and verdict and the mean score.
- The record is a file in the tree, so it survives the merge.
- A high-risk review leaves `risk:high` on the pull request and clears it when it is no longer high.

## Success criteria

- Two pull requests, two files; re-running rewrites one.
- `octospec ledger` needs no network.
- The JSON and the printed review agree.

## Out of scope

- Dashboards, DORA, cross-repository trends, anything server-side.
