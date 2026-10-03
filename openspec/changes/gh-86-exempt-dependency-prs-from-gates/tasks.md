## 1. Gate exemption

- [x] 1.1 In `scripts/check-gates.sh`, add a `PR_AUTHOR` env var (default empty) and skip G3 with a `SKIP` report when it is `dependabot[bot]`.
- [x] 1.2 In `scripts/check-gates.sh`, skip G4 with a `SKIP` report when `PR_AUTHOR` is `dependabot[bot]`, leaving the missing-body skip in place.
- [x] 1.3 In `.github/workflows/openspec.yml`, pass `PR_AUTHOR: ${{ github.event.pull_request.user.login }}` to the gate step.
- [x] 1.4 Mirror the gate script and workflow change into `cli/internal/seed/repo/` so `TestSeededGatesMatchRoot` stays green.

## 2. Tests

- [x] 2.1 In `scripts/check-gates-test.sh`, add a case that runs with `PR_AUTHOR=dependabot[bot]` and a `dependabot/...` branch and asserts G3 and G4 report `SKIP`.
- [x] 2.2 In `scripts/check-gates-test.sh`, assert a run without `PR_AUTHOR` still fails G3 and G4 on the same branch and body.

## 3. Verification

- [x] 3.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, `sh scripts/check-gates-test.sh`, and `sh scripts/check-gates.sh`.
