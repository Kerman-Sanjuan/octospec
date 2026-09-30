## 1. Version and changelog

- [ ] 1.1 Bump the CLI version to `1.0.0` (`cli/cmd/octospec/main.go`)
- [ ] 1.2 Finalize `CHANGELOG.md` for `1.0.0` with the release date (`CHANGELOG.md`)

## 2. Release

- [ ] 2.1 Tag `v1.0.0` and push it, then watch the release workflow (`docs/releasing.md`)
- [ ] 2.2 Verify the release has binaries and a checksum file (GitHub release)

## 3. Verify the install

- [ ] 3.1 Run the `curl | sh` installer and confirm `octospec version` reports `1.0.0` (`install.sh`)
- [ ] 3.2 Confirm `go install .../cli/cmd/octospec@latest` installs a working binary (`cli/`)
