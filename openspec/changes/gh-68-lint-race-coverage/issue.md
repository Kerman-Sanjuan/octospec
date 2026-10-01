# Issue #68: Add golangci-lint, go test -race, and coverage to CI

https://github.com/Kerman-Sanjuan/octospec/issues/68

## User story

As a maintainer, I want stronger automated checks, so a regression is caught before merge.

## Context / problem

The `cli` job runs `gofmt`, `go vet`, and `go test`. It does not run a linter, the race detector, or coverage.

## Requirements

- CI SHALL run a Go linter such as `golangci-lint`.
- CI SHALL run `go test -race ./...`.
- CI SHALL report coverage.

## Success criteria

- The linter, the race detector, and coverage run on every pull request.
- A lint or race failure blocks the pull request.

## Out of scope

- A coverage threshold, at least at first.
