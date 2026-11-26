# Tasks — budgets

## 1. The policy

- [ ] 1.1 `budget.Limits` with the five fields and conservative defaults.
- [ ] 1.2 `Load(repo)` reads `.octospec/budget.json`, per-stage overrides, defaults for the rest.
- [ ] 1.3 An unreadable file is an error, not a silent default.

## 2. The check

- [ ] 2.1 `Usage` with elapsed, tokens, files, diff lines and tasks.
- [ ] 2.2 `Check` returns the limit that was hit, or none, and the remaining allowance.
- [ ] 2.3 Negative tokens mean unknown: reported, never blocking, never zero.

## 3. The command

- [ ] 3.1 `octospec budget --stage <s>` with usage flags, `--json`, exit 6 when exhausted.
- [ ] 3.2 `octospec budget --plan` prints the limits for a wrapper to set timeouts.

## 4. Everywhere else

- [ ] 4.1 `seed` writes the default file.
- [ ] 4.2 `doctor` prints the effective budgets.
- [ ] 4.3 `/apply` and `/ship` check between tasks and stop with `budget_exhausted`.
- [ ] 4.4 `docs/commands.md`.

## 5. Verification

- [ ] 5.1 Tests: under budget, each limit, unknown tokens, unreadable file, exit code.
- [ ] 5.2 `go test ./...` and `scripts/e2e-test.sh`.
