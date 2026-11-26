# Tasks — stage result

## 1. The object

- [ ] 1.1 `result.Result` with schema, stage, issue, change, outcome, reason, artifacts, duration and
  instant.
- [ ] 1.2 Closed vocabularies: stages, outcomes and reasons, exported.
- [ ] 1.3 `Validate()` refuses an unknown stage, an unknown outcome, a stalled result without a
  reason, an `ok` with a reason, and an `ok` from `ship` or `archive` with no artefact.

## 2. The files

- [ ] 2.1 `Write(dir, res)` writes `<stage>-<sha>.json` atomically.
- [ ] 2.2 `Read(path)` parses and validates.
- [ ] 2.3 A result without a commit falls back to `<stage>.json`, and says so.

## 3. The command

- [ ] 3.1 `octospec result write --stage --issue --outcome [--reason] [--artifact] [--duration-ms]`,
  with `--json`.
- [ ] 3.2 `octospec result show <file>` prints it and exits non-zero when it does not validate.

## 4. The schema

- [ ] 4.1 `docs/schema/result.schema.json` committed.
- [ ] 4.2 Test: the schema's enums equal the Go vocabularies, field by field.

## 5. The stages

- [ ] 5.1 `apply`, `ship`, `review` and `archive` end by writing their result.
- [ ] 5.2 `docs/commands.md` documents the command and the vocabulary.

## 6. Verification

- [ ] 6.1 `go test ./...` and `scripts/e2e-test.sh` pass.
- [ ] 6.2 Two runs of the same failing stage produce the same reason (test).
