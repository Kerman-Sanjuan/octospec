# Issue #66: Decide and document Windows support

https://github.com/Kerman-Sanjuan/octospec/issues/66

## User story

As a Windows user, I want to know whether octospec works, so I do not waste time on an install that cannot run.

## Context / problem

The installer is a POSIX `sh` script, so it does not run on Windows. `goreleaser` can build a Windows binary, but there is no PowerShell installer and no statement of support. The README is silent.

## Requirements

- octospec SHALL either support Windows, with a PowerShell installer and a Windows release binary, or state clearly that Windows is not supported yet.
- The README SHALL state the outcome.

## Success criteria

- If supported, an install path works on Windows and `octospec version` reports the release.
- If unsupported, the README says so and no user is surprised.

## Out of scope

- Windows support in the `sh` installer.
