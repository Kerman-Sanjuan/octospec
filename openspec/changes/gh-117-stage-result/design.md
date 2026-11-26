# Design — stage result

## Context

A scheduler that reads comments is a scraper. A scheduler that reads `{outcome, reason}` is a control
plane. The difference is not the format, it is the vocabulary: what a stage is allowed to say about
why it stopped.

## Goals / Non-Goals

**Goals**

- One object per execution, in a file that survives the run, with a stable schema.
- A closed `reason` list, small enough that every value has a different retry.
- Vocabularies that cannot drift between the code and the schema consumers validate against.

**Non-Goals**

- Writing the result for the interactive stages (`/idea`, `/bug`): a person is in the session, and
  their answer is the result.
- Tokens and cost. They belong to the budget work (#121), which will add a field.
- Any aggregation. One record per execution is the contract.

## Decisions

### Two outcomes, six reasons

`ok` and `stalled` answer "did it move"; the reason answers "what now". Retry policy:

| reason | retry |
|---|---|
| `timeout` | yes, with a longer budget |
| `crashed` | yes, once; then a human |
| `empty_diff` | no: the stage had nothing to do |
| `no_pr_requested` | no: parking was the point |
| `needs_human` | not by a machine |
| `budget_exhausted` | yes, after raising the budget |

Collapsing these into `failed` is how an autonomous loop burns money retrying the same dead end.

### The file is keyed by stage and commit

`<stage>-<sha>.json` means a re-run of the same stage on the same commit rewrites its record, while
a new commit gets a new file. That is the shape a scheduler wants: "what happened last time I ran
apply on this commit".

### `ok` must point at something

`ship` with `ok` and no artefact, or `archive` with `ok` and nothing archived, is the exact shape of a
stage that returned zero because it did nothing. The check is not clever: `ok` from those two stages
requires at least one artefact, and the artefact is a URL or a path a person can open.

### The JSON Schema is a contract, and the test compares it with the code

`docs/schema/result.schema.json` is what a consumer in another language validates against. A test
reads that file and compares its `enum` lists with the Go vocabularies, so adding a reason in one
place and forgetting the other fails the build. We do not ship a generic validator: our own checks are
ten lines, and a validator would be a dependency for a schema we control.

## Risks

- A stage that forgets to write its result is invisible to the scheduler. The payload says so in each
  command, and #118 makes the stages resumable, where the missing result will surface as a hard error.
- `duration_ms` comes from the caller. A stage that lies about its duration is a stage that will lie
  about its reason; nothing we can check from here.
