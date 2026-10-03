## 1. Extraction helper

- [x] 1.1 Add `scripts/release-notes.sh` to print a version's `CHANGELOG.md` section, with `v` prefix handling and a fallback.
- [x] 1.2 Add `scripts/release-notes-test.sh` covering the real section, the prefix, the fallback, and a synthetic changelog.

## 2. Release wiring

- [x] 2.1 Disable GoReleaser's native changelog in `.goreleaser.yaml`.
- [x] 2.2 In `.github/workflows/release.yml`, build the notes with `release-notes.sh` and pass `--release-notes`.
- [x] 2.3 Run `scripts/release-notes-test.sh` in the `cli` CI job.

## 3. Checks

- [x] 3.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, `sh scripts/check-gates.sh`, `sh scripts/release-notes-test.sh`, and `sh scripts/e2e-test.sh`.
