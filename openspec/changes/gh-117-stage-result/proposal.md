# Proposal — stage result

## Why

Every stage ends with prose: "the pull request is open", "the archive failed because the lock was
held". A person reads that fine. A scheduler cannot, so today the only machine-readable thing a stage
leaves is a comment id in `.openspec.yaml`.

The expensive part is not the JSON. It is `reason`: a timeout, a crash, an empty diff and a change
deliberately parked need different retries, and an outcome of `failed` for all four teaches the
control plane to retry a change that will never work, or to give up on one that would.

## What Changes

- `octospec result write` writes `.octospec/runs/<stage>-<sha>.json` with the schema version, the
  stage, the issue, the change, the outcome, the reason, the artefacts and the duration; `--json`
  prints the same object.
- `octospec result show <file>` reads one back and fails when it does not validate.
- Vocabularies are closed and exported: three outcomes, six reasons, six stages. A typo fails here
  instead of becoming a state nobody handles.
- `ok` must be able to point at something: `ship` without a pull request or `archive` without an
  archived change is refused, because an `ok` with nothing attached is the result a broken stage
  writes.
- A JSON Schema file is committed under `docs/schema/`, and the test compares its enums with the
  Go vocabularies so the two cannot drift.

## Capabilities

- `stage-result`: added — the object, its vocabulary, its file, and the schema consumers validate
  against.

## Impact

- `cli/internal/result/` — the object, the validation, the writer and the reader.
- `cli/cmd/octospec/main.go` — `result write`, `result show`.
- `docs/schema/result.schema.json`, `docs/commands.md`.
- `cli/internal/payload/commands/{apply,ship,review,archive}.md` — each ends by writing its result.

## Changelog

### Added

- `octospec result`: every stage can leave a machine-readable record of what happened and why, with a
  closed vocabulary for reasons and a JSON Schema for consumers.
