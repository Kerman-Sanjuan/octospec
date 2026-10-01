## Context

The release publishes linux and darwin binaries. The installer is a POSIX `sh` script. Windows can build a Go binary with `goreleaser`, but there is no PowerShell installer, and the tool wiring (the agent directories) has not been checked on Windows.

## Goals / Non-Goals

**Goals:**
- A Windows user learns the status before trying.

**Non-Goals:**
- Adding Windows support, which needs a PowerShell installer and a Windows binary.

## Decisions

- **Document that native Windows is not supported yet, and do not add support here.** Support means a PowerShell installer, a Windows release binary, and verified tool wiring, which is its own change. WSL already works, because it is Linux, so the linux binary and the `sh` installer run there. *Alternative rejected:* ship a partial native Windows path now, which would be unverified. The README states the status plainly.

## Risks / Trade-offs

- [A Windows user is turned away] -> The status is explicit, so they do not waste time, and support can follow.

## Migration Plan

None.

## Open Questions

- Whether Windows support is worth a PowerShell installer and a Windows binary.
