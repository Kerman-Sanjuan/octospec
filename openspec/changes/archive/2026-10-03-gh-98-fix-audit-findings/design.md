## Context

The `cli` job in `.github/workflows/cli.yml` mounts `setup-node` but never runs
`npm install -g @fission-ai/openspec`, so `scripts/e2e-test.sh` hits its early
`command -v openspec` guard and prints `SKIP: openspec is not on PATH`, exiting
0. The README claims OpenSpec "1.3.1 or newer" (line 24) and ">= 1.3.1" (line
150), but the seed schema writes a `github:` key into `.openspec.yaml`, which
OpenSpec 1.14.0 rejects under `--strict`. The workflows pin `1.3.1`, hiding the
mismatch. `scripts/check-gates.sh` guards G2a with an `if/else` that reports
FAIL, but G2b and G7 wrap their whole body in `if command -v openspec` with no
else, so they pass silently. The G6 message references a nonexistent flag.

## Goals / Non-Goals

**Goals:**
- Make the e2e step run in CI and fail when OpenSpec is missing.
- Make the README version claim match what the gates support.
- Make the gate script consistent when OpenSpec is missing.

**Non-Goals:**
- No change to the OpenSpec schema or the `.openspec.yaml` `github:` key, which
  the workflow depends on.
- No new CLI features.

## Decisions

### Install OpenSpec in the cli job at the tested version
Add a step that runs `npm install -g @fission-ai/openspec@1.3.1` before the e2e,
matching the `openspec` workflow and `release-acceptance`. `setup-node` is
already present. This makes the green check real.

### Make the e2e guard fail instead of skip
`scripts/e2e-test.sh` currently prints a SKIP and `exit 0` when openspec is
absent. Change it to exit non-zero with a clear message, so a missing dependency
is a failure, not a quiet pass. The alternative, keeping the skip, is exactly
the decorative-green problem the audit found.

### State the tested version, not a range
The schema needs the `github:` key, and OpenSpec 1.14.0 warns on it under
`--strict`. Rather than claim "or newer", the README names the tested version
(1.3.1). This is truthful and matches the pins. The alternative, dropping the
`github:` key, would break the workflow's traceability and every existing
change.

### Consistent SKIP in the gate script
Give G2b and G7 an `else` branch that prints a SKIP like G2a, so a missing
`openspec` never reads as a pass. G2a keeps its FAIL, because structural
validity is a hard requirement; G2b and G7 depend on `openspec status` and a
SKIP is the honest report.

### Check the archived tasks and delete the merged branches
The three archived changes (gh-30, gh-31, gh-48) carry unchecked boxes left over
from their execution. Check them so `openspec validate --archived` passes.
Delete the merged remote branches 66, 85, 92, and 94, which `/archive` promises
to do.

## Risks / Trade-offs

- [Pinning openspec 1.3.1 while the README says "tested" could still rot] -> The
  README now names one version, and a future bump updates the README and the
  pins together; the mismatch is visible, not silent.
- [Making e2e fail without openspec could redden a contributor's local run] ->
  The message names the install command, and CI always installs it, so only a
  misconfigured local run fails, which is the point.

## Migration Plan

No data migration. Rollback is a revert. Deleting the merged branches is
irreversible but they are merged into `main`.

## Open Questions

- None. The audit pins each finding to a file and a line.
