---
description: File a new bug issue, after grilling what actually happens.
argument-hint: "[summary]"
---

Runs the `idea` stage in the current session. The `idea` agent defines its scope, tools, and model.

Create a bug issue. The issue is the intake artifact for a later fix.

1. If `$ARGUMENTS` is empty, ask for a summary.
2. Reproduce it before asking. Read the code path, run the command, read the
   log; a fact the tree can answer is not a question for the user.
3. Grill the decision frontier, in rounds: numbered questions, each with your
   recommended answer, worded so that a plain "yes" accepts it. A question that
   depends on one still open waits for the next round. For a bug the frontier
   covers what happens, what should happen, the smallest reproduction, the
   conditions it needs, who is affected and how badly, and the evidence that
   shows it (log line, failing test, trace).
4. File nothing while a decision on the frontier is unanswered. Name the open
   decision and stop. With nobody to answer, stop with `reason: needs_human`.
5. Create the issue from the sections as the body, with the labels. Write the
   body to a file to keep the formatting:
   `gh issue create --title "<summary>" --body-file <file> --label type:bug --label status:backlog`
   The body carries `## Grill` with one line per closed decision, and
   `## Open questions` with `None`.
   The YAML form in `.github/ISSUE_TEMPLATE/bug.yml` is the human web-form
   contract. `gh` cannot prefill a form from the CLI, so pass the body.
6. Print the issue URL and number. Stop.
