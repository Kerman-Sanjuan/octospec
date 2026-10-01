## Why

octospec installs pi and opencode globally, so their commands and agents leak into every project. Claude Code and Copilot are repo-local, so the rule is inconsistent. There is no warning and no choice. This change makes the scope a destination-based choice: repo-local by default, global on request, with a clear warning.

## What Changes

- Offer a repo-local and a global location for every tool, and default to repo-local.
- Warn when a tool installs globally, explaining that its agents and skills appear in every project.
- Add an explicit global flag for non-interactive installs.
- Record each tool's scope in `.octospec/octospec.json` and re-apply it on update.
- Stop installing agent files for pi, which has no agent mechanism.
- Offer to remove the files that earlier versions installed globally.

## Capabilities

### New Capabilities
- `install-scope`: the repo-local or global choice per tool, the global warning, and the cleanup of earlier global files.

### Modified Capabilities
- `cli-distribution`: every tool maps both a repo-local and a global location, a repo-local install writes nothing under home, and pi has no agent files.
- `consumer-docs`: the command reference documents the scopes and the warning.

## Impact

- `cli/internal/targets/`, `cli/internal/config/`, `cli/internal/plan/`, `cli/internal/install/`, `cli/internal/update/`, `cli/cmd/octospec/main.go`.
- `README.md` and `docs/commands.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/54
