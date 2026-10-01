## Context

The README has a "Why" section, the loop, the architecture, and then the install. A visitor must scroll to the install, and the install prerequisites list the OpenSpec CLI and `gh` but not Node.

## Goals / Non-Goals

**Goals:**
- A newcomer can try octospec from the top of the README.
- The prerequisites are complete.

**Non-Goals:**
- A demo recording (#67).
- Reworking the rest of the README.

## Decisions

- **A quickstart near the top, right after the intro.** It shows the three steps and the loop in one block. *Alternative rejected:* rely on the deeper install section, which a visitor must find.
- **Name Node in the prerequisites.** The OpenSpec CLI installs through npm, so Node is a real requirement. *Alternative rejected:* mention it only in the install section, where it is easy to miss.

## Risks / Trade-offs

- [The quickstart and the install section drift] -> The quickstart links to the install section for the details.

## Migration Plan

None.

## Open Questions

- None.
