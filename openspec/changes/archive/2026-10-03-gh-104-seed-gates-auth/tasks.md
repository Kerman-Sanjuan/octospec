## 1. Gate robustness

- [x] 1.1 In `scripts/check-gates.sh`, capture the `gh issue view` exit code and SKIP G1 when it fails, instead of treating an empty body as incomplete.
- [x] 1.2 Copy the change to `cli/internal/seed/repo/scripts/check-gates.sh` so the two stay byte-identical.

## 2. Seed workflow

- [x] 2.1 Add `GH_TOKEN: ${{ github.token }}` and `issues: read` to the seeded `.github/workflows/openspec.yml`, and bump `setup-node` to match the root.

## 3. Tests

- [x] 3.1 Add a case to `scripts/check-gates-test.sh`: a failing `gh` makes G1 SKIP and the run pass.
- [x] 3.2 Add a seed test that fails when the seeded workflow drifts from the root on the token, the permissions, or the `setup-node` major.

## 4. Checks

- [x] 4.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, `sh scripts/check-gates.sh`, and `sh scripts/check-gates-test.sh`.
- [x] 4.2 Reproduce the false negative before and confirm it is gone after.
