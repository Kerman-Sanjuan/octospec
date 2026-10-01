# Issue #54: Warn on global installs and offer a repo-local location for every tool

https://github.com/Kerman-Sanjuan/octospec/issues/54

## User story

As a user of octospec, I want to choose between a global and a repo-local install for every tool, and to be warned when a global install makes octospec's agents and skills appear in every project, so I can keep the workflow scoped to the current repository.

## Context / problem

octospec installs pi and opencode globally, so their commands and agents show up in every project and leak the workflow into unrelated work. Claude Code and GitHub Copilot are already repo-local, but there is no warning and no choice about the destination, so the rule is inconsistent across tools. pi also gets agent files in `~/.pi/agent/agents/`, a directory pi does not read, so those files are dead and global at the same time.

## Requirements

- `octospec install` SHALL warn when a tool's destination is a global user directory, and SHALL explain that the agents and skills then appear in every project.
- octospec SHALL offer both a repo-local and a global location for every supported tool that has a global config: pi, opencode, Claude Code, and GitHub Copilot.
- The repo-local location SHALL be the default for every tool.
- A non-interactive install SHALL NOT hang, and SHALL default to repo-local.
- The chosen locations SHALL be recorded in `.octospec/octospec.json`, so `octospec update` re-applies the same ones.
- pi has no agent mechanism, so its agent files SHALL NOT be installed.
- octospec SHALL detect the files that earlier versions installed globally and SHALL offer to remove them, as an optional, warned step.

## Success criteria

- Running `octospec install` on a terminal warns for any tool whose destination is global, and lets the user pick repo-local.
- A non-interactive install writes repo-local by default and does not hang.
- Choosing global writes to that tool's global directory and shows the warning.
- After a repo-local install, opening pi or opencode in another project shows no octospec commands or agents.
- Claude Code and Copilot default repo-local and warn only if global is chosen.
- The chosen locations are recorded in `.octospec/octospec.json`, so `octospec update` re-applies the same ones.
- The files that earlier versions installed globally can be removed through octospec, with a warning.

## Out of scope

- Changing the commands or the agents themselves.
- A new config file beyond `.octospec/octospec.json`.
- Windows-specific install paths.
