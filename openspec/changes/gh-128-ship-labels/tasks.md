# Tasks — label the pull request

## 1. The stage

- [ ] 1.1 `/ship` reads the issue's labels and copies `type:*`.
- [ ] 1.2 It copies `mode:autonomous` when present.
- [ ] 1.3 A `type:*`-less issue stops the stage with a named reason.

## 2. The agent and the docs

- [ ] 2.1 `agents/ship.md`: the rule and why the kind belongs on the pull request.
- [ ] 2.2 `docs/commands.md`: the `/ship` row says what it labels.

## 3. Verification

- [ ] 3.1 Dogfood: `/ship` this change and check `gh pr view --json labels` lists both labels.
- [ ] 3.2 `scripts/e2e-test.sh` and `go test ./...` still pass.
