---
mode: agent
description: Open a pull request linked to the issue.
---


Open the PR for the current branch.

1. Read `.openspec.yaml` for `github.issue`.
2. Push the branch: `git push -u origin HEAD`.
3. Create the PR. The body MUST contain `Closes #<issue>`, then:
   - a short summary,
   - the change name and artifact paths,
   - the `openspec validate` result.
   `gh pr create --title "<type>(#<issue>): <summary>" --body-file <file>`
4. Add label `status:in-review`; leave `status:in-progress` (the `/archive`
   step clears it).
5. Print the PR URL.
