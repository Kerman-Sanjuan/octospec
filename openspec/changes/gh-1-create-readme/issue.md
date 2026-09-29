# Issue #1: Create the README documenting the GitHub-native OpenSpec workflow

https://github.com/Kerman-Sanjuan/octospec/issues/1

## User story

As a developer adopting octospec, I want a README that explains the GitHub-native OpenSpec workflow, so that I can install it and follow the `/idea → /spec → /apply → /ship → /archive` loop without reading the design docs.

## Context / problem

The README today is a 5-line stub that only links to the design doc and defers full documentation to "Task 9". It is not a proper README: a new user has no way of knowing what octospec is, its architecture, how it works, how to install it, or what commands exist. The workflow is now specified and implemented, so the README is the missing front door.

## Requirements

- The README SHALL explain what octospec is and the problem it solves.
- The README SHALL describe the workflow loop (`/idea → /spec → /apply → /ship → /archive`).
- The README SHALL document the architecture (schema override, command shims, seed files, install script, CI gates).
- The README SHALL provide installation instructions for pi, opencode, and GitHub Copilot.
- The README SHALL list every command and its phase/behavior.
- The README SHALL link to the design doc and implementation plan as deeper references.

## Success criteria

- A first-time reader can state what octospec is and the workflow loop without opening the design docs.
- The README contains working install instructions for pi, opencode, and Copilot that a fresh clone can follow.
- Every command (`/idea`, `/bug`, `/explore`, `/spec`, `/apply`, `/ship`, `/archive`) is listed with its phase and behavior.
- The README describes the architecture and points to deeper docs for details.
- Links in the README resolve to existing files.

## Out of scope

- Replacing the design doc or implementation plan with the README (they remain the deep references).
- Writing a multi-page docs site; the README is the entry point only.
- Tutorials or onboarding guides beyond the install + workflow loop.
