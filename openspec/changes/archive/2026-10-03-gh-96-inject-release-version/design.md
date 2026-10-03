## Context

`cli/cmd/octospec/main.go` declares `const version = "1.1.0"` and the `version`
command prints it. `.goreleaser.yaml` builds the binary without any `ldflags`,
so the tag never reaches the binary. The `release-acceptance` workflow runs
`octospec version` and asserts the output contains the tag, which fails when the
const drifts from the tag.

## Goals / Non-Goals

**Goals:**
- The release binary reports the tag it was built from.
- A plain `go build` still works and prints a fallback.

**Non-Goals:**
- No change to the `version` command's output format beyond the value.
- No change to the installer or the workflows.

## Decisions

### A package-level `var version = "dev"` with ldflags injection
Change `const version` to `var version = "dev"`. GoReleaser sets it with
`-ldflags "-X main.version={{ .Version }}"`. `{{ .Version }}` is already the tag
without the leading `v` (for example `1.0.0`), so no string trimming is needed.
The alternative, a generated source file, adds a build step for no gain; the
alternative, keeping the const and bumping it by hand each release, is what
caused the bug.

### Keep the fallback as `dev`
A local `go build` has no tag, so it prints `octospec dev`. This is honest about
a non-release build and avoids pretending to be a version that was never tagged.

## Risks / Trade-offs

- [A release built outside GoReleaser would report `dev`] -> The release workflow
  is the only supported path and always passes the ldflags, so this is correct
  by design.
- [The ldflags path must match the variable's full name] -> The package is
  `main`, so `-X main.version=...` is exact; the test asserts the injected value
  is what `version` prints.

## Migration Plan

No data migration. The const stays source-compatible as a var; the only behavior
change is that the release now reports the tag. Rollback is a revert.

## Open Questions

- None. The issue is specific: inject the tag, keep a fallback, cover it.
