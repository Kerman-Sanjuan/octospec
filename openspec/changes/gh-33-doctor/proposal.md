## Why

An install spans several pieces: the OpenSpec schema, the commands per tool, the workflow labels, and the gate script. When one is missing the failure is confusing. `doctor` turns that into a clear list.

## What Changes

- Add `octospec doctor`, which checks the OpenSpec CLI, the schema, the installed commands, the workflow labels, and the gate script, and exits non-zero when something required is missing.

## Capabilities

### New Capabilities
- `diagnostics`: a read-only check of an octospec install.

### Modified Capabilities
<!-- None. -->

## Impact

- `cli/internal/doctor/doctor.go` (new), `cli/cmd/octospec/main.go`, `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/33
