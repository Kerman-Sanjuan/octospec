# Tasks — review ledger

## 1. The record

- [ ] 1.1 `review.Record` with schema, pr, issue, head, score, threshold, risk, risk_why, verdict,
  findings, at.
- [ ] 1.2 `WriteRecord(dir, rec)` creates `<pr>.json`, rewriting it atomically.
- [ ] 1.3 It refuses to write without a pull request number: a record with no key is not a record.

## 2. The summary

- [ ] 2.1 `ReadAll(dir)` skips files it cannot parse and says how many it skipped.
- [ ] 2.2 `Summarise(records)` returns counts by risk, counts by verdict, and the mean score.
- [ ] 2.3 `octospec ledger <dir>` prints the table, the counts and the mean; empty is not an error.

## 3. The command surface

- [ ] 3.1 `review --ledger <dir>` writes the record next to the printed verdict.
- [ ] 3.2 `--pr` on `review`, so the record has a key without guessing from the branch.

## 4. The label

- [ ] 4.1 `seed` provisions `risk:high`.
- [ ] 4.2 `/review` adds it on high and removes it otherwise.
- [ ] 4.3 `docs/commands.md` says so.

## 5. Tests

- [ ] 5.1 Writing twice rewrites one file.
- [ ] 5.2 A missing pull request number fails.
- [ ] 5.3 Unreadable files are skipped and counted, not fatal.
- [ ] 5.4 The summary maths: counts, mean, empty set.
- [ ] 5.5 `go test ./...` and `scripts/e2e-test.sh` pass.
