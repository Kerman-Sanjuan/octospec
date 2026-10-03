# Getting started

octospec turns OpenSpec into a GitHub-native workflow: GitHub Issues are the backlog, the repository holds the machine artifacts, and CI enforces the gates.

## 1. Install the CLI

```sh
curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
```

Contributors can use `go install github.com/kerman-sanjuan/octospec/cli/cmd/octospec@latest` instead.

## 2. Install the commands

Run this in your repository. Repeat `--tool`, or omit it to install for all.

```sh
octospec install --tool pi --tool opencode --tool copilot --tool claude --repo .
```

With no `--tool` on a terminal, `install` opens a multi-select of the tools,
pre-checked with your last install. Space toggles a tool, enter confirms, and
choosing none installs for all of them.

`install` also brings the OpenSpec skills into each tool's skill directory.

## 3. Seed the repository

```sh
octospec seed --repo .
```

This installs the OpenSpec schema, writes the issue forms and the CI workflow, and provisions the workflow labels.

## 4. Check the install

```sh
octospec doctor
```

Every line reports OK, or it tells you what is missing.

## 5. Choose the models

Each stage runs on a scoped agent, and each role (thinking, implementer,
reviewer) uses one model. Set them with the TUI, or by flag:

```sh
octospec models
octospec models --set thinking=sonnet --set implementer=sonnet
```

The TUI lists the known models for your installed tools, and `Custom...` lets
you type any other identifier. `(tool default)` leaves the role without a
model; an empty model means the tool default.

## 6. Run one change

Start in your agent (pi, opencode, GitHub Copilot, or Claude Code):

1. `/idea` to file the feature as an issue.
2. `/spec <issue>` to turn the issue into a change and publish the plan.
3. `/apply` to implement it on the change branch.
4. `/ship` to open the pull request once CI is green.
5. `/archive` after the PR merges, to sync the specs and close the issue.

## 7. Run several changes at once

One change at a time is fine to start, and it is the default. When you want to
work several issues in parallel, give each its own session. A session is a git
worktree on the issue's branch, so two agents never share a working tree, HEAD,
or index.

You are the orchestrator: open one tab per issue, and start its session there.

```sh
# tab 1                          # tab 2                          # tab 3
cd /path/to/repo                 cd /path/to/repo                 cd /path/to/repo
octospec session start 85        octospec session start 91        octospec session start 93
opencode                         pi                               claude
/spec 85                         /spec 91                         /spec 93
```

`session start` creates the worktree and, with `--tool`, launches that tool
inside it, so the tab is already rooted in the right checkout. Without `--tool`
it prints the path to open yourself.

Then run the loop (`/spec`, `/apply`, `/ship`, `/archive`) in each tab as
usual. The stages detect the session and stay in its worktree. When you finish
an issue, `/archive` ends the session; you can also end one by hand:

```sh
octospec session list    # issue, branch, worktree, pid, age
octospec session end 85  # remove the worktree and unregister
```

Two rules keep the tabs from colliding:

- Starting the same issue twice stops, because git refuses to check out a branch
  another worktree holds.
- A stage refuses to run on another issue's checkout, so a `/spec` for one issue
  never switches the branch out from under a different task. If it stops, run
  `octospec session start <issue>` and work in that tab.

Sessions are optional. If you never run `octospec session start`, every stage
behaves exactly as in step 6, in the current checkout.

## Next

- [Commands](commands.md) for the full reference.
- [Troubleshooting](troubleshooting.md) when something breaks.
