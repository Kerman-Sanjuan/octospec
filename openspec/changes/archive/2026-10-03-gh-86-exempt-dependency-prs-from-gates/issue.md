# Issue #86: Exempt automated dependency PRs from the G3 and G4 gates

https://github.com/Kerman-Sanjuan/octospec/issues/86

## User story

As a maintainer of this repository, I want the gates to exempt automated
dependency PRs, so that Dependabot updates can land without hand-editing branch
names and PR bodies.

## Context / problem

All five open dependency PRs (#73-#77) fail the `gates` check. G3 rejects the
`dependabot/...` branch name and G4 rejects Dependabot's auto-generated body for
lacking `Closes #n`. The dependency bumps themselves are fine; the convention
checks are what block them. This stalls security and tooling updates until
someone opens a matching issue and a `feat|fix` branch by hand.

## Requirements

- The gates SHALL skip the G3 branch-name check when the PR author is
  `dependabot[bot]`.
- The gates SHALL skip the G4 issue-link check when the PR author is
  `dependabot[bot]`.
- The gates MUST still run G2a, G2b, G7 and G8 on dependency PRs.
- The gate test script MUST cover the bot-exemption path.

## Success criteria

- PRs #73-#77 pass the `gates` check and can merge once `cli` is green.
- A non-bot PR with a bad branch name or a missing `Closes #n` still fails
  G3/G4.
- `sh scripts/check-gates-test.sh` passes and includes a bot-exemption case.

## Out of scope

- Changing the Dependabot schedule or ecosystems.
- Auto-merge or auto-close of bot PRs.
- The G1 and G6 spec-delta checks.
- Retroactively rewriting the existing Dependabot PRs.
