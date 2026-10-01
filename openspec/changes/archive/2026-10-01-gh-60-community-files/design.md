## Context

The repository is public. `.github/ISSUE_TEMPLATE/` has `feature.yml` and `bug.yml`, but there is no security policy, no code of conduct, no pull request template, and no issue chooser config. `CONTRIBUTING.md` carries a short inline code of conduct.

## Goals / Non-Goals

**Goals:**
- The standard community files exist and GitHub surfaces them.
- The templates keep contributions consistent.

**Non-Goals:**
- A contributor license agreement.
- A chat server.

## Decisions

- **Adopt the Contributor Covenant.** It is the common standard, and `CONTRIBUTING.md` links to it instead of restating it. *Alternative rejected:* keep the inline paragraph only, which GitHub does not surface.
- **Report vulnerabilities through GitHub's private reporting.** It avoids publishing a personal email and keeps the report private. *Alternative rejected:* a personal email address in a public file.
- **One pull request template.** The workflow already asks for `Closes #<issue>`, a summary, and the validation result, so the template mirrors that. *Alternative rejected:* several templates, which add choice without benefit.

## Risks / Trade-offs

- [Private vulnerability reporting must be enabled on the repository] -> Enable it in the settings; the file links to it, and the task records the step.
- [The code of conduct needs a real enforcement contact] -> Use the maintainer's public commit email.

## Migration Plan

None.

## Open Questions

- None.
