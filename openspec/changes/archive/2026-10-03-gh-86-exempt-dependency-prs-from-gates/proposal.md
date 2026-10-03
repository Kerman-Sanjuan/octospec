## Why

Dependency PRs open by Dependabot cannot land because the G3 branch-name and G4
issue-link gates reject the branch and body Dependabot generates. All five open
dependency PRs (#73-#77) are stalled on this, even though the bumps themselves
are fine.

## What Changes

- Skip the G3 branch-name check when the pull request author is `dependabot[bot]`.
- Skip the G4 issue-link check when the pull request author is `dependabot[bot]`.
- Keep G2a, G2b, G7 and G8 running on dependency PRs, so structure and
  uniqueness are still enforced.
- Add a bot-exemption case to `scripts/check-gates-test.sh`.

## Capabilities

### New Capabilities

### Modified Capabilities
- `gates`: G3 and G4 gain a bot exemption for `dependabot[bot]`; every other
  gate keeps its behaviour.

## Impact

- `scripts/check-gates.sh`: read the PR author and skip G3 and G4 for
  `dependabot[bot]`.
- `.github/workflows/openspec.yml`: pass the PR author to the gate script.
- `scripts/check-gates-test.sh`: cover the bot path.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/86

## Changelog

### Fixed
- Dependency PRs from Dependabot no longer fail the G3 and G4 gates.
