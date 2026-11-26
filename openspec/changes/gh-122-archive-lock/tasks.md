# Tasks — archive lock

## 1. Lock primitive

- [ ] 1.1 Add `cli/internal/lock/lock.go` with the lock record type, paths, and the `O_EXCL`
      acquire primitive
- [ ] 1.2 Add stale detection in `cli/internal/lock/lock.go`: pid + host for the local case,
      TTL for the remote case
- [ ] 1.3 Add `Release` guarded by change name and `ReleaseForce` for the explicit human path

## 2. Tests

- [ ] 2.1 Add `cli/internal/lock/lock_test.go` covering acquire, contention, waiting, release
      by another change, stale detection and force release

## 3. Command surface

- [ ] 3.1 Add the `lock` subcommand (`status`, `acquire`, `release`) in
      `cli/cmd/octospec/main.go`
- [ ] 3.2 Extend `cli/cmd/octospec/main_test.go` for the new subcommand's argument validation
- [ ] 3.3 Report the lock state from `octospec doctor` in `cli/internal/doctor/doctor.go`

## 4. Wire the stage

- [ ] 4.1 Add the acquire/park/release steps to `cli/internal/payload/commands/archive.md`
- [ ] 4.2 Add the diff-scope check to the same payload, after the fold and before the commit
- [ ] 4.3 Ensure `octospec seed` writes `.octospec/locks/` to `.gitignore`

## 5. Documentation

- [ ] 5.1 Document the lock, the reasons and the command in `docs/commands.md`
- [ ] 5.2 Document the stale-lock runbook in `docs/troubleshooting.md`
