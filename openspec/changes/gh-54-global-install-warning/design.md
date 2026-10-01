## Context

`targets.Target` has one scope (global or repo-local) and one directory per kind. pi and opencode are global; Claude Code and Copilot are repo-local. `install` writes whatever `plan` produces and records the file hashes in `.octospec/octospec.json`. There is no scope choice and no warning. pi gets agent files in `~/.pi/agent/agents/`, which pi does not read, so they are dead and global.

## Goals / Non-Goals

**Goals:**
- Repo-local is the default for every tool.
- The user can choose global, with a clear warning.
- The choice is recorded and re-applied.
- Earlier global files can be removed.

**Non-Goals:**
- Change the commands or the agents.
- A new config file.
- Windows paths.

## Decisions

- **The warning depends on the destination, not the tool.** Any install that writes to a global user directory warns. *Alternative rejected:* warn per tool, which is inconsistent and misses a tool that later gains a global location.
- **Repo-local is the default for every tool.** *Alternative rejected:* keep the current per-tool default, which leaks.
- **Both locations per tool.** Each target gains a repo-local directory and a global directory for commands and for agents. *Alternative rejected:* keep one fixed scope per tool, which is what the issue complains about.
- **Interactive choice on a terminal, a flag otherwise.** On a terminal, a `huh` prompt asks the scope; without a terminal, repo-local is the default and `--global` selects global. *Alternative rejected:* always prompt, which hangs in CI.
- **Record the scope in the state.** Add a `scopes` map (tool to `local` or `global`) to `.octospec/octospec.json`; `update` re-applies it. *Alternative rejected:* infer the scope from the files, which is brittle.
- **pi gets no agent files.** pi has no agent mechanism, so its agent location is empty and the planner skips it. *Alternative rejected:* keep writing dead files.
- **Cleanup reads the recorded files.** On install or update, the previous `files` map is inspected for paths under the home directory and offered for removal. *Alternative rejected:* a manual removal only.

## Risks / Trade-offs

- [The Copilot and Claude user-scope paths need confirmation] -> Verify against each tool's docs before coding; the design lists them as targets.
- [Moving every default to repo-local changes behavior for existing users] -> `update` re-applies the recorded scope, so an existing global install stays global until the user chooses otherwise, and cleanup is offered.
- [The cleanup could remove a file the user edited] -> It uses the recorded managed-file hashes, and it is an explicit, warned choice.

## Migration Plan

A fresh `install` defaults to repo-local. An existing `.octospec/octospec.json` keeps its recorded scope through `update`. When the user moves a tool to repo-local, octospec offers to remove the old global files.

## Open Questions

- Is the scope prompt once for all tools, or one prompt per tool?
- The exact Copilot user-scope directories.
