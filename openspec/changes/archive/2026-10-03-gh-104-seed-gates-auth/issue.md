# Issue #104: The seeded gates workflow cannot read the issue

https://github.com/Kerman-Sanjuan/octospec/issues/104

## User story

As a consumer who seeds octospec into my repository, I want the gates to read
the linked issue, so that G1 does not fail on an issue that already has every
required section.

## Context / problem

The root `.github/workflows/openspec.yml` authenticates `gh` (`GH_TOKEN:
${{ github.token }}`, `issues: read`, `setup-node@v7`), but the copy the CLI
seeds into consumer repositories does not. Commit `6701cff` ("ci(#60):
authenticate gh in the gates job") changed only the root workflow and never
updated the seed.

On a freshly seeded repository, `gh` is on the runner but unauthenticated, so
`gh issue view` fails, the issue body is empty, and G1 reports a false negative
on an issue that has all four sections:

```
FAIL G1 issue #1 missing section(s): User story Context Requirements Success criteria
```

That is the worst failure mode: a gate blaming the user for something they did
not do, on the first `/spec` of the workflow the README tells them to follow. No
test catches it: `check-gates-test.sh` injects a fake `gh` that always succeeds,
and the e2e runs the gates with no change in the diff, so G1 never runs.

## Requirements

- The seeded `.github/workflows/openspec.yml` SHALL authenticate `gh` with
  `GH_TOKEN: ${{ github.token }}` and request `issues: read`, matching the root.
- The seeded workflow SHALL use the same `actions/setup-node` major as the root.
- `scripts/check-gates.sh` SHALL NOT fail G1 when `gh` cannot read the issue
  body (auth, network, or a missing issue); it SHALL SKIP with a clear message,
  so a false negative is impossible.
- The seeded `check-gates.sh` SHALL stay byte-identical to the root copy.
- A gate test SHALL cover the unauthenticated-`gh` path and assert G1 does not
  FAIL.
- A test SHALL assert the seeded workflow keeps the token, the `issues: read`
  permission, and the same `setup-node` major as the root.

## Success criteria

- A freshly seeded repository runs the gates and G1 passes on a complete issue.
- `PATH` with a failing `gh`, G1 reports SKIP, not FAIL, and the run passes.
- `go test ./...` includes a seed test that fails if the workflow drifts again.
- `openspec validate --all --strict`, `scripts/check-gates.sh`,
  `scripts/check-gates-test.sh`, and CI pass.

## Out of scope

- No change to the gate set or its rules beyond G1's failure handling.
- No new CLI features.

