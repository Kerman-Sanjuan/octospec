# Issue #94: Prepare the polished 1.0.0 public release

https://github.com/Kerman-Sanjuan/octospec/issues/94

## User story

As a maintainer, I want octospec packaged and presented as a polished 1.0
public release, so that a newcomer sees a credible, attractive project and can
install and run it end to end.

## Context / problem

octospec already has enough features to be useful, but it is not yet presentable
as a first public release. There are two prior GitHub releases (v1.0.0 and
v1.1.0, both dated 2026-09-30) that predate the polish. The README has no logo
and no badges, so the project's landing page does not communicate what it is or
that it is healthy. There is no recorded end-to-end verification of the published
binary. This issue prepares the debut release: present it well, prove it works,
and reset the release history to a single clean version.

## Requirements

- The repository SHALL contain a logo asset (SVG and PNG) under `docs/` or an
  `assets/` directory, and the README SHALL display it.
- The README SHALL show status badges for CI, the latest release, the license,
  and the Go version.
- The maintainers SHALL remove the existing v1.0.0 and v1.1.0 GitHub releases
  and their tags before the debut.
- The debut release SHALL be tagged v1.0.0.
- The release SHALL pass an end-to-end check on the published binary: a fresh
  clone, the `curl | sh` install, `octospec seed`, `octospec install` into a
  tool, one full change through the loop, and `octospec version` matching the
  tag.
- The CHANGELOG SHALL present a single 1.0.0 section for the debut, generated
  from the archived changes.
- The README SHALL link the docs, the changelog, and the About section, and the
  repository About SHALL be filled in with a description, topics, and the
  project link.

## Success criteria

- A visitor opens the repository and sees a logo and green badges.
- The release page lists exactly one release, v1.0.0, with binaries, a
  checksums file, and release notes from the changelog.
- `curl | sh` on a clean machine installs octospec, and `octospec version`
  prints 1.0.0.
- A newcomer can follow the README and run one change from idea to archive.
- `openspec validate --all --strict`, `scripts/check-gates.sh`, and the `cli`
  and `gates` CI checks all pass.

## Out of scope

- No new CLI features and no behavior changes.
- No changes to the workflow commands, the agents, or the OpenSpec schema.
- No new tools.
