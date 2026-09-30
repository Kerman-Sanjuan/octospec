---
description: File a new bug issue.
argument-hint: "[summary]"
---

Runs the `idea` agent. The agent definition sets the scope and the tool surface.

Create a bug issue. The issue is the intake artifact for a later fix.

1. If `$ARGUMENTS` is empty, ask for a summary.
2. Interview for the bug form fields: summary, steps to reproduce, expected,
   actual, impact.
3. `gh issue create --template bug.yml --title "<summary>" --body "<fields>"`
4. Add labels `type:bug` and `status:backlog`.
5. Print the issue URL and number. Stop.
