## Why

An audit found that the `e2e` CI step prints `SKIP: openspec is not on PATH`
and exits 0 (a green check that proves nothing), the README claims OpenSpec
"1.3.1 or newer" while the gates fail on 1.14.0, the README prose is mutilated,
three archived changes have unchecked tasks, the gate script is inconsistent
when openspec is missing, and four merged branches were never deleted.

## What Changes

- Install the OpenSpec CLI in the `cli` CI job so the `e2e` step actually runs,
  and make `e2e-test.sh` fail (not skip) when `openspec` is unavailable.
- Correct the OpenSpec version range in the README to the range the gates
  support, and keep the workflows inside it.
- Repair the mutilated prose in the README skills paragraph.
- Check the open tasks in the archived changes gh-30, gh-31, and gh-48.
- Make G2b and G7 report a SKIP like G2a when `openspec` is missing.
- Fix the G6 message to name `skip_specs: true`.
- Delete the merged branches 66, 85, 92, and 94.

## Capabilities

### New Capabilities

### Modified Capabilities
- `release-pipeline`: the `cli` job installs OpenSpec so the e2e runs, and the
  e2e fails when OpenSpec is missing.
- `readme`: the OpenSpec version range is accurate and the prose is complete.
- `gates`: G2b and G7 report a SKIP consistently, the G6 message names the real
  control, and archived changes carry no unchecked tasks.

## Impact

- `scripts/e2e-test.sh`, `.github/workflows/cli.yml`, `README.md`,
  `scripts/check-gates.sh`, and the three archived changes.
- Remote branches 66, 85, 92, and 94.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/98

## Changelog

### Fixed
- The `e2e` CI step now installs OpenSpec and fails instead of skipping when it
  is missing.
- The README states the OpenSpec version the gates support and repairs a
  truncated sentence.
- G2b and G7 report a SKIP when OpenSpec is missing, and the G6 message names
  `skip_specs: true`.
