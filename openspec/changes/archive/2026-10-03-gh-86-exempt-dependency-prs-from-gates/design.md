## Context

`scripts/check-gates.sh` runs G1 to G8 against a pull request. G3 builds the
branch name from `HEAD_REF` and fails anything that is not
`feat|fix/<issue>-<slug>`. G4 reads `PR_BODY` and fails when it has no
`Closes`/`Fixes`/`Resolves #n`. Dependabot opens pull requests on a branch named
`dependabot/...` with an auto-generated body that has neither, so every
dependency pull request fails the `gates` check even when it is otherwise clean.
The workflow `.github/workflows/openspec.yml` already passes `GH_TOKEN`,
`BASE_REF`, `HEAD_REF` and `PR_BODY` to the script; it does not pass the author.

## Goals / Non-Goals

**Goals:**
- Skip G3 and G4 for pull requests opened by `dependabot[bot]`.
- Keep G1, G2a, G2b, G6, G7 and G8 unchanged for every pull request.
- Cover the exemption in `scripts/check-gates-test.sh`.

**Non-Goals:**
- No change to the Dependabot schedule or ecosystems.
- No auto-merge or auto-close.
- No change to the branch-name or issue-link rules for human pull requests.

## Decisions

### Read the author from a new `PR_AUTHOR` env var
G3 and G4 already read `HEAD_REF` and `PR_BODY` from the environment, so the
author follows the same pattern: add `PR_AUTHOR`, set it in the workflow from
`github.event.pull_request.user.login`, and default it to empty. The script
already reads env with `:-` defaults, so no author means no exemption and the
current behaviour holds. The alternative, asking `gh` for the author inside the
script, adds a network call and a second source of truth when the workflow
already has the value.

### Skip only G3 and G4, not the whole script
The exemption is two `case` guards that report `SKIP` and skip the check, not an
early exit. G2a and G2b still validate the repository, and G7 and G8 still
inspect unarchived changes, which is what the issue asks for. The alternative,
an early exit for bots, would let a malformed change ride in on a dependency
pull request.

### Detect the bot by the exact login
The guard matches `dependabot[bot]`, not a prefix. A broader match on
`dependabot` or a `*[bot]` glob would exempt other automation (release bots,
GitHub Actions bots) that was not part of this issue. Widening later is a
one-line change if it is ever wanted.

### Test the exemption hermetically
`check-gates-test.sh` already runs the script with fake `BASE_REF`, `HEAD_REF`
and a fake `gh` on `PATH`. The bot case follows the same shape: run with
`PR_AUTHOR=dependabot[bot]` and a non-conforming branch, then assert G3 and G4
report `SKIP` instead of `FAIL`, and a second run without `PR_AUTHOR` asserts
they still `FAIL`.

## Risks / Trade-offs

- [Exempting the bot could let a broken dependency change through] -> Only G3
  and G4 are exempt. G2a and G2b still validate structure, and the `cli` check
  still compiles and tests the bump, so a bad dependency still fails.
- [A hardcoded login goes stale if GitHub renames the bot] -> The login has been
  stable for years; if it changes, the fix is one string and the new test fails
  loudly, which surfaces it.
- [Adding an env var the workflow can forget to set] -> The script defaults it
  to empty, so a missing value degrades to today's behaviour, not a crash.

## Migration Plan

No data migration. Land the script change and the workflow change in the same
pull request. Rollback is a revert: the existing five dependency pull requests
stay blocked until the revert is undone.

## Open Questions

- None.
