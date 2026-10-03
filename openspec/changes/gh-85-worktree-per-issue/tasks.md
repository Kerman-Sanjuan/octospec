## 1. Session registry and worktrees

- [ ] 1.1 Add `cli/internal/session/session.go`: resolve the shared git dir
      (`git rev-parse --git-common-dir`), read and write
      `octospec/sessions.json`, and serialize writes with an exclusive lock
      file.
- [ ] 1.2 In `cli/internal/session/session.go`, add `Start(issue, tool)`: derive
      `<issue>-<slug>` from the issue title, create the worktree at
      `.worktrees/<issue>-<slug>` on `feat|fix/<issue>-<slug>`, record the
      session (branch, worktree, pid, host, stage, start time), and warn when the
      issue already has a live session.
- [ ] 1.3 In `cli/internal/session/session.go`, add `List()` and
      `End(issue)`: list issue, branch, worktree, pid, and age; remove the
      worktree and unregister the session.
- [ ] 1.4 In `cli/internal/session/session.go`, add the launch path: with a
      tool, exec the tool's agent inside the worktree; without it, print the
      path.

## 2. CLI surface

- [ ] 2.1 Add `newSessionCmd` with `start`, `list`, and `end` subcommands in
      `cli/cmd/octospec/main.go`, and register it on the root command.
- [ ] 2.2 Add `.worktrees/` to `.gitignore` and to
      `cli/internal/seed/repo/` if the seed ships a gitignore.

## 3. Session-aware stages

- [ ] 3.1 Add a shared session-check step to
      `cli/internal/payload/commands/spec.md`: if a session exists, run in its
      worktree; if the checkout is on another issue's branch, stop.
- [ ] 3.2 Apply the same session-check step to
      `cli/internal/payload/commands/apply.md`, `ship.md`, and `archive.md`.
- [ ] 3.3 Add the bounded fetch-rebase-push retry to
      `cli/internal/payload/commands/archive.md`, stopping on a rebase conflict.

## 4. Docs and tests

- [ ] 4.1 Document `session start|list|end` and the worktree model in
      `README.md` (Install or a new Parallel work section) and
      `docs/commands.md`.
- [ ] 4.2 Add `cli/internal/session/session_test.go`: registry round-trip,
      concurrent write, duplicate-session warning, and worktree creation.
- [ ] 4.3 Run `gofmt -l . && go vet ./... && go test ./...` in `cli/`,
      `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.
