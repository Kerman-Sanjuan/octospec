---
description: Run CI, open a pull request linked to the issue, and iterate failures back to apply.
argument-hint: "[change]"
---

Runs the `ship` stage in the current session. The `ship` agent defines its scope, tools, and model.

Open the PR for the current branch once the gates are green.

1. Read `.openspec.yaml` for `github.issue`.
2. Push the branch: `git push -u origin HEAD`.
3. Create or update the pull request. The body MUST contain `Closes #<issue>`,
   then:
   - a short summary,
   - the change name and artifact paths,
   - the `openspec validate` result.
   `gh pr create --title "<type>(#<issue>): <summary>" --body-file <file>`
4. Copy the issue's kind and mode onto the pull request, verbatim:

   ```
   gh issue view <issue> --json labels --jq '[.labels[].name | select(startswith("type:") or . == "mode:autonomous") | "--add-label=" + .] | .[]'
   ```

   Pass those arguments to `gh pr edit <pr>`. If the issue carries no `type:*` label, stop and name
   the issue: a pull request with no kind cannot be filtered, and guessing it from the title is how
   labels stop meaning anything. Do not copy `status:*`: status describes the issue's journey, and
   the pull request already says it is under review by existing.
5. Run or observe CI (`gh pr checks --watch`). If any gate fails, hand the
   specific failing check back to `/apply` to fix and iterate - do not report a
   bare failure, and do not continue until it is green.
6. Once every gate is green, add label `status:in-review` to the **issue**; leave
   `status:in-progress` (the `/archive` step clears it).
7. Write the result you are handing over:

   ```
   octospec result write --stage ship --issue <issue> --change <change> --outcome ok \
     --artifact <pull request url>
   ```

   An `ok` from this stage without a pull request is refused on purpose: it is what a stage writes
   when it returned without working. If CI never goes green, write
   `--outcome stalled --reason timeout`; if the change should not ship after all,
   `--reason no_pr_requested`.
8. Print the PR URL.
