## Context

The CLI already has one canonical payload: `cli/internal/payload/commands/*.md`, embedded in the binary and rendered per tool by `targets.Target`. `plan.Files` turns a tool list into concrete files, and `install`/`update` write them with a content hash per managed file in `.octospec/octospec.json`. The `openspec` CLI owns the OpenSpec skills.

Two things are missing. There is no definition of an agent per stage, and there is no place to pick a model. The loop is driven by commands that carry the whole instruction set, so an agent loaded for one stage sees the whole manual. The repo supports pi, opencode, GitHub Copilot, and Claude Code, and every tool gets its own rendered copy from the one payload.

## Goals / Non-Goals

**Goals:**
- One canonical definition per stage agent (idea, spec, apply, ship, archive).
- One model configuration, role-based, editable through an interactive TUI.
- Render every agent for every supported tool from the single source.
- Add a new tool by adding a mapping only.
- Leave the gate rules untouched.

**Non-Goals:**
- Local or unsupported model providers.
- A rewrite of the existing OpenSpec skills.
- Migration of already-installed consumer repos beyond what `octospec update` does automatically.
- Any change to `scripts/check-gates.sh`.

## Decisions

- **The agent definitions live beside the commands in the canonical payload.** A new `cli/internal/payload/agents/*.md` holds them, with front matter (name, stage, role, description, skills, writes) and a body. The same embed and render pipeline covers them. *Alternative rejected:* a top-level `agents/` directory copied per tool, which reintroduces the per-tool copy the payload exists to remove.
- **The model configuration lives in `.octospec/octospec.json`.** Add a `models` map of role to model id. Changing a model means changing one file, and `update` re-stamps every agent. *Alternative rejected:* a separate `agents.yaml`, which adds a second state file and a second source of truth.
- **Roles, not stages.** The roles are `thinking`, `implementer`, and `reviewer`. The stages map as idea and spec to thinking, apply to implementer, and ship and archive to reviewer. *Alternative rejected:* one role per stage, which triples the configuration for little gain.
- **An empty model means the tool default.** Install works before the user picks anything, and a rendered agent omits the model field when its role has no value. *Alternative rejected:* requiring a model at install, which blocks a first run.
- **The multi-tool standard is a small set of canonical fields plus a per-tool render.** Extend `targets.Target` with an agent directory and an agent format, the way it already carries a command directory and format. A new tool is a new mapping. *Alternative rejected:* a generic file for every tool, which loses each tool's native agent mechanism.
- **Context injection is static and auditable.** Each agent references only its stage command and skills, so the rendered file shows exactly what the agent receives. *Alternative rejected:* a runtime prompt builder, which hides the effective prompt from inspection.
- **The model TUI uses a small interactive library.** Use `github.com/charmbracelet/huh` for the role-to-model form, with a flag path for non-interactive use. *Alternative rejected:* hand-rolled `bufio` prompts, which are not a TUI, and a full `bubbletea` program, which is heavier than a form needs.
- **The developing repo adopts the standard.** `AGENTS.md` and the workflow commands of this repo point at the canonical agents and the model config. *Alternative rejected:* leaving the repo on the old rules while consumers get the new ones, which breaks dogfooding.

## Risks / Trade-offs

- [Tool agent formats differ and change over time] -> Keep the canonical fields small; keep each format quirk inside a target mapping; cover every renderer with a test.
- [A TUI library grows the binary and its dependency tree] -> One small library, and keep the non-interactive path flag-driven so CI never needs a terminal.
- [A model id a tool does not know breaks a rendered agent] -> Empty means the tool default; document the accepted values; `doctor` reports an unknown role.
- [The change is large] -> Land it in groups: standard and definitions, then config and TUI, then docs and self-hosting.

## Migration Plan

`octospec update` re-renders the agents for an existing install and keeps local edits through the managed-file hashes. An existing `.octospec/octospec.json` without a `models` key reads as empty, so no model is stamped and behavior is unchanged until the user picks one. Rollback: revert the CLI version and run `update` again, which rewrites the managed files from the previous payload.

## Open Questions

- Does Codex expose a native agent format, or does it need a single-file render like a prompt set?
- Which default model ids does the TUI offer per tool, and who maintains that list?
