# Issue #98: Fix the audit findings

https://github.com/Kerman-Sanjuan/octospec/issues/98

## User story

As a maintainer, I want the CI and the docs to tell the truth, so that a green
check proves something and a newcomer who follows the README does not hit a
broken setup.

## Context / problem

An audit of the repository found the following. The first two are release
blockers tracked elsewhere; this issue covers the rest.

- The `cli` CI job mounts `setup-node` but never installs the OpenSpec CLI, so
  the `e2e` step prints `SKIP: openspec is not on PATH` and exits 0. The green
  check proves nothing.
- The README says "OpenSpec CLI 1.3.1 or newer", but with OpenSpec 1.14.0
  `openspec validate --strict` fails because `.openspec.yaml` carries a
  `github:` key the strict schema rejects. The workflows pin 1.3.1, which hides
  this from CI.
- The README prose at line 234 is mutilated: an open backtick and a line that
  begins with `;`.
- `openspec validate --archived` fails on three archived changes (gh-30, gh-31,
  gh-48) that carry unchecked checkboxes, which is exactly what gate G7 forbids.
- In `scripts/check-gates.sh`, G2b and G7 pass silently when `openspec` is
  missing, while G2a reports FAIL. The gates are inconsistent.
- The G6 failure message says "use --skip-specs", but the real control is
  `skip_specs: true` in `.openspec.yaml`.
- Four merged branches are still on the remote (66, 85, 92, 94), though
  `/archive` promises to delete them.

## Requirements

- The `cli` CI job SHALL install the OpenSpec CLI before the `e2e` step, so the
  step runs instead of skipping.
- The `e2e` test SHALL fail, not pass, when `openspec` is unavailable.
- The README SHALL state the OpenSpec version range the gates actually support,
  and the workflows SHALL install a version inside that range.
- The README prose at the skills paragraph SHALL be a complete, readable
  sentence with no stray backtick.
- `openspec validate --archived` SHALL pass: the archived changes gh-30, gh-31,
  and gh-48 SHALL have every task checked.
- G2b and G7 SHALL report the same way as G2a when `openspec` is unavailable:
  a clear SKIP or FAIL, never silence.
- The G6 message SHALL name the real control, `skip_specs: true`.
- The merged branches 66, 85, 92, and 94 SHALL be deleted from the remote.

## Success criteria

- The `cli` and `gates` checks run the e2e against a real OpenSpec install and
  pass.
- Running `scripts/e2e-test.sh` without openspec on PATH exits non-zero with a
  clear message.
- `openspec validate --all --strict` and `openspec validate --archived` pass on
  the supported OpenSpec version.
- The README version claim matches the version the workflows install, and the
  mutilated sentence reads correctly.
- `git ls-remote --heads` lists no merged branches.

## Out of scope

- No new CLI features or behavior changes.
- No change to the workflow commands, agents, or the OpenSpec schema content.

