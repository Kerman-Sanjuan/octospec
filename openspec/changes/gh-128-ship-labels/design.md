# Design — label the pull request

## Context

A pull request opened by a machine looks exactly like one opened by a person, except that ours
carry no labels at all. The issue says what kind of change it is; the pull request does not.

## Goals / Non-Goals

**Goals**

- One glance on the pull request list tells the kind and whether a machine may own it.
- No new vocabulary: the labels are the ones `seed` already provisions.

**Non-Goals**

- Inferring labels from the title or the diff. A guessed kind is worse than a missing one.
- Mirroring `status:*`. Status describes the issue's journey; the pull request has its own signals
  (checks, review verdict).

## Decisions

### Copy, never invent

`/ship` copies the issue's `type:*` and `mode:autonomous` verbatim. Anything else — new labels,
inferred kinds, a mapping from conventional-commit prefixes — would need a policy and a place to
keep it, and `ship` is not the place.

### A missing `type:*` stops the stage

The alternative is opening the pull request unlabelled and hoping somebody notices. Stopping with
"issue #N has no `type:*` label" costs one command and keeps the rule true.

### `status:in-review` stays on the issue

Adding it to the pull request as well would create two sources of truth for one fact. The pull
request already says it is under review by existing.

## Risks

- A repository that named its kinds differently gets nothing copied. It will notice on the first
  pull request, which is the right moment.
