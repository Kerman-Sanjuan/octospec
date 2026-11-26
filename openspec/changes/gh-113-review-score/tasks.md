# Tasks — review score

## 1. Scoring engine

- [ ] 1.1 Add `cli/internal/review/findings.go`: the finding type, the closed axis and type
      vocabularies, and validation
- [ ] 1.2 Add `cli/internal/review/policy.go`: load `.octospec/review.json`, merge with defaults
- [ ] 1.3 Add `cli/internal/review/score.go`: the arithmetic, the per-axis breakdown and the verdict
- [ ] 1.4 Add `cli/internal/review/risk.go`: the sensitive-path / size / familiarity table

## 2. Tests

- [ ] 2.1 Add `cli/internal/review/review_test.go`: determinism, clean change, weights, unknown
      axis, threshold boundary, each risk row, defaults without a policy file

## 3. Command surface

- [ ] 3.1 Add the `review` subcommand in `cli/cmd/octospec/main.go` with `--findings`, `--json`,
      `--enforce` and the exit-code mapping
- [ ] 3.2 Extend `cli/cmd/octospec/main_test.go` for the review exit codes

## 4. The stage

- [ ] 4.1 Add `cli/internal/payload/commands/review.md`: the five axes, the JSON shape, the
      one-comment rule and the agentic-process notice
- [ ] 4.2 Update the payload command count in `cli/internal/payload/embed_test.go`
- [ ] 4.3 Seed a default `.octospec/review.json` from `cli/internal/seed/repo/` without
      overwriting an existing one

## 5. Documentation

- [ ] 5.1 Add the command and the workflow stage to `docs/commands.md`
- [ ] 5.2 Document the policy keys and the calibration advice in `docs/commands.md`
