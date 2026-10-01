## Context

The repository is public. It has no `dependabot.yml` and no CodeQL workflow, so dependency updates and code scanning are off.

## Goals / Non-Goals

**Goals:**
- Dependabot opens updates for Go modules and Actions.
- CodeQL scans Go on pull requests and on a schedule.

**Non-Goals:**
- The seed default for other repositories (#39).

## Decisions

- **Dependabot for `gomod` and `github-actions`, weekly.** Weekly is enough for a small project and keeps the noise low. *Alternative rejected:* daily, which is noisy for a small dependency set.
- **A committed CodeQL workflow, not the default setup.** A committed workflow is reproducible and reviewable, and matches how the rest of CI is managed. *Alternative rejected:* the repository's default code scanning setup, which lives only in settings.

## Risks / Trade-offs

- [CodeQL needs `security-events: write`] -> The workflow grants it, and the repository is not a fork, so uploads are allowed.
- [Dependabot pull requests need review like any other] -> They run the same `cli` and `gates` checks.

## Migration Plan

None.

## Open Questions

- None.
