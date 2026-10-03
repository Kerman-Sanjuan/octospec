## 1. Logo and badges

- [x] 1.1 Add `docs/logo.svg`, a simple geometric mark for octospec.
- [x] 1.2 Render and commit `docs/logo.png` from the SVG.
- [x] 1.3 Display the logo at the top of `README.md`.
- [x] 1.4 Add shields.io badges to `README.md` for CI, the latest release, the license, and the Go version.

## 2. README links

- [x] 2.1 Ensure `README.md` links the docs, `CHANGELOG.md`, and the About link, without changing existing prose.

## 3. Release reset

- [ ] 3.1 Remove the v1.1.0 GitHub release and tag, and the pre-polish v1.0.0 release and tag.
- [ ] 3.2 Fold the `CHANGELOG.md` 1.1.0 section into a single dated 1.0.0 section.

## 4. Release cut

- [ ] 4.1 Cut and push the `v1.0.0` tag and confirm the release workflow publishes binaries and checksums.
- [ ] 4.2 Run the end-to-end check on the published binary: fresh clone, `curl | sh`, `seed`, `install`, one change through the loop, and `octospec version`.
- [ ] 4.3 Fill in the repository About (description, topics, project link).

## 5. Checks

- [ ] 5.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.
- [ ] 5.2 Update `docs/releasing.md` if the reset changes any step, and confirm the `cli` and `gates` checks are green on the pull request.
