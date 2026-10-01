---
description: File a new feature issue in the backlog and interview for its content.
argument-hint: "[idea]"
---

Runs the `idea` agent. The agent definition sets the scope and the tool surface.

Create a feature issue in the GitHub backlog. The issue is the human source
of truth; there is no OpenSpec change yet.

1. If `$ARGUMENTS` is empty, ask what the user wants to build.
2. Interview for the feature fields, one question at a time: user story,
   context, requirements (SHALL/MUST), success criteria, out of scope.
3. Create the issue from the sections as the body, with the labels. Write the
   body to a file to keep the formatting:
   `gh issue create --title "<title>" --body-file <file> --label type:feature --label status:backlog`
   The YAML form in `.github/ISSUE_TEMPLATE/feature.yml` is the human web-form
   contract. `gh` cannot prefill a form from the CLI, so pass the body.
4. Print the issue URL and number. Stop.

Do not create an OpenSpec change. `/spec` does that later.
