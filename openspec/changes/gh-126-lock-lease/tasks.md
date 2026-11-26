# Tasks — lock lease

## 1. The record and its lease

- [ ] 1.1 Add `expires_at` to `lock.Record`, with a default TTL constant.
- [ ] 1.2 `Acquire` takes a TTL and writes the expiry; `tryCreate` stays the only creation path.
- [ ] 1.3 Staleness is `now >= expires_at`, with `acquired_at + default TTL` as the fallback for
  records written before this change.

## 2. Heartbeat

- [ ] 2.1 `lock.Heartbeat(repo, change, ttl)` extends a valid lease and fails on an expired one.
- [ ] 2.2 The rewrite is atomic: temporary file plus rename.
- [ ] 2.3 A heartbeat naming another change fails, and says who holds it.

## 3. Command surface

- [ ] 3.1 `octospec lock heartbeat --change <change>` with `--ttl`.
- [ ] 3.2 `octospec lock acquire --ttl`, defaulting to 30m.
- [ ] 3.3 `lock status` prints holder, age and expiry; `--force` appears only when expired.
- [ ] 3.4 `doctor` reports the same, and stops calling a fresh lock stale.

## 4. The stage

- [ ] 4.1 `/archive` heartbeats after the fold and after the changelog entry.
- [ ] 4.2 The payload stops offering `--force` as a remedy for a lock it just took.

## 5. Tests

- [ ] 5.1 Regression: a dead PID inside a valid lease is not stale.
- [ ] 5.2 `acquire` then `status` in different processes reports held.
- [ ] 5.3 Expiry reports stale; heartbeat on an expired lease fails and leaves it stale.
- [ ] 5.4 Heartbeat extends; release is guarded by name; `--force` breaks a valid lease.
- [ ] 5.5 `scripts/e2e-test.sh` and `go test ./...` pass.

## 6. Documentation

- [ ] 6.1 `docs/commands.md`: the three lock verbs and the TTL.
- [ ] 6.2 `docs/troubleshooting.md`: an expired lease, and why `--force` is the human's call.
