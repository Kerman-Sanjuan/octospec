---
description: Score the pull request from a list of findings and say whether it needs a human.
argument-hint: "<pr-number|url>"
---

Runs the `review` stage in the current session. The `review` agent defines its scope, tools, and model.

You find; the command scores. Never write the number yourself: it comes from
`octospec review`, which is the same arithmetic in a laptop and in CI.

1. Resolve the pull request from `$ARGUMENTS` (`gh pr view <n> --json number,title,body,files,additions,deletions,headRefName`).
   If it has no issue, stop and say so: a change without an issue has no contract to review.
2. Read the diff (`gh pr diff <n>`) and the change it belongs to. Judge only against what the
   change declares: the requirements and scenarios in `openspec/changes/<change>/`, the design
   decisions, and the repository's own rules in `AGENTS.md`.
3. Write the findings as JSON to a file. Each finding needs one of the five axes
   (`correctness`, `architecture`, `domain`, `tests`, `ui`), one type (`blocking`, `correctness`,
   `test`, `duplication`, `naming`, `style`, `docs`), the file it is in and the line when there is
   one, and one sentence that a person could check:

   ```json
   {"findings": [
     {"axis": "tests", "type": "test", "file": "cli/internal/lock/lock_test.go", "line": 42,
      "note": "asserts the error message instead of the behaviour, so a rename passes"}
   ]}
   ```

   Report what you can defend. An empty list is a legitimate answer and means the change is clean.
4. Score it:

   ```
   octospec review --findings <file> \
     --changed-files <comma-separated paths> \
     --diff-lines <additions + deletions> \
     [--familiar] \
     --json
   ```

   Pass `--familiar` only when the person who will merge has already touched those paths.
5. Publish exactly one comment, headed `## Review`, containing:
   - the score, the threshold and the verdict, on the first line;
   - the costed findings grouped by axis, each with `file:line` and the sentence;
   - the risk level and the facts that produced it;
   - a closing line stating that a process, not a person, produced this score.

   Update the comment in place on re-run: keep its id in `.openspec.yaml` as
   `github.comments.review`.
6. If the verdict is `human`, leave the pull request unmerged and say which axis cost the most. If
   it is `pass`, say so and stop. Do not merge: merging is a person's decision.
7. If the command failed because a finding named an unknown axis or type, fix the finding, not the
   command.

The repository's numbers live in `.octospec/review.json`. Changing the threshold is an edit there,
never a release of octospec.
