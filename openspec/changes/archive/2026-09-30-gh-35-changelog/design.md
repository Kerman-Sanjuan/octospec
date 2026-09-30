## Context

`goreleaser` builds and publishes the binaries on a `v*` tag, and the `curl | sh` installer reads the latest release. The missing piece is the human side: what changed, and what the maintainer does to cut a release.

## Goals / Non-Goals

**Goals:**
- A changelog a reader can scan.
- A short, repeatable release procedure.

**Non-Goals:**
- Automating the changelog from commits.

## Decisions

- **A hand-written changelog grouped by impact.** Commits are for us; the changelog is for consumers. `git log` buries what matters. *Alternative rejected:* generating it from commit messages.
- **Semantic versioning, tag drives the version.** The tag is the source of truth; the version is not hand-edited in scattered files.

## Risks / Trade-offs

- [The changelog falls behind] → The AGENTS.md definition of done can require an entry when behaviour changes.

## Migration Plan

None.

## Open Questions

- None.
