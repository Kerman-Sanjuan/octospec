# Issue #134: Grill the issue, then test the spec with a naive reader

https://github.com/Kerman-Sanjuan/octospec/issues/134

## User story

As the person who has the idea, I want the definition phase to interrogate me and then be tested by a
reader with no context, so that the specification is buildable by somebody who was not in the
conversation — a machine included.

## Context / problem

Three inputs, and none of them is ours yet:

1. **`grilling`** (Matt Pocock, MIT), the practice 540 adapts in their software-factory series: before a
   spec exists, the agent interrogates the person until the outcome is testable, instead of writing a
   beautiful specification out of a vague sentence.
2. **google/elephant-goldfish-model** (Apache-2.0), whose reviewers are stateless and read-only, and whose
   `comprehension` persona contains the idea worth most: **paraphrase first**.
3. **Our reality.** `/idea` interviews lightly, `/spec` drafts from a thin issue, and `/apply` fills the
   gaps by guessing.

This change takes the first half of the issue: the grill. The naive reader is a separate change, because
it touches `/spec` and the review arithmetic.

## Requirements

- `/idea` and `/bug` SHALL work the request as a decision tree and ask, in rounds, only the decisions
  whose prerequisites are settled, and SHALL NOT file the issue while any decision on the frontier is
  unanswered.
- Each question in a round SHALL be numbered and SHALL carry the agent's recommended answer.
- The agent SHALL look up every fact it can read (repository, files, existing issues) instead of asking
  the user for it.
- The frontier SHALL cover at least: who acts, the observable outcome, what done looks like, the failures
  that matter, the boundary (what is out of scope and why), the constraints already known, and the
  evidence that will prove it.
- The answers SHALL land in the issue body under `## Grill`, one closed decision per line, so that
  `/spec` reads them and never re-asks a question it could have read.
- In an unattended run the grill SHALL NOT invent answers: it SHALL stop with `reason: needs_human`.
- G1 SHALL fail when the linked issue has no `## Grill` section, when a requirement under
  `## Requirements` carries neither `SHALL` nor `MUST`, when `## Success criteria` is empty, or when
  `## Open questions` still has content.

## Success criteria

- On a request whose outcome is not testable, the agent asks the frontier in rounds and files nothing
  until the last decision is answered.
- A filed issue carries `## Grill` with the closed decisions and no open question.
- G1 fails an issue whose requirements say what to do without SHALL/MUST, and passes one that states
  them.
- On five consecutive real changes, `/spec` asks nothing that `## Grill` already answers.

## Grill

- Scope: the grill and its gate in this change; the naive reader is a separate change.
- Round shape: the whole settled frontier per round, each question numbered and carrying a recommended
  answer, instead of one question per turn.
- Facts versus decisions: the agent reads the repository, the files and the tracker; the user is asked
  only for decisions.
- Stop rule: the grill ends when the frontier is empty, not when the form fields are filled. Unattended
  stops with `reason: needs_human`.
- Record: the closed decisions live in the issue body under `## Grill`, one line each, so `/spec` reads
  them and does not re-ask.
- Gate: G1 is hard and checks the shape a script can see — `## Grill` present, SHALL/MUST in every
  requirement line, success criteria non-empty, `## Open questions` empty.
- Old issues: an issue filed before this change fails G1 until it gains the section. Accepted.
- Attribution: `grilling` is Matt Pocock's skill (MIT); 540 adapts it in their factory.


## Out of scope

- The naive reader at the end of `/spec` (separate change, touches the review scoring).
- Porting EGM's five-phase pipeline and nine personas.
- Proving that a requirement is testable. A grep checks the form; only a reader judges the substance.

## Open questions

None
