## Why

The gate system is described in the README but has no living spec, so a change to a gate has no contract to hold it. G1 is documented but not implemented, and the release acceptance workflow verifies the binary without re-running the gates. The result is a quality system that is hard to trust for v1.0.

## What Changes

- Implement G1 in `scripts/check-gates.sh`: the linked issue carries its required sections.
- Run the gates in `release-acceptance.yml` before the binary is verified.
- Add a living spec, `gates`, that documents G1 through G8 and their strength.
- Add `scripts/e2e-test.sh`, a hermetic test of `octospec seed`, `octospec install`, and the gates.
- Remove `/explore`; its definition duplicates the OpenSpec explore skill.

## Capabilities

### New Capabilities

- `gates`: the gate set, its checks, and its strength.

### Modified Capabilities

- `release-pipeline`: the release acceptance test also runs the gates.

## Impact

- `scripts/check-gates.sh`, `scripts/check-gates-test.sh`, `scripts/e2e-test.sh`
- `.github/workflows/release-acceptance.yml`
- `cli/internal/payload/commands/explore.md` (removed) and `README.md`
- `cli/internal/payload/embed_test.go`

## Changelog

### Added
- G1 in `scripts/check-gates.sh`: the linked issue must carry its required sections.
- A living `gates` spec documenting G1 through G8.
- `scripts/e2e-test.sh`, a hermetic end-to-end test of seed, install, and the gates.

### Changed
- The release acceptance workflow runs the gates before it verifies the binary.

### Removed
- `/explore`, which duplicated the OpenSpec explore skill.
