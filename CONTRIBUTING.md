# Contributing

octospec makes OpenSpec GitHub-native. GitHub Issues are the backlog and the human source of truth, the repository holds the machine artifacts under `openspec/`, and CI enforces the gates.

This repository dogfoods itself: the workflow in `cli/internal/payload/commands/` is the workflow used here. `AGENTS.md` is the guide for agents working on the repo; this file is for humans.

## The loop

1. `/idea` files the issue from the feature form.
2. `/spec` turns the issue into an OpenSpec change on a branch, commits the artifacts, and publishes one `## Plan` comment.
3. `/apply` works the tasks, commits group by group, and keeps one `## Implementation` comment.
4. `/ship` opens the pull request once the gates are green.
5. `/archive` archives the change, syncs the specs, and closes the issue.

## Branches

Name the branch `feat|fix/<issue>-<slug>`. The prefix follows the issue label: `type:feature` is `feat`, `type:bug` is `fix`. The G3 gate checks this.

## Commits

Use `<type>(#<issue>): <summary>`, for example `feat(#50): verify the installer checksum`. The conventional prefixes feed the generated release notes.

## Run the checks

```sh
cd cli && gofmt -l . && go vet ./... && go test ./...
sh scripts/install-test.sh
openspec validate --all --strict
sh scripts/check-gates.sh
```

CI runs the same on every pull request, in the `cli` and `gates` jobs. Both are required on `main`.

## Cut a release

See [docs/releasing.md](docs/releasing.md). In short: update `CHANGELOG.md`, bump `const version` in `cli/cmd/octospec/main.go`, land it through a pull request, then tag and push. The release workflow publishes the binaries, and the acceptance workflow verifies the published binary on linux and macOS.

## Where things live

- `cli/internal/payload/commands/` - the canonical workflow commands.
- `cli/internal/payload/agents/` - the canonical stage agents.
- `openspec/specs/` - the living contracts.
- `openspec/changes/archive/` - the shipped changes.
- `docs/` - the guides.
- `AGENTS.md` - the guide for agents working on the repo.
