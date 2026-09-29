## Why

The README is a 5-line stub that only links to the design doc; a new user cannot discover what octospec is, how to install it, or how to run the workflow without reading internal plan files. The workflow is now specified and implemented, so the README is the missing front door.

## What Changes

- Explain what octospec is and the problem it solves.
- Describe the workflow loop (`/idea → /spec → /apply → /ship → /archive`).
- Document the architecture (schema override, command shims, seed files, install script, CI gates).
- Provide installation instructions for pi, opencode, and GitHub Copilot.
- List every command and its phase/behavior.
- Link to the design doc and implementation plan as deeper references.

## Capabilities

### New Capabilities
- `readme`: The README that documents octospec's workflow, architecture, installation, and command surface.

### Modified Capabilities
<!-- None; no existing specs change. -->

## Impact

- Only `README.md` changes (documentation).
- No code, APIs, dependencies, or systems affected.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/1
