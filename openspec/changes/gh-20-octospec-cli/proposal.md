## Why

octospec ships the same command file several times — once per tool — and installs them with a shell script tied to a checked-out repo. The copies drift, the repo root is cluttered, and there is no installable artifact. A Go CLI that renders one canonical payload into each tool's directory removes the copies and the `install.sh` dependency. Because it introduces a code artifact, it also needs its own tests, CI/CD, and a `/ship` loop that feeds CI failures back into `/apply`.

## What Changes

- Add a Go CLI, `octospec/cli`, distributed via `curl | sh` (and `go install`), that installs the commands into pi, opencode, GitHub Copilot, and Claude Code from one canonical payload.
- Remove the per-tool copies (`.github/prompts/`, `repo-template/.github/prompts/`, the `.opencode/`/`.pi/` mirrors) and `install.sh`; the CLI generates them.
- Add the CLI's own Go tests and CI (build, `gofmt`, `go vet`, `go test`, plus the octospec gates), and CD that publishes release binaries.
- Make `/ship` run/observe CI and hand failures back to `/apply` to iterate, proceeding only on green.

## Capabilities

### New Capabilities
- `cli-distribution`: the Go CLI — installing/updating the commands for each tool from one canonical payload, and how it is distributed.
- `release-pipeline`: the CLI's tests, PR CI gates, and tagged release binaries.
- `ship-feedback-loop`: `/ship` observing CI and looping failures back to `/apply`.

### Modified Capabilities
<!-- None; these are new capabilities. -->

## Impact

- New: the Go CLI source, tests, and release workflow.
- Removed: `install.sh`, `.github/prompts/`, `repo-template/.github/prompts/`, the `.opencode/`/`.pi/` mirror folders.
- Changed: `commands/ship.md` (+ mirrors) for the feedback loop; `README.md`; `.github/workflows/`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/20
