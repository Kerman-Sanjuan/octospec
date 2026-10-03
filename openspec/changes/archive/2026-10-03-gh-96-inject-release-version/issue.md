# Issue #96: Inject the release version at build time

https://github.com/Kerman-Sanjuan/octospec/issues/96

## User story

As a maintainer cutting a release, I want the binary to report the tag it was
built from, so that `octospec version` matches the release and the acceptance
workflow can verify it.

## Context / problem

`octospec version` prints a hardcoded `const version = "1.1.0"` in
`cli/cmd/octospec/main.go`. GoReleaser does not inject a version, so a release
built from a `v1.0.0` tag still reports `1.1.0`. The release spec already says
the tag is the source of truth and the version should not be hand-edited in
scattered files, but the code does exactly that. The `release-acceptance`
workflow asserts the version contains the tag, so the mismatch fails the
acceptance check. This blocks the 1.0.0 debut: tagging v1.0.0 today would
publish a binary that reports 1.1.0.

## Requirements

- The binary SHALL report the release tag when built by the release workflow.
- `octospec version` SHALL print the tag without the leading `v` (for example
  `octospec 1.0.0`).
- A local `go build` without the injected value SHALL still produce a runnable
  binary with a sensible version string (for example `octospec dev`).
- The version SHALL be injected at build time from the tag, not written by hand
  in a source file on each release.

## Success criteria

- `go build` with `-ldflags "-X main.version=vX.Y.Z"` makes `octospec version`
  print that version.
- The GoReleaser config injects the tag through `ldflags`.
- A plain `go build` prints the fallback and does not fail.
- `openspec validate --all --strict`, `scripts/check-gates.sh`, and the `cli`
  and `gates` CI checks pass.

## Out of scope

- No change to the workflow commands, agents, or the OpenSpec schema.
- No new CLI commands or flags.
