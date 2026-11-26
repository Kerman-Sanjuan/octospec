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
4. Score it, and keep the record:

   ```
   octospec review --findings <file> \
     --changed-files <comma-separated paths> \
     --diff-lines <additions + deletions> \
     [--familiar] \
     --ledger .octospec/ledger --pr <pr> --issue <issue> --head <head commit> \
     --json
   ```

   Pass `--familiar` only when the person who will merge has already touched those paths. The
   ledger is the point of the run: the comment dies with the pull request, the record does not.
5. Move the `risk:high` label so the pull request list says where to look:

   ```
   gh pr edit <pr> --add-label risk:high      # when risk is high
   gh pr edit <pr> --remove-label risk:high   # when it is not
   ```

   The label points at the record; it never replaces it. A label is current state and disappears
   with the pull request, and what we want here is a series.
6. Publish exactly one comment, headed `## Review`, containing:
   - the score, the threshold and the verdict, on the first line;
   - the costed findings grouped by axis, each with `file:line` and the sentence;
   - the risk level and the facts that produced it;
   - a closing line stating that a process, not a person, produced this score.

   Update the comment in place on re-run: keep its id in `.openspec.yaml` as
   `github.comments.review`.
7. Commit the record (`git add .octospec/ledger && git commit -m "chore(#<issue>): record the
   review"`), so it lands with the merge instead of dying on the branch.
8. If the verdict is `human`, leave the pull request unmerged and say which axis cost the most. If
   it is `pass`, say so and stop. Do not merge: merging is a person's decision.
9. If the command failed because a finding named an unknown axis or type, fix the finding, not the
   command.

10. Write the result, so a scheduler can tell a clean review from one that never happened:

    ```
    octospec result write --stage review --issue <issue> --change <change> --outcome ok \
      --artifact <comment url>
    ```

    With `--outcome stalled --reason needs_human` when the verdict asks for a person.

The repository's numbers live in `.octospec/review.json`. Changing the threshold is an edit there,
never a release of octospec. The series lives in `.octospec/ledger/` and is read with
`octospec ledger`.
