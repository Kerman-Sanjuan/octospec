## Context

The `cli` job runs `gofmt`, `go vet`, `go test`, `go build`, and two shell tests. `golangci-lint` is not installed locally, so its config is verified by CI.

## Goals / Non-Goals

**Goals:**
- A linter, the race detector, and coverage run in the `cli` job.

**Non-Goals:**
- A coverage threshold, at least at first.

## Decisions

- **`golangci-lint` with a conservative, explicit linter set.** `staticcheck`, `ineffassign`, and `unused` are high-signal; the noisy defaults stay off, so existing code is not drowned in style findings. *Alternative rejected:* the full default set, which would flag style before the project has agreed on it.
- **`go test -race -coverprofile` in one step, with a coverage summary line.** One run keeps CI fast. *Alternative rejected:* a separate coverage job, which duplicates the test run.
- **The config lives in `cli/.golangci.yml`.** The action runs with `working-directory: cli`, so it reads that file.

## Risks / Trade-offs

- [The linter flags existing code] -> The set is small; if it does, fix the finding rather than widen the ignore list.
- [`-race` slows the tests] -> The suite is small; the value outweighs the seconds.

## Migration Plan

None.

## Open Questions

- A coverage threshold, later.
