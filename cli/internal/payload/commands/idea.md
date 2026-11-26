---
description: File a new feature issue in the backlog, after grilling its decisions.
argument-hint: "[idea]"
---

Runs the `idea` stage in the current session. The `idea` agent defines its scope, tools, and model.

Create a feature issue in the GitHub backlog. The issue is the human source
of truth; there is no OpenSpec change yet.

1. If `$ARGUMENTS` is empty, ask what the user wants to build.
2. Look the facts up before asking. Read the repository, the existing issues
   and every file the idea names; a question the tree can answer is not asked.
3. Grill the decision frontier, in rounds: the decisions whose prerequisites are
   already settled, numbered, each with your recommended answer, worded so that
   a plain "yes" accepts it. A question that depends on one still open waits for
   the next round. The frontier covers who acts, the observable outcome, what
   done looks like, the failures that matter, the boundary and why it is out of
   scope, the constraints already known, and the evidence that will prove it.
4. File nothing while a decision on the frontier is unanswered. Name the open
   decision and stop. With nobody to answer, stop with `reason: needs_human`.
5. Create the issue from the sections as the body, with the labels. Write the
   body to a file to keep the formatting:
   `gh issue create --title "<title>" --body-file <file> --label type:feature --label status:backlog`
   The body carries `## Grill` with one line per closed decision, and
   `## Open questions` with `None`.
   The YAML form in `.github/ISSUE_TEMPLATE/feature.yml` is the human web-form
   contract. `gh` cannot prefill a form from the CLI, so pass the body.
6. Print the issue URL and number. Stop.

Do not create an OpenSpec change. `/spec` does that later.
