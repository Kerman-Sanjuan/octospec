# Tasks — G9

## 1. The check

- [x] 1.1 G9 reads the requirement titles from every `specs/*/spec.md` in the change.
- [x] 1.2 It reads `coverage.md` and reports a requirement with no line.
- [x] 1.3 It reports a symbol that does not exist outside `openspec/`.
- [x] 1.4 A change with no spec delta skips.

## 2. Everywhere else

- [x] 2.1 `/apply` creates `coverage.md` as it works the tasks.
- [x] 2.2 `README.md` lists G9 and its strength.

## 3. Verification

- [x] 3.1 Run the gate on a change that is fully covered: PASS.
- [x] 3.2 Run it on a change with a missing line and a fake symbol: the warnings name both.
- [x] 3.3 `scripts/e2e-test.sh` and `go test ./...` pass.
