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

An empty model means the tool default.

## 6. Run one change

Start in your agent (pi, opencode, GitHub Copilot, or Claude Code):

1. `/idea` to file the feature as an issue.
2. `/spec <issue>` to turn the issue into a change and publish the plan.
3. `/apply` to implement it on the change branch.
4. `/ship` to open the pull request once CI is green.
5. `/archive` after the PR merges, to sync the specs and close the issue.

## Next

- [Commands](commands.md) for the full reference.
- [Troubleshooting](troubleshooting.md) when something breaks.
