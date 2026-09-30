# Agent definitions

One canonical agent per stage of the octospec loop. octospec renders these
files into each tool's agent directory, so the repository never stores a copy
per tool.

## Canonical fields

| Field | Meaning |
|---|---|
| `name` | The file name, without extension. |
| `stage` | One of `idea`, `spec`, `apply`, `ship`, `archive`. |
| `role` | The model role: `thinking`, `implementer`, or `reviewer`. |
| `description` | One line, shown by the tool. |
| `skills` | The OpenSpec skills the agent may load, comma-separated. |
| `writes` | The artifact the agent produces. |

The body is the mission. It references only the stage command and the skills
listed for that stage, so the rendered file shows exactly what the agent
receives.

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
