# Contributing to octospec

We would love for you to contribute to octospec and help make it better. As a contributor, here are the guidelines we would like you to follow:

- [Code of Conduct](#code-of-conduct)
- [Got a question or a problem?](#got-a-question-or-a-problem)
- [Found a bug?](#found-a-bug)
- [Missing a feature?](#missing-a-feature)
- [Use of AI assistance](#use-of-ai-assistance)
- [Submitting an issue](#submitting-an-issue)
- [Development setup](#development-setup)
- [The workflow loop](#the-workflow-loop)
- [Branches and commits](#branches-and-commits)
- [Coding rules](#coding-rules)
- [Running the checks](#running-the-checks)
- [Submitting a pull request](#submitting-a-pull-request)
- [Reviewing a pull request](#reviewing-a-pull-request)
- [After your pull request is merged](#after-your-pull-request-is-merged)
- [Cutting a release](#cutting-a-release)
- [Where things live](#where-things-live)
- [License](#license)

octospec makes OpenSpec GitHub-native. GitHub Issues are the backlog and the human source of truth, the repository holds the machine artifacts under `openspec/`, and CI enforces the gates. This repository dogfoods itself: the workflow in `cli/internal/payload/commands/` is the workflow used here. `AGENTS.md` is the guide for agents working on this repository; this file is for humans.

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Be respectful and constructive. Harassment, personal attacks, and dismissive behaviour are not welcome. Report unacceptable behaviour to the maintainer.

## Got a question or a problem?

Do not open an issue for general support. Issues are for bug reports and feature requests. For a question, use GitHub Discussions on this repository, or ask in a fork. If the question turns out to be a real bug, open an issue once you have a reproduction.

## Found a bug?

Before you open an issue:

1. Search the open and closed issues for a duplicate. A closed issue often answers the question.
2. Reproduce the bug on the latest release. Reinstall with `curl | sh`, or run `octospec update`.
3. Record the version (`octospec version`), your OS and architecture, and the exact steps.
4. Open a bug issue from the bug form, with a minimal reproduction and the expected and actual behaviour.

Even better: send a pull request with the fix, linked to the issue.

## Missing a feature?

You can request a feature by opening an issue from the feature form.

If you want to implement it:

- For a **large feature**, open the issue first and outline the approach. This avoids duplicate work and lets us agree on the design before you build it.
- **Small features** can go straight to a pull request.

## Use of AI assistance

octospec is built for agents, so AI-assisted contributions are welcome. Two rules:

1. **Disclose it.** Say in the pull request that you used an AI assistant, and which one. You are responsible for the result: read it, understand it, and be ready to discuss and revise it in review.
2. **No bulk, queue-driven pull requests.** Do not point an agent at the issue tracker and open patches across unrelated issues. One person, one issue, one change, shepherded through review. Pull requests that look like a queue will be closed.

## Submitting an issue

The issue is the intake artifact and the human source of truth. A good issue has a user story, the context, the requirements (SHALL or MUST), the success criteria, and what is out of scope. The issue form asks for exactly that. Everything after that, the plan and the implementation, lives in comments.

## Development setup

Prerequisites:

- Go 1.23 or newer.
- The [OpenSpec CLI](https://github.com/Fission-AI/OpenSpec) 1.3.1 or newer.
- The [`gh` CLI](https://cli.github.com) 2.x, authenticated (`gh auth status`).
- Git.

Clone and build:

```sh
git clone https://github.com/Kerman-Sanjuan/octospec.git
cd octospec/cli
go build ./...
```

Try your build against a scratch repository:

```sh
cd cli
go build -o /tmp/octospec ./cmd/octospec
cd /path/to/a/scratch/repo
/tmp/octospec install
/tmp/octospec doctor
```

## The workflow loop

1. `/idea` files the issue from the feature form.
2. `/spec` turns the issue into an OpenSpec change on a branch, commits the artifacts, and publishes one `## Plan` comment.
3. `/apply` works the tasks, commits group by group, and keeps one `## Implementation` comment.
4. `/ship` opens the pull request once the gates are green.
5. `/archive` archives the change, syncs the specs, and closes the issue.

Every change lands through a pull request. One issue, one change.

## Branches and commits

Name the branch `feat|fix/<issue>-<slug>`. The prefix follows the issue label: `type:feature` is `feat`, `type:bug` is `fix`. The G3 gate checks this.

Commit as `<type>(#<issue>): <summary>`, for example `feat(#50): verify the installer checksum`. The conventional prefix feeds the generated release notes. Keep commits focused: one group per commit when a change has several.

## Coding rules

- **Go.** Format with `gofmt`, and keep `go vet` clean. Follow the standard library style.
- **Tests.** Every behaviour change comes with a test. The CLI tests live beside the code (`*_test.go`), and the installer is covered by `scripts/install-test.sh`.
- **Keep changes small.** One capability or fix per change. Small diffs are easier to review and revert.
- **Update the docs.** A change to the CLI, the commands, or the workflow updates the README and the relevant spec under `openspec/specs/`.
- **One canonical payload.** Commands and agents live once in `cli/internal/payload/` and render per tool. Never add a per-tool copy.

## Running the checks

```sh
cd cli && gofmt -l . && go vet ./... && go test ./... && go build ./...
sh scripts/install-test.sh
openspec validate --all --strict
sh scripts/check-gates.sh
```

CI runs the same on every pull request, in the `cli` and `gates` jobs. Both are required on `main`.

## Submitting a pull request

The easy path is the loop: `/spec`, then `/apply`, then `/ship`. If you drive it by hand:

1. Work on a branch named `feat|fix/<issue>-<slug>`.
2. Keep the issue as the source of truth. Link the pull request with `Closes #<issue>`.
3. Run the checks above.
4. Open the pull request. The body describes the problem, the change, and how you verified it.

A pull request checklist:

- [ ] It links the issue (`Closes #<issue>`).
- [ ] The branch is named `feat|fix/<issue>-<slug>`.
- [ ] `gofmt -l`, `go vet`, and `go test ./...` are clean.
- [ ] `openspec validate --all --strict` passes.
- [ ] The docs and the spec are updated when behaviour changes.
- [ ] The text follows the `humanize` style: short sentences, no em dashes.

## Reviewing a pull request

Review the diff, the spec delta, and the verification. Comment with specifics, not vibes.

If a reviewer asks for changes, push more commits to the same branch. The pull request updates and CI re-runs. There is no need to open a second pull request.

## After your pull request is merged

Run `/archive` if you own the change. It archives the change, syncs the specs, and closes the issue. Then delete the branch and pull `main`.

## Cutting a release

See [docs/releasing.md](docs/releasing.md). In short: update `CHANGELOG.md`, bump `const version` in `cli/cmd/octospec/main.go`, land it through a pull request, then tag and push. The release workflow publishes the binaries, and the acceptance workflow verifies the published binary on linux and macOS.

## Where things live

- `cli/internal/payload/commands/` - the canonical workflow commands.
- `cli/internal/payload/agents/` - the canonical stage agents.
- `cli/internal/targets/` - the per-tool mapping.
- `openspec/specs/` - the living contracts.
- `openspec/changes/archive/` - the shipped changes.
- `docs/` - the guides.
- `AGENTS.md` - the guide for agents working on this repository.

## License

octospec is released under the [MIT License](LICENSE). By contributing, you agree that your contribution is licensed under the same terms.
