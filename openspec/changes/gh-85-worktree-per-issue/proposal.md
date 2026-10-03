## Why

Every octospec stage mutates one checkout, so two agents on two different issues
fight over one HEAD and one index. A developer cannot run several agents at once
without them clashing.

## What Changes

- Add `octospec session start <issue> [--tool <tool>]`: create a git worktree
  per issue, register it, and optionally launch the tool's agent inside it.
- Add `octospec session list` and `octospec session end <issue>`.
- Store the registry under the shared git dir
  (`$(git rev-parse --git-common-dir)/octospec/sessions.json`), never committed,
  with serialized writes.
- Make `/spec`, `/apply`, `/ship`, and `/archive` session-aware: run in the
  session's worktree when one exists, behave as today when none does.
- Add a bounded fetch-rebase-push retry to `/archive` so two independent
  archives do not lose to a race on `main`.
- Keep sessions optional; every stage works with no session.

## Capabilities

### New Capabilities
- `parallel-sessions`: one git worktree per issue, a shared registry, optional
  agent launch, and session-aware stages.

### Modified Capabilities
- `archive-flow`: `/archive` gains a bounded fetch-rebase-push retry before it
  pushes the archive, and stops on a rebase conflict.

## Impact

- New package `cli/internal/session` and a `session` command in
  `cli/cmd/octospec/main.go`.
- `cli/internal/payload/commands/{spec,apply,ship,archive}.md` learn to detect a
  session and work in its worktree.
- `.gitignore` ignores `.worktrees/`.
- No behaviour change when no session is used. Link:
  https://github.com/Kerman-Sanjuan/octospec/issues/85

## Changelog

### Added
- `octospec session start|list|end` to run several agents at once, each in its
  own git worktree, without clashing.
