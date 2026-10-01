## Context

`CONTRIBUTING.md` says to use GitHub Discussions for questions. Discussions are disabled on the repository, so a contributor finds nothing. The issue chooser config links the security policy and the contributing guide, but not a question channel.

## Goals / Non-Goals

**Goals:**
- A real question channel exists and is linked from the places that name it.

**Non-Goals:**
- A chat server.
- Moderating categories.

## Decisions

- **Enable GitHub Discussions.** It is built into the repository, needs no external service, and matches what `CONTRIBUTING.md` already says. *Alternative rejected:* reword `CONTRIBUTING.md` to remove the reference, which drops the channel instead of providing it.
- **Link Discussions from the issue chooser config.** The chooser is the first place a visitor meets, so it points at the question channel. *Alternative rejected:* only the CONTRIBUTING mention, which is easy to miss.

## Risks / Trade-offs

- [Discussions adds a surface to moderate] -> Default categories are enough; the code of conduct covers behaviour.
- [The chooser link breaks if Discussions are disabled again] -> The same change enables them.

## Migration Plan

None.

## Open Questions

- None.
