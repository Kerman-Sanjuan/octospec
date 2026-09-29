# AGENTS.md

octospec makes OpenSpec GitHub-native: **GitHub Issues are the backlog and the
human source of truth**, the repository holds the machine artifacts under
`openspec/`, and CI enforces the gates. This file is the entry point for anyone
— human or agent — working **on** octospec. If you are **using** octospec in
another repository, read `README.md` and the commands instead.

This repo **dogfoods** itself: the workflow in `cli/internal/payload/commands/`
is the workflow used on this repo.

## Where the explanations live

- `README.md` — what octospec is, the workflow loop, install, and the command surface.
- `openspec/specs/` — the living contracts, one `spec.md` per capability.
- `openspec/changes/archive/` — every shipped change, with its proposal, design,
  specs and tasks. Read these to see how a past decision was made.
- `docs/plans/` — the original (now historical) design and plan, from before the
  CLI. Treat `README.md` + `openspec/specs/` as current.

## Layout

- `cli/` — the Go CLI.
  - `cli/internal/payload/commands/` — **the single canonical source** of the commands.
  - `cli/internal/targets/` — maps each tool to its directory and front matter.
  - `cli/internal/seed/` — the embedded OpenSpec schema and repo seed.
  - `cli/internal/{install,update,config,plan}/` — install, in-place update, state.
- `scripts/check-gates.sh` — the gates (G2a, G2b, G3, G4, G6), run locally and in CI.
- `.github/workflows/` — `cli.yml` (Go checks), `openspec.yml` (gates), `release.yml`.
- `install.sh` — only the `curl | sh` installer for the CLI.

## Invariants

- **One canonical payload.** Commands live once in `cli/internal/payload/commands/`
  and are rendered per tool. Never add a per-tool copy in the repo.
- **Generated files are not committed.** `openspec update` output
  (`.pi/skills/`, `.opencode/skills/`, `.github/skills/`) is gitignored.
- **The schema and repo seed are embedded** under `cli/internal/seed/`.
- **Gates are enforced.** `main` requires the `cli` and `gates` checks; every
  change lands through a pull request.

## Build and check

```sh
cd cli && gofmt -l . && go vet ./... && go test ./...
openspec validate --all --strict
sh scripts/check-gates.sh
```

## Definition of done

A change is not done until:

- [ ] `gofmt -l` is clean; `go vet ./...` and `go test ./...` pass
- [ ] `openspec validate --all --strict` passes and the change is archived
- [ ] **Changed a command?** Update the README **Commands** table.
- [ ] **Added or removed a tool?** Update the README **Install** table and the target tests.
- [ ] **Superseded a file or script?** Delete it — no dead files.
- [ ] **Changed the repo structure?** Update this file and the README.
- [ ] **Changed behaviour?** Update the living spec in `openspec/specs/`.
- [ ] `cli` and `gates` are green on the pull request.

## Conventions

- Commits: `<type>(#<issue>): <summary>`.
- Branches: `feat|fix/<issue>-<slug>` (the G3 gate).
- One OpenSpec change per issue (the G8 advisory).
