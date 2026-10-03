## 1. Version injection

- [ ] 1.1 Change `const version` to `var version = "dev"` in `cli/cmd/octospec/main.go`.
- [ ] 1.2 Add `ldflags` to the `octospec` build in `.goreleaser.yaml` to inject `-X main.version={{ .Version }}`.
- [ ] 1.3 Add a test that the version variable is settable and the command prints it.

## 2. Checks

- [ ] 2.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.
- [ ] 2.2 Verify the injection locally: build with `-ldflags "-X main.version=v1.0.0"` and confirm `octospec version` prints `1.0.0`, and a plain build prints `dev`.
