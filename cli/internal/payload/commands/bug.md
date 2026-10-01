---
description: File a new bug issue.
argument-hint: "[summary]"
---

Runs the `idea` agent. The agent definition sets the scope and the tool surface.

Create a bug issue. The issue is the intake artifact for a later fix.

1. If `$ARGUMENTS` is empty, ask for a summary.
2. Interview for the bug fields, one question at a time: summary, steps to
   reproduce, expected, actual, impact.
3. Create the issue from the sections as the body, with the labels. Write the
   body to a file to keep the formatting:
   `gh issue create --title "<summary>" --body-file <file> --label type:bug --label status:backlog`
   The YAML form in `.github/ISSUE_TEMPLATE/bug.yml` is the human web-form
   contract. `gh` cannot prefill a form from the CLI, so pass the body.
4. Print the issue URL and number. Stop.
