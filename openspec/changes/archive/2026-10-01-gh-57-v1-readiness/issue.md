# Issue #57: v1.0 readiness: gate spec, G1, release gates, and an end-to-end test

https://github.com/Kerman-Sanjuan/octospec/issues/57

## User story

As a maintainer, I want the gate system documented as a living spec and covered by an end-to-end test, so that v1.0 ships with enforced, verifiable quality controls.

## Context / problem

The gates are described in the README and enforced by `scripts/check-gates.sh`, but:

- G1 (issue sections) is described but not implemented.
- the release acceptance workflow does not re-run the gates.
- there is no living spec for the gate set, so changes to gates have no contract.
- there is no end-to-end test that exercises install, seed, and the gates together.

## Requirements

- The gate script shall implement G1: the linked issue carries the required sections.
- The release acceptance workflow shall run the gates before it verifies the binary.
- A living spec named `gates` shall document G1 through G8 and their strength.
- An end-to-end test shall exercise `octospec seed`, `octospec install`, and the gate script in a temporary repository.

## Success criteria

- `scripts/check-gates.sh` reports G1.
- The release acceptance workflow runs `scripts/check-gates.sh`.
- `openspec/specs/gates/spec.md` documents every gate.
- `scripts/e2e-test.sh` passes locally and in CI.

## Out of scope

- Changing the gate names or the strength of an existing gate.
- A test that calls GitHub; the end-to-end test stays hermetic.
