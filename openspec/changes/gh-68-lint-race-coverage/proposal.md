## Why

The `cli` job runs `gofmt`, `go vet`, and `go test`. It misses a linter, the race detector, and coverage, so some regressions slip to review.

## What Changes

- Add `golangci-lint` to the `cli` job with a conservative config.
- Run `go test -race ./...` and report coverage.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `release-pipeline`: the `cli` job runs a linter, the race detector, and coverage.

## Impact

- `cli/.golangci.yml` and `.github/workflows/cli.yml`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/68
