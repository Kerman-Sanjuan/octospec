# Issue #113 — /review: score the pull request and gate on risk before merge

https://github.com/Kerman-Sanjuan/octospec/issues/113

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As a reviewer, I want the machine to score a pull request and flag the risky ones, so that I read
only the ones that need me and the rest merge on evidence.

## Context

Today the human reads every PR. G7 and G8 only warn, and G5 says whether the build is green, not
whether the change is any good. The pattern that works elsewhere is layered: deterministic checks
that block, a scored review, and a risk assessment that decides who has to look. The model detects
findings; a script converts them into a score, so the number is auditable and the same findings
always give the same score.

## Requirements

- `/review` SHALL publish exactly one comment per pull request, updated in place on every re-run,
  with: score 0-10, the findings that cost points (file:line), and the risk level.
- The score SHALL be computed by a deterministic script from the findings, never by the model
  grading itself.
- The five axes SHALL be explicit: correctness, architecture, domain fit, tests, and UI/style where
  the change has a UI surface.
- Risk SHALL combine at least: files inside a sensitive area (declared per repository in config),
  diff size, and whether the author has never touched those paths before.
- The merge gate SHALL block at or below a configurable threshold. The threshold and the
  sensitive-area list SHALL live in the adopting repository's config, never in octospec.
- The comment SHALL state that its content comes from an agentic process.

## Success criteria

- A PR with a deliberately weak test fails the gate, and the finding is cited by file and line.
- A clean PR passes with a score above the threshold and no human read.
- Changing the threshold requires editing config only, with no CLI release.

## Out of scope

- Replacing G5. A green build is a precondition, not the decision.
- Reviewing code the human then has to re-read anyway: if the score cannot be defended, the gate is
  wrong.
- Cross-repository learnings of previous reviews.
