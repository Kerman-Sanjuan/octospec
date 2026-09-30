# Issue #50: Harden distribution: checksum verification, version pinning, and a release acceptance test

https://github.com/Kerman-Sanjuan/octospec/issues/50

## User story

As a maintainer of octospec, I want the distribution to be verified and hardened (an installer that checks the checksum and accepts a version, a contributing guide for humans, and an automated acceptance test of the published binary on every release), so that delivery does not depend on a manual check and a user never receives a broken binary.

## Context / problem

The distribution path has no safety net. `curl | sh` fetches `releases/latest` without verifying the checksum, and it accepts no version, so a user cannot pin. The only end-to-end acceptance was a manual smoke repo, now deleted, so nothing checks that a published binary runs. There is no contributing guide for humans, and `goreleaser` is minimal. A broken release would only surface when a user tries to install it.

## Requirements

- The installer SHALL verify the release `checksums.txt` before it installs the binary.
- The installer SHALL accept a version to install, so a user can pin (for example `--version v1.1.0`). This absorbs #37.
- A release workflow SHALL run an acceptance test against the published binary: on linux and macOS, download it, run `octospec version`, run `octospec install` and `octospec seed --no-labels` in a temporary directory, and assert the commands, the 5 agents, and the schema.
- `goreleaser` SHALL generate the release notes from the commits.
- The repository SHALL have a `CONTRIBUTING.md` that explains the loop, the branch naming, how to run the tests, and how to cut a release.
- The docs SHALL list the supported platforms and every install channel.

## Success criteria

- `curl | sh` refuses to install when the checksum does not match.
- `curl | sh --version v1.1.0` installs that exact version.
- After a tag, the release workflow downloads the published binary on linux and macOS, and `octospec version`, `octospec install`, and the file assertions all pass.
- The GitHub release carries generated release notes.
- `CONTRIBUTING.md` exists and a newcomer can follow it to send a change.
- The README install section lists the supported platforms and the channels.
- A broken binary fails the release workflow instead of reaching the user.

## Out of scope

- Homebrew and other package managers.
- Security scanning (Dependabot, CodeQL), tracked in #39.
- Windows support.
- Label provisioning against a real repository (the A2 level).
- Automating the full agent loop (`/idea` to `/archive`).
