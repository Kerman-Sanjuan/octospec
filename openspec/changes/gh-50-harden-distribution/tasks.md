## 1. Harden the installer

- [ ] 1.1 Download and verify the release checksum before installing (`install.sh`)
- [ ] 1.2 Add a version flag that pins the release and keep latest as the default (`install.sh`)
- [ ] 1.3 Document the checksum limit and the pinning usage (`docs/releasing.md`)

## 2. Verify the release

- [ ] 2.1 Add a release acceptance workflow for linux and macOS that downloads the published binary (`.github/workflows/release-acceptance.yml`)
- [ ] 2.2 Assert `octospec version`, the commands, the 5 agents, and the schema in a temporary directory (`.github/workflows/release-acceptance.yml`)
- [ ] 2.3 Generate the release notes from the commits in `goreleaser` (`.goreleaser.yaml`)

## 3. Contributing guide and docs

- [ ] 3.1 Add `CONTRIBUTING.md` covering the loop, the branch naming, the tests, and the release (`CONTRIBUTING.md`)
- [ ] 3.2 List the supported platforms and the install channels in the README (`README.md`)

## 4. Tests and gates

- [ ] 4.1 Add a shell test for the installer checksum and pinning paths (`scripts/install-test.sh`)
- [ ] 4.2 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
