## Why

`install.sh --repo` seeds the issue forms, prompts, CI workflow and config, but not the six labels the workflow uses. On a fresh repo, `/idea`, `/spec` and `/apply` fail with `'status:...' not found` until the labels are created by hand.

## What Changes

- Provision the six workflow labels from `install.sh --repo`: `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, `status:in-review`.
- Make provisioning idempotent so re-running the seed does not fail on existing labels.
- Keep the label set defined in one place.

## Capabilities

### New Capabilities
- `label-provisioning`: seeding a repository creates the workflow's `type:*` and `status:*` labels.

### Modified Capabilities
<!-- None; openspec/specs/ has no label capability today. -->

## Impact

- `install.sh` (the `--repo` branch) and a new `scripts/seed-labels.sh`.
- `README.md` install section (describe what `--repo` seeds).
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/12
