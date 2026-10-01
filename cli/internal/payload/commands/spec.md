---
description: Turn a GitHub issue into an OpenSpec change, commit it on the change branch, publish the plan, and validate.
argument-hint: "<issue-number|url>"
---

Runs the `spec` stage in the current session. The `spec` agent defines its scope, tools, and model.

Load the `openspec-propose` skill and follow it. The artifact semantics come
from the `octospec` schema, not from this file.

1. Resolve the issue from `$ARGUMENTS`:
   `gh issue view <n> --json number,title,body,updatedAt,url,labels`
   If it does not exist, stop and tell the user to run `/idea`.
2. If the issue is ambiguous - a requirement, a user story, or a bug is not
   clear enough to plan - ask clarifying questions one at a time and route the
   answers to the issue body. Do not invent requirements, and do not finalize
   the plan until the issue body is clear.
3. Derive `<slug>` from the title; the change name is `gh-<issue>-<slug>`
   (OpenSpec requires a leading letter, so the number cannot come first).
   Derive the branch prefix from the issue's `type:*` label
   (`type:feature` -> `feat`, `type:bug` -> `fix`); the branch is
   `feat|fix/<issue>-<slug>`.
4. If `openspec/changes/gh-<issue>-<slug>/` already exists, reuse it and its
   branch (idempotent re-run) - do not scaffold a second change. Otherwise:
   `openspec new change "gh-<issue>-<slug>" --schema octospec`
5. Create or switch to the change branch:
   `git checkout -b feat|fix/<issue>-<slug>` (or `git checkout` it if it
   already exists).
6. Let the skill create every artifact required by
   `openspec status --change "<change>" --json` (`issue`, `proposal`,
   `specs`, `design`, `tasks`).
7. Commit the artifacts on the branch and push:
   `git add openspec/changes/<change>`
   `git commit -m "docs(#<issue>): add spec artifacts"`
   `git push -u origin HEAD`
8. Publish a single **plan comment** to the issue, headed with the commit SHA
   the artifacts were committed at (`git rev-parse HEAD`), with fixed
   headings:
   - `## Plan`
   - `### Proposal`
   - `### Capabilities` (one `#### <capability>` per capability)
   - `### Design`
   - `### Tasks` (the checklist)
   Record its id in `.openspec.yaml` as `github.comments.plan`. On re-run,
   update that comment in place instead of adding a new one:
   `gh api --method PATCH /repos/{owner}/{repo}/issues/comments/<id> -F body=@<file>`
9. Run `openspec validate "<change>" --strict`. If it fails, fix, commit,
   and push again, then refresh the plan comment.
10. Add label `status:spec-ready`; remove `status:backlog`. Do **not** add an
    approval label: `status:spec-ready` means the approach is drafted. The
    human approves the approach by running `/apply`.

Report the change path, the branch, the issue URL, and the validation result.
