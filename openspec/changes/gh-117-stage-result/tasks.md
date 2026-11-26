# Tasks — stage result

## 1. The object

- [x] 1.1 `result.Result` with schema, stage, issue, change, outcome, reason, artifacts, duration and
  instant.
- [x] 1.2 Closed vocabularies: stages, outcomes and reasons, exported.
- [x] 1.3 `Validate()` refuses an unknown stage, an unknown outcome, a stalled result without a
  reason, an `ok` with a reason, and an `ok` from `ship` or `archive` with no artefact.

## 2. The files

- [x] 2.1 `Write(dir, res)` writes `<stage>-<sha>.json` atomically.
- [x] 2.2 `Read(path)` parses and validates.
- [x] 2.3 A result without a commit falls back to `<stage>.json`, and says so.

## 3. The command

- [x] 3.1 `octospec result write --stage --issue --outcome [--reason] [--artifact] [--duration-ms]`,
  with `--json`.
- [x] 3.2 `octospec result show <file>` prints it and exits non-zero when it does not validate.

## 4. The schema

- [x] 4.1 `docs/schema/result.schema.json` committed.
- [x] 4.2 Test: the schema's enums equal the Go vocabularies, field by field.

## 5. The stages

- [x] 5.1 `apply`, `ship`, `review` and `archive` end by writing their result.
- [x] 5.2 `docs/commands.md` documents the command and the vocabulary.

## 6. Verification

- [x] 6.1 `go test ./...` and `scripts/e2e-test.sh` pass.
- [x] 6.2 Two runs of the same failing stage produce the same reason (test).
