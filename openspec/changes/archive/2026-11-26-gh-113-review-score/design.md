# Design — review score

## Context

Verification today stops at "the build is green". The decision that actually consumes a human is
"does this change need me", and nothing produces it. The pattern that works elsewhere is two
layers: the model finds things, and something deterministic turns those findings into a number that
can be argued with, audited and compared between runs.

## Goals / Non-Goals

**Goals**

- One number per pull request, reproducible from the findings it cites.
- A risk level that answers "who has to look" using facts, not impressions.
- The numbers belong to the repository, so tuning never needs a release.

**Non-Goals**

- Replacing G5 or the deterministic gates. A red build is a precondition, not a score.
- Judging whether the feature is worth building. That is discovery, not review.
- Learning across repositories.

## Decisions

**Decision: the arithmetic is Go inside octospec, not a script in the adopting repository.**
The finding type→weight mapping and the risk table are the contract; if each repository
reimplemented them, two repositories would disagree about what 7.5 means and the number would stop
being comparable. Rejected: a `score.sh` per repository (drifts, untestable) and a model call to
produce the number (two runs disagree, which is the failure 540 hit and fixed).

**Decision: findings arrive as JSON on stdin or a file, with a closed vocabulary for axis and type.**
The agent is good at finding; it is not trusted with the arithmetic. An unknown axis fails loudly
rather than being silently scored as "other".

**Decision: risk is a small decision table, not a score.**
Three inputs — sensitive path hit, diff size, author familiarity — and four outcomes. A risk
*score* invites tuning the weights until everything is medium; a table can be argued about and
tested.

**Decision: report-only by default, `--enforce` to block.**
The issue's own caveat: 540's thresholds come from application repositories with a team calibrating
weekly. A single maintainer must first see twenty PRs. Blocking on day one with a guessed number is
how a gate gets disabled.

**Decision: exit codes 0 and 5.**
`5` means "blocked by the review", distinct from `1` (the command failed) and from the lock's `3`
and `4`. A CI job can block on `5` without parsing English.

## Risks / Trade-offs

- [The model over-reports, so everything scores low] → weights are policy: lower them in
  `.octospec/review.json`, which is exactly the calibration the first twenty PRs are for.
- [The model under-reports and a bad change passes] → G9 (requirement ↔ test) and the deterministic
  gates stay in front; the score is the last layer, never the only one.
- [A finding cites a file:line that does not exist] → the score still lands, but the comment shows
  the citation so a human can catch it in the PRs that matter.

## Migration Plan

Additive. No `.octospec/review.json` means defaults, and nothing blocks until a repository adds
`--enforce` to its workflow.

## Open Questions

- Should the five axes carry their own floors (a correctness finding can never score above 6, say)?
  Deferred until the first twenty PRs show whether one axis dominates.
