# Issue #117 — Emit a machine-readable result from every stage

https://github.com/Kerman-Sanjuan/octospec/issues/117

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As the scheduler, I want one JSON object per stage execution, so that I can decide whether to retry,
park, or move on without parsing an agent's prose.

## Context

The stages end with prose. A scheduler needs `{stage, outcome, reason, artifacts, duration}` with the
same vocabulary everywhere, or every consumer infers state from comments and labels.

The `reason` vocabulary is the important part: a timeout, a crash, an empty diff and a change stopped
on purpose need different retries, and collapsing them into "failed" is how a control plane starts
lying.

## Requirements

- Every non-interactive stage SHALL write a result object, by default to
  `.octospec/runs/<stage>-<sha>.json`.
- Fields: `schema`, `stage`, `issue`, `change`, `outcome` (`ok | stalled`), `reason`, `artifacts[]`,
  `duration_ms`, `at`.
- `reason` is a closed list: `timeout`, `crashed`, `empty_diff`, `no_pr_requested`, `needs_human`,
  `budget_exhausted`.
- The schema is versioned, and a JSON Schema file is committed for consumers.
- `--json` prints the object, so a scheduler never reads the file.
- A result SHALL never contradict its artefacts: `ok` from `ship` without a pull request, or from
  `archive` without an archived change, is invalid.

## Success criteria

- The emitted object validates against the committed JSON Schema, and the test compares the schema's
  vocabularies with the code's, so they cannot drift.
- Two runs of the same failing stage produce the same `reason`.
- A consumer tells a timeout from an empty diff without reading the comment.

## Out of scope

- Aggregation and dashboards.
