# Proposal — review score

## Why

The human is still the bottleneck at merge: G5 says the build is green, not that the change is
good, and G7/G8 only warn. A scored review with a risk assessment lets the machine say "this one
needs you" and "this one does not", which is the precondition for merging without reading.

## What Changes

- A `review` command that turns a list of findings into a score 0-10, a risk level and a verdict.
- The arithmetic lives in Go (`cli/internal/review`), so the same findings always produce the same
  score in a laptop and in CI.
- The policy lives in the repository, in `.octospec/review.json`: threshold, weights per finding
  type, sensitive paths and what counts as a large diff. octospec ships defaults and never owns
  the numbers.
- `/review` publishes one comment per pull request, updated in place, with the score, the costed
  findings and the risk, and states that it comes from an agentic process.
- First release is report-only: the command always reports, and blocks only with `--enforce`, so a
  threshold can be set from a real distribution instead of a guess.

## Capabilities

**New capabilities**

- `review-score` — `specs/review-score/spec.md`

**Modified capabilities**

None.

## Impact

- New package `cli/internal/review`, a `review` subcommand, a new `/review` stage command in the
  payload (the tool payload grows from 6 to 7 commands), and `.octospec/review.json` seeded with
  defaults.
- Affected code: `cli/cmd/octospec/main.go`, `cli/internal/payload`, `cli/internal/seed`.
- No change to the gates: a repository wires `octospec review --enforce` where it wants the block.

## Changelog

### Added

- `octospec review`: score a pull request from a list of findings and decide whether it needs a
  human.
- `/review`: publish that score on the pull request, as one comment updated in place.
