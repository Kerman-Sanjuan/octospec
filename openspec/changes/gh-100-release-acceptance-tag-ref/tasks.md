## 1. Gate fix

- [x] 1.1 In `scripts/check-gates.sh`, check G3 only for `feat/`, `fix/`, and `dependabot/` refs; SKIP any other ref (main, a tag, or another long-lived ref).
- [x] 1.2 Copy the change to the seeded gate script `cli/internal/seed/repo/scripts/check-gates.sh`.

## 2. Test and docs

- [x] 2.1 Add a case to `scripts/check-gates-test.sh` that runs G3 with a tag ref and asserts a SKIP.
- [x] 2.2 Pin the README install example to `v1.0.0`.

## 3. Checks and verification

- [x] 3.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, `sh scripts/check-gates.sh`, and `sh scripts/check-gates-test.sh`.
- [x] 3.2 Verify `HEAD_REF=v1.0.0 sh scripts/check-gates.sh` reports SKIP G3 and does not fail.
- [ ] 3.3 After merge, dispatch `release-acceptance` with `v1.0.0` and confirm the check is green.
