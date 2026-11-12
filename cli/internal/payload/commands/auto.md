---
description: Run the full idea-to-archive loop unattended, retrying failed stages.
argument-hint: "[issue-number]"
---

Runs the cycle autonomously in the current session: select an active issue, drive it through the remaining stage commands (`/spec`, `/apply`, `/ship`, `/archive`), retry a failing stage, and leave the same GitHub trail a manual run would. No human confirmation between stages. The stage agents keep their scope, tools, and models.

1. Select the issue to work. If `$ARGUMENTS` gives a number, use it.
   Otherwise list the open issues and take the one with the lowest number
   that carries an active status label (`status:spec-ready`,
   `status:in-progress`, or `status:in-review`):
   `gh issue list --state open --json number,title,labels`
   Ignore issues in `status:backlog`: they wait for a human decision.
   If nothing is active, say so and stop.
2. Resume from the issue's state, do not redo finished work:
   - `status:spec-ready` -> start at `/apply` (the plan is drafted).
   - `status:in-progress` without `status:in-review` -> start at `/apply`
     (mid-apply; re-check the change state before working).
   - `status:in-progress` with `status:in-review` -> start at `/ship`
     (there is a pull request; observe its checks).
3. Run the remaining stages in order, one issue at a time, following each
   stage command's own contract unchanged - including its single-comment
   and label rules. Do not restate or weaken them.
4. Retry rule: when a stage fails or a gate does not pass, re-run THAT
   stage, not the next one. A stage gets at most 3 attempts (the retry
   bound). The CLI subcommand `octospec auto --max-retries <n>` raises or
   lowers it for terminal runs.
5. Give up visibly: past the retry bound, post a comment on the issue
   headed `## Autonomous run failed`, naming the stage, the attempt count,
   and the last error. Then stop working that issue, select the next
   active one, and continue. Never fail silently.
6. Bugs found while working: when behaviour contradicts `openspec/specs/`,
   file a bug issue through the `/bug` intake (`type:bug`,
   `status:backlog`), quoting the spec it contradicts. Never invent
   requirements. Then continue the current issue.
7. Record everything on GitHub - comments, commits, pull requests, label
   changes - exactly as the manual stage commands do. Keep no run state
   outside GitHub; the issue's labels ARE the run state.
8. When no active issue remains, print a run summary: issues completed,
   issues abandoned with their failure comments, and stop.
