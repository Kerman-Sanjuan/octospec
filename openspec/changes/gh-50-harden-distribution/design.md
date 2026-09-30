## Context

The installer downloads `releases/latest/download/octospec_<os>_<arch>` and moves it into place. The release workflow publishes four binaries and `checksums.txt` through `goreleaser`, which is minimal today. The only end-to-end acceptance was a manual smoke repo, now deleted. There is no human-facing contributing guide; `AGENTS.md` is for agents.

## Goals / Non-Goals

**Goals:**
- A checksum-verified installer that accepts a version.
- An automated acceptance test of the published binary on each release.
- Generated release notes.
- A contributing guide and accurate platform docs.

**Non-Goals:**
- Homebrew and other package managers.
- Security scanning (#39).
- Windows.
- Label provisioning against a real repo (A2).
- The full agent loop (B).

## Decisions

- **Verify `checksums.txt` with sha256.** The installer downloads `checksums.txt`, computes the hash of the binary, and refuses to install on a mismatch. *Alternative rejected:* cosign or signature verification, which adds tooling and a key to maintain.
- **A `--version` flag on the installer.** The installer resolves the tag URL when given a version and `releases/latest` otherwise. *Alternative rejected:* a separate installer per version.
- **The acceptance test is file-level and token-free.** A workflow job runs on linux and macOS, installs the published binary, and asserts `octospec version`, the commands, the 5 agents, and the schema in a temporary directory. It does not create a repo, so no PAT. *Alternative rejected:* a throwaway GitHub repo per release, which needs a PAT and cleanup.
- **Trigger the acceptance test on `release: published`.** It depends on the assets being present. *Alternative rejected:* run it inside the release job, which mixes build and verify.
- **Generate release notes with `goreleaser` from the commit history.** The repo already uses conventional commit prefixes. *Alternative rejected:* hand-written notes, which drift.
- **`CONTRIBUTING.md` for humans, `AGENTS.md` for agents.** They serve different readers. *Alternative rejected:* one document for both, which serves neither.

## Risks / Trade-offs

- [The checksum and the binary come from the same origin] -> The checksum guards against a corrupted or truncated download, not a compromised release. Note the limit in the docs.
- [A goreleaser changelog needs commit prefixes] -> The repo already uses `<type>(#<issue>): <summary>`. Add a note to `CONTRIBUTING.md`.
- [The acceptance test depends on the release existing] -> Trigger on `release: published`, which fires after the assets are uploaded.

## Migration Plan

None. The installer keeps the latest-by-default behavior.

## Open Questions

- Does the acceptance test cover both latest and a pinned version on every run, or only latest in CI and pinning in a single job? Proposed: cover both.
