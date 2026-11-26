# Design — G9

## Context

`spec.md` says what must be true. `tests` say what is checked. Nothing connects them, so the review of
a change can only ask "does this look tested", which is the question the whole harness exists to
replace.

## Goals / Non-Goals

**Goals**

- One file per change, one line per requirement, checked by a script.
- A warning that names the requirement, not a count.
- Zero new metadata inside the specs themselves.

**Non-Goals**

- Proving that the test *actually* covers the requirement. A human or a review says that; the gate says
  that something exists and is named.
- Language-specific parsing. The symbol is a string that must be in the tree; whether it is a Go test,
  a shell script or a Playwright case is the repository's business.

## Decisions

### The mapping lives in the change, not in the spec

Adding `tests:` to a requirement would change the OpenSpec schema and put coverage in the artefact that
is archived. The mapping is a fact about this change's work, so it lives with it and disappears with
it — the specs stay readable prose.

### The symbol is a string, verified by grep

`TestThePublishedSchemaMatchesTheCode` is a symbol. So is `scripts/check-gates.sh`. G9 checks only that
the string appears somewhere outside `openspec/`, which is enough to catch the real failure: a mapping
written from memory, pointing at a test renamed three commits ago.

### Advisory first

The precedent in this repository is G7 and G8: they shipped as warnings, found their false positives,
and then became rules people trusted. G9 will do the same, and flipping it to hard is one word.

## Risks

- A requirement that is truly untestable (documentation) cannot pass. The honest answer is a symbol
  pointing at the document check — or splitting the requirement. We will see which case dominates.
