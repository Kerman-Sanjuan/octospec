## 1. Real e2e in CI

- [x] 1.1 Add an OpenSpec install step (pinned version) to the `cli` job in `.github/workflows/cli.yml` before the e2e step.
- [x] 1.2 Make `scripts/e2e-test.sh` fail (non-zero) with a clear message when `openspec` is not on PATH, instead of printing a skip and exiting 0.

## 2. Truthful docs

- [x] 2.1 State the tested OpenSpec version in `README.md` and drop the "or newer" / ">= 1.3.1" claim.
- [x] 2.2 Repair the mutilated prose in the README skills paragraph.

## 3. Gate consistency

- [x] 3.1 Give G2b and G7 an else branch that reports a SKIP when `openspec` is missing, consistent with G2a.
- [x] 3.2 Fix the G6 failure message to name `skip_specs: true`.

## 4. Archived changes

- [x] 4.1 Check the open tasks in the archived changes gh-30, gh-31, and gh-48 so the archived changes have no unchecked boxes.

## 5. Cleanup

- [x] 5.1 Delete the merged remote branches 66, 85, 92, and 94.

## 6. Checks

- [x] 6.1 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.
- [x] 6.2 Confirm `scripts/e2e-test.sh` fails without openspec on PATH and passes with it.
