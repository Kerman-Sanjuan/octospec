# Design — the intake grill

## Context

`/idea` asks for five form fields, one at a time, and writes the issue when the fields are filled. The
fields are a template, not a test of decidability: "improve the intake" fills `## Requirements` and
passes G1. Nothing carries the conversation to `/spec`, so every gap is reconstructed by guesswork one
stage later, where it costs ten times more.

`grilling` (Matt Pocock) is the practice worth taking: a decision tree worked in rounds, the frontier of
decisions whose prerequisites are settled, an answer recommendation per question. 540 chains it with
`to-spec` and `to-tickets`; we keep our own stages and take only the interrogation.

## Goals / Non-Goals

**Goals**

- Stop the intake when nothing is silently assumed, not when the form is full.
- Ask the user only for decisions; every fact is looked up.
- Leave the answers in the issue, so `/spec` reads instead of re-asking.
- Make G1 able to tell a decidable issue from a filled template.

**Non-Goals**

- The naive reader at the end of `/spec` (separate change; it touches the review arithmetic).
- The `to-tickets` breakdown: `/spec` already produces tasks, and 540's acceptance-criteria table belongs
  to a separate decision.
- Proving a requirement is testable. G1 checks form; judgement stays with a reader.

## Decisions

### The frontier, not one question per turn

The current payload says "one question at a time", which is a well-meant simplification of a worse
practice: it serialises decisions that do not depend on each other into twenty turns, and a tired answer
is worse than no answer. The source asks the whole frontier in one round, each question with a
recommendation. We adopt the round: it is faster, it shows the shape of the problem, and "yes" accepts
the recommendation.

The cost is a longer message. Accepted: a round is read once, a twenty-turn interrogation is read twenty
times.

### Facts are the agent's job, decisions are the user's

The `idea` agent already has `read`, `search` and `shell`. Every question that can be answered by looking
at the tree or the tracker is a question wasted, and it teaches the user that the agent is a form. The
payload states the rule; the agent's tool surface already allows it.

### `## Grill` is a section of the issue, not a comment

A comment is not read by `/spec` (it reads the issue body), and the issue body is what gets snapshotted
into `issue.md` and read by a fresh session. One line per closed decision keeps it stable under G1 and
readable by a human in ten seconds.

### G1 grows a second half, and says what it cannot do

Three checks are deterministic: the `## Grill` heading, `SHALL`/`MUST` in every requirement line, and an
empty `## Open questions`. They are cheap and they catch the real failure, an issue that says what to do
without saying what must be true.

A fourth — that a requirement is *testable* — is not. A grep cannot tell "the gate SHALL fail when the
body is empty" from "the system SHALL be faster". We do not fake it with a word list (fragile, and it
punishes good prose); the spec review in the next change is where that judgement lives.

### The hard part: existing issues

G1 is hard for new requirements. An issue filed before this change has no `## Grill` and fails until it
is edited. That is the standard talking, not a bug: the repository has 16 open issues and one user, and
an issue that cannot answer "what is out of scope" is not ready for a change.

## Risks / Trade-offs

- [A round of questions can overwhelm] → the recommendation per question means the user can answer
  "yes, yes, no" without writing prose. Rounds stay small because only the settled frontier is asked.
- [SHALL/MUST becomes ritual] → it is one grep, not the standard. The seven groups and the reader do the
  real work; the gate only stops the empty issue.
- [G1 fails every old issue] → intended, and named in the failure message so the fix is obvious.

## Migration Plan

No data migration. Existing open issues get a `## Grill` section the next time a change touches them; the
gate's message names the missing section and the offending requirement line.

## Open Questions

None.
