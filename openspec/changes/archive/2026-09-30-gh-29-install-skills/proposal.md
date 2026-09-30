## Why

Every command says "load the `openspec-propose` skill", but octospec installs only the commands. On a fresh repo, `/spec` fails at the first skill load, so the loop cannot start. The skills must reach each tool.

## What Changes

- `octospec install` ensures the OpenSpec skills are present for the selected tools by running `openspec init --tools <list>`, which is OpenSpec's own, non-interactive generator.
- The tool names map to OpenSpec's: `pi`, `opencode`, `copilot` to `github-copilot`, and `claude`.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `cli-distribution`: `install` now also brings the OpenSpec skills into each tool's skill directory.

## Impact

- `cli/internal/install/install.go` and a new `cli/internal/skills` helper.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/29
