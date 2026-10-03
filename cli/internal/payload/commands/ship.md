---
description: Run CI, open a pull request linked to the issue, and iterate failures back to apply.
argument-hint: "[change]"
---

Runs the `ship` stage in the current session. The `ship` agent defines its scope, tools, and model.

Open the PR for the current branch once the gates are green.

1. Read `.openspec.yaml` for `github.issue`.
2. Push the branch: `git push -u origin HEAD`.
3. Create or update the pull request. The body MUST contain `Closes #<issue>`,
   then:
   - a short summary,
   - the change name and artifact paths,
   - the `openspec validate` result.
   `gh pr create --title "<type>(#<issue>): <summary>" --body-file <file>`
4. Run or observe CI (`gh pr checks --watch`). If any gate fails, hand the
   specific failing check back to `/apply` to fix and iterate - do not report a
   bare failure, and do not continue until it is green.
5. Once every gate is green, add label `status:in-review`; leave
   `status:in-progress` (the `/archive` step clears it).
6. Print the PR URL.
