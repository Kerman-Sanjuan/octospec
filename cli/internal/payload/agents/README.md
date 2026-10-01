# Agent definitions

One canonical agent per stage of the octospec loop. octospec renders these
files into each tool's agent directory, so the repository never stores a copy
per tool.

## How a stage runs

An agent definition is the stage's contract: its scope, its skills, its tool
surface, and its model role. The stage command runs the stage in the current
session under that contract. octospec does not dispatch the stage to a subagent
yet; some harnesses have no subagent mechanism at all, and pi is one of them.
Real dispatch is a follow-up.

## Canonical fields

| Field | Meaning |
|---|---|
| `name` | The file name, without extension. |
| `stage` | One of `idea`, `spec`, `apply`, `ship`, `archive`. |
| `role` | The model role: `thinking`, `implementer`, or `reviewer`. |
| `description` | One line, shown by the tool. |
| `skills` | The OpenSpec skills the agent may load, comma-separated. |
| `tools` | The capability surface the agent may use, comma-separated. |
| `writes` | The artifact the agent produces. |

The body is the mission. It references only the stage command and the skills
listed for that stage, so the rendered file shows exactly what the agent
receives. It also names no tool outside the `tools` surface.

## Tool surface

The `tools` field is a capability-group allowlist (B15 TOOL SUBSET). The
canonical capabilities are `read`, `edit`, `search`, `shell`, `web`, and
`agent`. The deployer maps them to each harness's native permission syntax, so
the canonical body stays portable.

| Stage | Role | tools |
|---|---|---|
| idea | thinking | read, search, shell |
| spec | thinking | read, search, edit, shell |
| apply | implementer | read, search, edit, shell |
| ship | reviewer | read, search, shell |
| archive | reviewer | read, search, edit, shell |

## Portability

Each harness realizes the tool surface to a different degree. A harness that
cannot enforce it is declared, and the render does not emit a field the harness
ignores.

| Harness | Agent file | Tool surface |
|---|---|---|
| claude | `.claude/agents/<name>.md` | enforced, `tools` |
| opencode | `~/.config/opencode/agents/<name>.md` | enforced, `permissions` (provisional shape) |
| copilot | `.github/agents/<name>.agent.md` | enforced, `tools` |
| pi | `~/.pi/agent/agents/<name>.md` | advisory, in the body only |
| codex | none | unsupported, Codex has no persona file |

## Model roles

The model is not part of the definition. octospec stamps it at render time
from the one model configuration, keyed by role. An empty model means the tool
default, so a rendered agent omits the `model` field.

| Role | Stages |
|---|---|
| `thinking` | idea, spec |
| `implementer` | apply |
| `reviewer` | ship, archive |

## Render rules

Each target maps an agent directory and an agent format.

- `default`: YAML front matter with `name`, `description`, and `model` when
  set.
- `copilot`: YAML front matter with `mode: agent`, `description`, and `model`
  when set.

Adding a tool means adding a target mapping. The agent definitions do not
change.
