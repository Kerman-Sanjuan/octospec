## Context

The README documents eight gates, but only G2 to G6 are implemented, and only `release-pipeline` mentions them. There is no living spec for the gate set, so a change to a gate has no contract. The release acceptance workflow verifies the published binary but does not re-run the gates.

## Goals / Non-Goals

**Goals:**
- G1 exists and runs on the same terms as the other stack-independent gates.
- A `gates` spec is the contract for every gate and its strength.
- The release checks the gates before it trusts the binary.
- One hermetic test proves seed, install, and the gates work together.

**Non-Goals:**
- Change the strength of an existing gate.
- Add a gate that needs GitHub to run.
- Test agent behaviour, which is not deterministic code.

## Decisions

- **G1 reads the linked issue through `gh`.** The change's `.openspec.yaml` names the issue; the gate fetches the body and checks for the four required sections. *Alternative rejected:* parse the issue from the change's `issue.md`, which is a snapshot and can be stale.
- **G1 skips when `gh` is unavailable.** CI provides `gh`; a local run without it skips rather than fails. *Alternative rejected:* fail on a missing `gh`, which breaks local runs.
- **The release acceptance keeps its OS matrix.** A `gates` job runs once; the `acceptance` job keeps linux and macOS, so the darwin binary stays covered. *Alternative rejected:* fold the gates into the matrix, which runs them once per OS for no gain.
- **The end-to-end test is hermetic.** It builds the CLI, seeds and installs into a temporary repository, and runs the gates. It does not call GitHub. *Alternative rejected:* a test against a real repository, which needs a token and a network.

## Risks / Trade-offs

- [G1 depends on the issue body staying in sync with the change] -> The snapshot `issue.md` is still written by `/spec`; G1 is a merge-time check, not a runtime one.
- [The end-to-end test adds CI time] -> It builds the CLI once and runs in seconds.

## Migration Plan

Add G1, the `gates` spec, the end-to-end test, and the release job. Remove `/explore`.

## Open Questions

- None.
