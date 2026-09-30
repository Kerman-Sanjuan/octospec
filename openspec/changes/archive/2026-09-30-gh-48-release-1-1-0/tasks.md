## 1. Version and changelog

- [ ] 1.1 Bump the CLI version to `1.1.0` (`cli/cmd/octospec/main.go`)
- [ ] 1.2 Add the `1.1.0` changelog entry, grouped by impact (`CHANGELOG.md`)
- [ ] 1.3 Make the install check version-agnostic (`openspec/changes/gh-48-release-1-1-0/specs/release-process/spec.md`)

## 2. Release

- [x] 2.1 Tag `v1.1.0` and push it, then watch the release workflow (`docs/releasing.md`)
- [x] 2.2 Verify the release has binaries and a checksum file (GitHub release)

## 3. Verify the install

- [x] 3.1 Run the `curl | sh` installer and confirm `octospec version` reports `1.1.0` (`install.sh`)
- [x] 3.2 Confirm the installed binary writes the agents in a scratch repo (`cli/internal/payload/agents/`)
