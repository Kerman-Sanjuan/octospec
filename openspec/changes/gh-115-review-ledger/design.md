# Design — review ledger

## Context

A review produces a number, a risk and a verdict. All three are printed and published in a comment,
and all three are lost when the pull request closes. Any question starting with "how many" or "how
much risk" is unanswerable today.

## Goals / Non-Goals

**Goals**

- One record per pull request, in the tree, surviving the merge.
- Cheap to write, cheap to read, no service, no database.
- Honest about *why* the risk was high, not only that it was.

**Non-Goals**

- A dashboard. `octospec ledger` prints a table and five numbers; that is the ceiling.
- Cross-repository aggregation. `find` and `jq` do that better than we would.
- Replacing the comment. A person reads the comment; the file is for questions about many pull
  requests.

## Decisions

### One file per pull request, not an append-only log

`<pr>.json` means the branch writes exactly one file and two open pull requests never touch the same
line. An append-only `.jsonl` would conflict on every concurrent branch, and a merge conflict in a
ledger is the worst kind: silent data loss in the one place you wanted history.

Rewriting the same file also makes the operation idempotent: a re-review replaces its record instead
of adding a second opinion nobody asked for.

### `risk_why` is part of the record

A row that says `high` tells the future nothing. A row that says `high — touches
cli/internal/payload/commands` teaches whether the sensitive-path list is right, or whether it is a
list of everything we fear because we touched it once.

### The label mirrors the record, and is not the record

`risk:high` is a pointer for the eye in a list of pull requests. It is set and removed by the same
step that writes the file, so the two cannot drift. The history lives in the files because a label
cannot hold a series.

### The schema is versioned

`schema: 1` costs one field and means a later reader can tell an old record from a broken one.

## Risks

- A repository that never merges keeps the files on branches, so nothing records them. Acceptable:
  a pull request nobody merged is not a risk we took.
- The label and the file can drift if somebody labels by hand. The next review fixes it, which is the
  same self-healing property the rest of the flow has.
