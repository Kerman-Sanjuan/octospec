# Tasks — the intake grill

## 1. The prompts

- [x] 1.1 Rewrite `cli/internal/payload/commands/idea.md` around the frontier: rounds, recommendation per
      question, facts looked up, no filing until the frontier is empty.
- [x] 1.2 Rewrite `cli/internal/payload/commands/bug.md` the same way, for the bug fields.
- [x] 1.3 Teach `cli/internal/payload/commands/spec.md` to read `## Grill` and not re-ask it.

## 2. The agent contract

- [x] 2.1 `cli/internal/payload/agents/idea.md`: the seven groups, facts versus decisions, the stop rule,
      and `needs_human` when nobody is there.

## 3. The gate

- [x] 3.1 `scripts/check-gates.sh`: G1 also requires `## Grill`, `SHALL`/`MUST` per requirement line, a
      filled `## Success criteria`, and an empty `## Open questions`.
- [x] 3.2 Mirror it in `cli/internal/seed/repo/scripts/check-gates.sh`, character for character.
- [x] 3.3 `scripts/check-gates-test.sh`: a requirement with no SHALL fails, an open question fails, and
      the complete fixture still passes.

## 4. The seed

- [x] 4.1 `cli/internal/seed/repo/.github/ISSUE_TEMPLATE/feature.yml` and `bug.yml` gain the `Grill`
      section and the SHALL/MUST wording.
- [x] 4.2 `cli/internal/seed/schema/templates/issue.md` mirrors the new section.

## 5. Everywhere else

- [x] 5.1 `README.md`: step 1 of the loop and the G1 row.
- [x] 5.2 `docs/commands.md`: the `/idea` and `/bug` rows.
- [x] 5.3 `docs/getting-started.md`: the first step, in one line.

## 6. Verification

- [x] 6.1 `cd cli && go test ./...` and `go vet ./...`.
- [x] 6.2 `sh scripts/check-gates.sh` and `sh scripts/check-gates-test.sh`.
- [x] 6.3 `sh scripts/e2e-test.sh`.
- [x] 6.4 Run the new gate against issue #134 itself: it must pass.
