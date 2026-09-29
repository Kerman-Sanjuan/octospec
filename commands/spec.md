---
description: Turn a GitHub issue into an OpenSpec change, publish it to the issue, and validate.
argument-hint: "<issue-number|url>"
---

Load the `openspec-propose` skill and follow it. The artifact semantics come
from the `octospec` schema, not from this file.

1. Resolve the issue from `$ARGUMENTS`:
   `gh issue view <n> --json number,title,body,updatedAt,url`
   If it does not exist, stop and tell the user to run `/idea`.
2. Derive `<slug>` from the title; the change name is `gh-<issue>-<slug>`
   (OpenSpec requires a leading letter, so the number cannot come first).
3. `openspec new change "gh-<issue>-<slug>" --schema octospec`
4. Let the skill create every artifact required by
   `openspec status --change "<change>" --json` (`issue`, `proposal`,
   `specs`, `design`, `tasks`).
5. Publish each artifact to the issue as a comment, each headed with the
   base commit SHA (`git rev-parse HEAD`); `/spec` does not commit, so the
   artifacts are uncommitted until `/apply`:
   - proposal and design: one comment each.
   - specs: one comment per capability.
   - tasks: one comment containing the checklist.
   Use `gh issue comment <n> --body-file <file>`.
6. Run `openspec validate "<change>" --strict`. If it fails, fix and repeat.
7. Add label `status:spec-ready`; remove `status:backlog`.

Do not commit. Report the change path, the issue URL, and the validation
result.
