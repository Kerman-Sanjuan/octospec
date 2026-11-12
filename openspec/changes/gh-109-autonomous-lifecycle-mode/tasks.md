## 1. Orchestration entry point

- [x] 1.1 Add a new payload command (e.g. `cli/internal/payload/commands/auto.md`) that documents the autonomous mode's contract: select an active issue, run idea → spec → apply → ship → archive in sequence, retry on gate failure, stop on a bounded failure.
- [x] 1.2 Wire the new command into `cli/internal/targets/` so it renders per tool like the other five commands.
- [x] 1.3 Decide and document the retry bound per stage in the command body (`cli/internal/payload/commands/auto.md`).

## 2. Issue selection

- [x] 2.1 Implement issue selection logic (CLI or command instructions) that queries GitHub for issues in an active status (`status:in-progress` or equivalent) via `gh issue list`.
- [x] 2.2 Document the bug-discovery path: how the mode files a new issue through the existing intake when it finds behaviour contradicting `openspec/specs/`, without inventing requirements.

## 3. Stage sequencing and retry

- [x] 3.1 Implement the sequencing loop that runs `/idea`, `/spec`, `/apply`, `/ship`, `/archive` in order for the selected issue, reusing each stage's existing contract unchanged.
- [x] 3.2 Implement the retry behaviour: on a failing gate (G2a, G2b, G3, G4, G6), re-run the current stage instead of advancing, up to the bound from 1.3.
- [x] 3.3 Implement the give-up path: past the retry bound, post a failure comment on the issue (reusing the plan/implementation comment pattern) and stop working that issue.

## 4. CLI invocation

- [x] 4.1 Add a CLI subcommand (e.g. `octospec auto`) in `cli/` that invokes the autonomous mode, so a human can start the overnight run from the terminal without going through a per-tool slash command.
- [x] 4.2 Document the new subcommand's flags (issue selector, retry bound override) in the CLI's help output and in the README **Install**/**Commands** tables.

## 5. GitHub trail parity

- [ ] 5.1 Verify the autonomous run produces the same artifacts a manual run would: status labels, plan comment, implementation comment, pull request, archive state.
- [ ] 5.2 Add or update tests in `cli/` covering issue selection and retry/give-up behaviour.

## 6. Docs

- [ ] 6.1 Update the README **Commands** table with the new autonomous mode command (per AGENTS.md definition of done).
- [ ] 6.2 Update `cli/internal/payload/agents/README.md` if the autonomous mode introduces a new agent definition or changes the stage table.
