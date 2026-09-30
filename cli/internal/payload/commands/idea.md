---
description: File a new feature issue in the backlog and interview for its content.
argument-hint: "[idea]"
---

Runs the `idea` agent. The agent definition sets the scope and the tool surface.

Create a feature issue in the GitHub backlog. The issue is the human source
of truth; there is no OpenSpec change yet.

1. If `$ARGUMENTS` is empty, ask what the user wants to build.
2. Interview for the feature form fields, one question at a time: user story,
   context, requirements (SHALL/MUST), success criteria, out of scope.
3. Create the issue from the form:
   `gh issue create --template feature.yml --title "<title>" --body "<fields>"`
   If the form is unavailable, pass the same sections as the body.
4. Add labels `type:feature` and `status:backlog`.
5. Print the issue URL and number. Stop.

Do not create an OpenSpec change. `/spec` does that later.
