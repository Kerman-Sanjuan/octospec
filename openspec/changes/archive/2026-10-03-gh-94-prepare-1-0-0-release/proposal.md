## Why

octospec is feature-complete enough to be useful but is not yet presentable as
a first public release: the README has no logo or badges, two pre-polish
releases (v1.0.0 and v1.1.0) sit in the history, and there is no recorded
end-to-end check of the published binary.

## What Changes

- Add a logo asset (SVG and PNG) and display it in the README.
- Add status badges to the README for CI, the latest release, the license, and
  the Go version.
- Remove the existing v1.0.0 and v1.1.0 GitHub releases and their tags, then cut
  the debut release as v1.0.0 with a single dated CHANGELOG section.
- Verify the published binary end to end: fresh clone, `curl | sh` install,
  `octospec seed`, `octospec install`, one full change through the loop, and
  `octospec version` matching the tag.
- Fill in the repository About with a description, topics, and the project link.

## Capabilities

### New Capabilities
- `release-debut`: reset the release history to one polished v1.0.0 and prove
  the published binary works end to end.

### Modified Capabilities
- `readme`: the README now displays a logo and status badges.

## Impact

- README.md, a new logo asset under the repository, CHANGELOG.md, and
  docs/releasing.md.
- GitHub releases and tags (v1.0.0, v1.1.0 removed, v1.0.0 re-cut) and the
  repository About.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/94

## Changelog

### Added
- A logo and CI, release, license, and Go badges in the README.

### Changed
- The release history is a single v1.0.0 debut; the earlier v1.0.0 and v1.1.0
  releases and tags were removed.
