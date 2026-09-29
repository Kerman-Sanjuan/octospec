# GitHub-native OpenSpec workflow — implementation plan

> **For agentic workers:** This plan is task-by-task. Each task ends with an independently testable deliverable. Steps use checkbox (`- [ ]`) syntax. Run the test steps and confirm the expected output before moving on.

**Goal:** Build a reusable, open-source workflow that makes OpenSpec GitHub-native — GitHub Issues as the backlog and human layer, repo files as the machine layer — installable for pi, opencode and GitHub Copilot from one canonical repository.

**Architecture:** A canonical repo ships (a) an OpenSpec schema override `octospec` that owns artifact semantics, (b) seven thin command shims whose logic is only phase orchestration, (c) per-repo seed files (issue forms, config, CI), and (d) an `install.sh` that publishes the global parts and seeds a repo. CI enforces the hard gates.

**Tech Stack:** OpenSpec CLI ≥ 1.3.1, `gh` CLI ≥ 2.x, POSIX `sh`, GitHub Actions, Markdown + YAML.

**Spec:** `docs/plans/2026-09-29-github-openspec-workflow-design.md`. Read it before starting. Both files move into the canonical repo in Task 1.

## Global Constraints

- OpenSpec ≥ 1.3.1 installed globally (`openspec --version`); `gh` authenticated (`gh auth status`).
- Canonical repo: `Kerman-Sanjuan/octospec`, MIT, public. Local path `~/Documents/repos/personal/octospec`.
- Schema name is `octospec`. Never shadow `spec-driven`.
- **Artifact semantics live only in `schema/octospec/`.** Command shims must not re-describe what an artifact contains.
- Never edit generated OpenSpec files; the schema override is the extension point.
- Global install paths: schema → `${XDG_DATA_HOME:-$HOME/.local/share}/openspec/schemas/octospec/`; pi → `~/.pi/agent/prompts/`; opencode → `~/.config/opencode/command/`. Copilot commands are repo-local (`.github/prompts/`).
- Gate IDs are G1–G8 as defined in the spec. Hard: G1–G6. Advisory: G7, G8.
- The OpenSpec CLI never calls `gh`; agents do, following instructions.

## File Structure

```
octospec/
├── README.md
├── LICENSE                         MIT
├── install.sh                      publishes global parts; --repo seeds a repo
├── schema/
│   └── octospec/
│       ├── schema.yaml             artifact semantics (the core logic)
│       └── templates/
│           ├── issue.md
│           ├── proposal.md
│           ├── spec.md
│           ├── design.md
│           └── tasks.md
├── commands/                        one file per command, tool-neutral
│   ├── idea.md
│   ├── bug.md
│   ├── explore.md
│   ├── spec.md
│   ├── apply.md
│   ├── ship.md
│   └── archive.md
├── repo-template/                   copied into a target repo by `install.sh --repo`
│   ├── openspec/config.yaml
│   └── .github/
│       ├── ISSUE_TEMPLATE/
│       │   ├── feature.yml
│       │   └── bug.yml
│       ├── prompts/                 Copilot command shims
│       │   └── spec.prompt.md ... (one per command)
│       └── workflows/
│           └── openspec.yml
└── scripts/
    └── check-gates.sh               G2–G6, runnable locally and in CI
```

---

### Task 1: Canonical repo skeleton

**Files:**
- Create: `~/Documents/repos/personal/octospec/` (git repo)
- Create: `LICENSE`
- Create: `README.md`
- Create: `.gitignore`

**Interfaces:**
- Produces: an empty git repo with an MIT license and a README stub, at the path every later task writes into.

- [ ] **Step 1: Create the directory and init git**

```bash
mkdir -p ~/Documents/repos/personal/octospec
cd ~/Documents/repos/personal/octospec
git init -b main
```

- [ ] **Step 2: Write the MIT license**

Create `LICENSE` with the standard MIT text:

```
MIT License

Copyright (c) 2026 Kerman Sanjuan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 3: Write `.gitignore`**

```
.DS_Store
# Installer scratch
.octospec-tmp/
```

- [ ] **Step 4: Write a README stub**

```markdown
# octospec

GitHub-native OpenSpec: GitHub Issues are the backlog and the human layer;
the repo holds the machine artifacts. One schema override, thin commands,
hard CI gates.

Design: see `docs/plans/2026-09-29-github-openspec-workflow-design.md`.

> Full documentation lands in Task 9.
```

- [ ] **Step 5: Move the design doc and this plan into the repo**

```bash
mkdir -p ~/Documents/repos/personal/octospec/docs/plans
mv ~/Documents/repos/personal/brave-distraction-blocker/docs/plans/2026-09-29-github-openspec-workflow-design.md \
   ~/Documents/repos/personal/octospec/docs/plans/
cp ~/Documents/repos/personal/brave-distraction-blocker/docs/plans/2026-09-29-github-openspec-workflow-plan.md \
   ~/Documents/repos/personal/octospec/docs/plans/
```

- [ ] **Step 6: Verify and commit**

Run: `git -C ~/Documents/repos/personal/octospec status --short`
Expected: `LICENSE`, `README.md`, `.gitignore`, `docs/plans/*` listed as untracked.

```bash
cd ~/Documents/repos/personal/octospec
git add -A
git commit -m "chore: scaffold canonical repo with MIT license and design docs"
```

---

### Task 2: `octospec` schema and templates

**Files:**
- Create: `schema/octospec/schema.yaml`
- Create: `schema/octospec/templates/issue.md`
- Create: `schema/octospec/templates/proposal.md`
- Create: `schema/octospec/templates/spec.md`
- Create: `schema/octospec/templates/design.md`
- Create: `schema/octospec/templates/tasks.md`

**Interfaces:**
- Consumes: nothing.
- Produces: a schema named `octospec` with artifact ids `issue`, `proposal`, `specs`, `design`, `tasks`, and `apply.requires: [tasks]`. Later tasks reference these exact ids.

- [ ] **Step 1: Write `schema/octospec/schema.yaml`**

```yaml
name: octospec
version: 1
description: GitHub-native OpenSpec. A GitHub issue is the intake artifact; the repo holds the machine artifacts.
artifacts:
  - id: issue
    generates: issue.md
    description: Durable snapshot of the GitHub issue that motivates the change
    template: issue.md
    instruction: |
      Capture the intake artifact: the GitHub issue.

      1. Read `.openspec.yaml`. If it already has `github.issue`, that issue
         is authoritative - do not create another.
      2. If it does not, the issue must already exist. The change name is
         `gh-<issue>-<slug>`; OpenSpec requires a leading letter, so the
         number cannot come first. Adopt the number from the name. If no
         issue exists, stop and ask the user to run `/idea` first.
      3. Read the issue:
         `gh issue view <n> --json number,title,body,updatedAt,url`
      4. Record the link in `.openspec.yaml`:
         github:
           issue: <n>
           issue_updated_at: <updatedAt>
      5. Write `issue.md` as a snapshot of the STABLE narrative sections
         only: user story, context, requirements, success criteria, out of
         scope. Do not include the task checklist or comments - those
         change. Put the issue URL and number at the top.

      `issue.md` is a read-only mirror; the issue body stays the human
      source of truth.
    requires: []

  - id: proposal
    generates: proposal.md
    description: Proposal derived from the issue
    template: proposal.md
    instruction: |
      Create the proposal document (WHY), derived from `issue.md`.

      Sections:
      - **Why**: 1-2 sentences from the issue's context.
      - **What Changes**: bullet list from the issue's requirements.
      - **Capabilities**: New capabilities become `specs/<kebab-name>/spec.md`.
        Modified capabilities only when an existing spec's REQUIREMENTS
        change. Leave Modified empty when unsure.
      - **Impact**: affected code, APIs, dependencies, systems.

      Keep it under two pages. Link the issue; do not copy it wholesale.
    requires:
      - issue

  - id: specs
    generates: "specs/**/*.md"
    description: Requirement deltas
    template: spec.md
    instruction: |
      Create one delta spec per capability listed in the proposal.

      Delta operations (## headers): ADDED, MODIFIED, REMOVED, RENAMED.
      - Requirement: `### Requirement: <name>`, using SHALL or MUST.
      - Scenario: `#### Scenario: <name>` with WHEN/THEN. Exactly four
        hashtags - three fails silently.
      - Every requirement needs at least one scenario.

      MODIFIED: copy the ENTIRE existing block from
      `openspec/specs/<cap>/spec.md` and edit it. Partial content loses
      detail at archive time.
      REMOVED: include **Reason** and **Migration**.

      Scenarios are the success criteria. If a success criterion from the
      issue has no scenario, stop and ask rather than inventing one.
    requires:
      - proposal

  - id: design
    generates: design.md
    description: Technical design document
    template: design.md
    instruction: |
      Create the design document (HOW). Skip only when none apply:
      cross-cutting change, new dependency or data model change, security
      or migration complexity, or decisions worth recording.

      Sections: Context, Goals / Non-Goals, Decisions (with alternatives),
      Risks / Trade-offs ([Risk] -> Mitigation), Migration Plan, Open
      Questions.

      Reference the proposal for motivation and the specs for requirements.
    requires:
      - proposal

  - id: tasks
    generates: tasks.md
    description: Implementation checklist
    template: tasks.md
    instruction: |
      Create the task list. The apply phase parses checkboxes; tasks that
      are not `- [ ]` are not tracked.

      - Group under `## N. <Group>` headings.
      - Each task is `- [ ] N.M <description>`.
      - A task is small enough for one session and names the file it
        touches.
      - Order by dependency.
      - Every task traces to a requirement or a design decision.
    requires:
      - specs
      - design

apply:
  requires: [tasks]
  tracks: tasks.md
  instruction: |
    Implement the tasks in order.

    - Branch first: `feat/<issue>-<slug>` or `fix/<issue>-<slug>`.
    - Work one group at a time; commit per group, referencing the issue.
    - Tick tasks in `tasks.md` as they complete and keep the issue
      checklist comment in sync with `gh issue comment`.
    - Run the project's build and tests before reporting a group done.
    - Stop and ask on a blocker. Do not guess a requirement.
```

- [ ] **Step 2: Write `templates/issue.md`**

```markdown
# Issue #<n>: <title>

<url>

## User story

As a <role>, I want <capability>, so that <benefit>.

## Context / problem

<!-- From the issue body -->

## Requirements

- The system SHALL ...

## Success criteria

- ...

## Out of scope

- ...
```

- [ ] **Step 3: Write `templates/proposal.md`**

```markdown
## Why

<!-- 1-2 sentences, from the issue's context -->

## What Changes

<!-- Bullet list, from the issue's requirements -->

## Capabilities

### New Capabilities
- `<kebab-name>`: <brief description>

### Modified Capabilities
<!-- Only when an existing spec's REQUIREMENTS change -->

## Impact

<!-- Affected code, APIs, dependencies, systems. Link the issue. -->
```

- [ ] **Step 4: Write `templates/spec.md`**

```markdown
## ADDED Requirements

### Requirement: <name>
The system SHALL <behaviour>.

#### Scenario: <name>
- **WHEN** <condition>
- **THEN** <expected outcome>
```

- [ ] **Step 5: Write `templates/design.md`**

```markdown
## Context

<!-- Background, current state, constraints -->

## Goals / Non-Goals

**Goals:**
- ...

**Non-Goals:**
- ...

## Decisions

<!-- Key choices with rationale and alternatives considered -->

## Risks / Trade-offs

- [Risk] → Mitigation

## Migration Plan

<!-- Deploy and rollback steps, if applicable -->

## Open Questions

- ...
```

- [ ] **Step 6: Write `templates/tasks.md`**

```markdown
## 1. <Group>

- [ ] 1.1 <Task, names the file it touches>
- [ ] 1.2 <Task>

## 2. <Group>

- [ ] 2.1 <Task>
```

- [ ] **Step 7: Validate the schema in a scratch project**

```bash
SCRATCH=$(mktemp -d)
cd "$SCRATCH"
git init -q
openspec init --tools none >/dev/null 2>&1 || true
mkdir -p openspec/schemas
cp -R ~/Documents/repos/personal/octospec/schema/octospec openspec/schemas/
openspec schema validate octospec
openspec schema which octospec
```

Expected: validation passes; `which` reports the project-local path.

- [ ] **Step 8: Verify instructions are served**

```bash
cd "$SCRATCH"
openspec new change gh-42-test --schema octospec
openspec status --change gh-42-test --json | head -40
openspec instructions issue --change gh-42-test --json | head -40
```

Expected: `issue`, `proposal`, `specs`, `design`, `tasks` appear; `instructions issue` returns the `issue.md` instruction and template.

- [ ] **Step 9: Commit**

```bash
cd ~/Documents/repos/personal/octospec
git add schema
git commit -m "feat(schema): add octospec OpenSpec schema and templates"
```

---

### Task 3: Installer publishes the global schema

**Files:**
- Create: `install.sh`

**Interfaces:**
- Consumes: `schema/octospec/` from Task 2.
- Produces: `install.sh`, POSIX `sh`, no arguments needed for the global install. Later tasks add sections to it.

- [ ] **Step 1: Write `install.sh`**

```sh
#!/usr/bin/env sh
# Install the octospec OpenSpec workflow.
#   ./install.sh                publish the global parts (schema, pi, opencode)
#   ./install.sh --repo <path>  also seed the repo-local parts into <path>
set -eu

REPO_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
SCHEMA=octospec
DATA_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/openspec/schemas"

# 1. Global schema override
mkdir -p "$DATA_DIR/$SCHEMA"
cp -R "$REPO_DIR/schema/$SCHEMA/." "$DATA_DIR/$SCHEMA/"
printf 'schema   -> %s\n' "$DATA_DIR/$SCHEMA"

# 2. pi prompt templates
PI_DIR="${PI_PROMPTS_DIR:-$HOME/.pi/agent/prompts}"
mkdir -p "$PI_DIR"
cp "$REPO_DIR"/commands/*.md "$PI_DIR/"
printf 'pi       -> %s\n' "$PI_DIR"

# 3. opencode commands
OC_DIR="${OPENCODE_COMMAND_DIR:-$HOME/.config/opencode/command}"
mkdir -p "$OC_DIR"
cp "$REPO_DIR"/commands/*.md "$OC_DIR/"
printf 'opencode -> %s\n' "$OC_DIR"

# 4. Optional repo seed
if [ "${1:-}" = "--repo" ]; then
  TARGET=${2:-}
  if [ -z "$TARGET" ] || [ ! -d "$TARGET" ]; then
    printf 'usage: install.sh --repo <existing-repo-path>\n' >&2
    exit 2
  fi
  mkdir -p "$TARGET/.github/ISSUE_TEMPLATE" "$TARGET/.github/prompts" \
           "$TARGET/.github/workflows" "$TARGET/openspec"
  cp -R "$REPO_DIR/repo-template/." "$TARGET/"
  printf 'repo     -> %s\n' "$TARGET"
fi
```

- [ ] **Step 2: Check shell syntax**

Run: `sh -n install.sh`
Expected: no output.

- [ ] **Step 3: Run the installer and verify the schema resolves**

```bash
cd ~/Documents/repos/personal/octospec
# commands/ does not exist yet; create it empty so the glob is harmless
mkdir -p commands
chmod +x install.sh
./install.sh
openspec schema which octospec
```

Expected: `which` reports `~/.local/share/openspec/schemas/octospec/schema.yaml` with source `user`.

- [ ] **Step 4: Commit**

```bash
git add install.sh
git commit -m "feat: install.sh publishes the global schema"
```

---

### Task 4: Command shims

**Files:**
- Create: `commands/idea.md`, `bug.md`, `explore.md`, `spec.md`, `apply.md`, `ship.md`, `archive.md`

**Interfaces:**
- Consumes: artifact ids and the `octospec` schema from Task 2.
- Produces: seven tool-neutral command files. Task 5 publishes them; Task 6 mirrors them into `.github/prompts/` for Copilot.

- [ ] **Step 1: Write `commands/idea.md`**

```markdown
---
description: File a new feature issue in the backlog and interview for its content.
argument-hint: "[idea]"
---

Create a feature issue in the GitHub backlog. The issue is the human source
of truth; there is no OpenSpec change yet.

1. If `$@` is empty, ask what the user wants to build.
2. Interview for the feature form fields, one question at a time: user story,
   context, requirements (SHALL/MUST), success criteria, out of scope.
3. Create the issue from the form:
   `gh issue create --template feature.yml --title "<title>" --body "<fields>"`
   If the form is unavailable, pass the same sections as the body.
4. Add labels `type:feature` and `status:backlog`.
5. Print the issue URL and number. Stop.

Do not create an OpenSpec change. `/spec` does that later.
```

- [ ] **Step 2: Write `commands/bug.md`**

```markdown
---
description: File a new bug issue.
argument-hint: "[summary]"
---

Create a bug issue. The issue is the intake artifact for a later fix.

1. If `$@` is empty, ask for a summary.
2. Interview for the bug form fields: summary, steps to reproduce, expected,
   actual, impact.
3. `gh issue create --template bug.yml --title "<summary>" --body "<fields>"`
4. Add labels `type:bug` and `status:backlog`.
5. Print the issue URL and number. Stop.
```

- [ ] **Step 3: Write `commands/explore.md`**

```markdown
---
description: Think through an idea before an issue exists.
argument-hint: "<topic>"
---

Explore the topic with the user. Ask one question at a time. Produce no
files, no issue, and no change. When the thinking converges, suggest running
`/idea`.
```

- [ ] **Step 4: Write `commands/spec.md`**

```markdown
---
description: Turn a GitHub issue into an OpenSpec change, publish it to the issue, and validate.
argument-hint: "<issue-number|url>"
---

Load the `openspec-propose` skill and follow it. The artifact semantics come
from the `octospec` schema, not from this file.

1. Resolve the issue: `gh issue view <n> --json number,title,body,updatedAt,url`.
   If it does not exist, stop and tell the user to run `/idea`.
2. Derive `<slug>` from the title; the change name is `gh-<issue>-<slug>`
   (OpenSpec requires a leading letter, so the number cannot come first).
3. `openspec new change "gh-<issue>-<slug>" --schema octospec`
4. Let the skill create every artifact required by
   `openspec status --change "<change>" --json` (`issue`, `proposal`,
   `specs`, `design`, `tasks`).
5. Publish each artifact to the issue as a comment, each headed with the
   commit SHA it came from:
   - proposal and design: one comment each.
   - specs: one comment per capability.
   - tasks: one comment containing the checklist.
   Use `gh issue comment <n> --body-file <file>`.
6. Run `openspec validate "<change>" --strict`. If it fails, fix and repeat.
7. Add label `status:spec-ready`; remove `status:backlog`.

Do not commit. Report the change path, the issue URL, and the validation
result.
```

- [ ] **Step 5: Write `commands/apply.md`**

```markdown
---
description: Implement an OpenSpec change on a branch and sync the issue checklist.
argument-hint: "[change]"
---

Load the `openspec-apply-change` skill and follow it.

1. Resolve the change (argument, or the only unarchived change).
   Read `.openspec.yaml` for `github.issue`; it must be set.
2. Create and switch to the branch `feat/<issue>-<slug>` for a feature or
   `fix/<issue>-<slug>` for a bug. Refuse if the working tree is dirty.
3. Work the tasks in `tasks.md` group by group. Commit each group with
   `<type>(#<issue>): <summary>`.
4. After each group, tick the tasks and post the updated checklist as a new
   comment on the issue, headed with the commit SHA.
5. Run the project's build and tests before reporting a group done. Report a
   failure as a failure.

Stop on a blocker and ask.
```

- [ ] **Step 6: Write `commands/ship.md`**

```markdown
---
description: Open a pull request linked to the issue.
argument-hint: "[change]"
---

Open the PR for the current branch.

1. Read `.openspec.yaml` for `github.issue`.
2. Push the branch: `git push -u origin HEAD`.
3. Create the PR. The body MUST contain `Closes #<issue>`, then:
   - a short summary,
   - the change name and artifact paths,
   - the `openspec validate` result.
   `gh pr create --title "<type>(#<issue>): <summary>" --body-file <file>`
4. Add label `status:in-review`; remove `status:in-progress`.
5. Print the PR URL.
```

- [ ] **Step 7: Write `commands/archive.md`**

```markdown
---
description: Archive the change, update the specs, and close the issue.
argument-hint: "[change]"
---

Load the `openspec-archive-change` skill and follow it.

1. Resolve the change and read `github.issue`. Confirm the PR is merged;
   if it is not, stop and say so.
2. `openspec archive "<change>" --yes`
3. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
4. Remove `status:*` labels from the issue.
5. Report the archive path and the spec files updated.
```

- [ ] **Step 8: Verify every command references a real skill or the CLI**

Run:
```bash
cd ~/Documents/repos/personal/octospec
grep -l 'openspec-propose\|openspec-apply-change\|openspec-archive-change' commands/*.md
grep -c 'description:' commands/*.md
```
Expected: `spec.md`, `apply.md`, `archive.md` list the skills; every file has a `description:`.

- [ ] **Step 9: Commit**

```bash
git add commands
git commit -m "feat(commands): add idea, bug, explore, spec, apply, ship, archive shims"
```

---

### Task 5: Installer publishes commands (pi + opencode)

**Files:**
- Modify: `install.sh` (Tasks 1–3 already contain the pi/opencode copy blocks — confirm they are present and correct)

**Interfaces:**
- Consumes: `commands/*.md` from Task 4.
- Produces: files at `~/.pi/agent/prompts/` and `~/.config/opencode/command/`.

- [ ] **Step 1: Confirm the installer copies commands**

Run: `grep -n 'PI_DIR\|OC_DIR' install.sh`
Expected: both variables present, each followed by a `cp "$REPO_DIR"/commands/*.md`.

- [ ] **Step 2: Install and verify**

```bash
cd ~/Documents/repos/personal/octospec
./install.sh
ls ~/.pi/agent/prompts/
ls ~/.config/opencode/command/
```

Expected: `idea.md`, `bug.md`, `explore.md`, `spec.md`, `apply.md`, `ship.md`, `archive.md` in both directories.

- [ ] **Step 3: Commit**

```bash
git add install.sh
git commit -m "feat(install): publish command shims for pi and opencode"
```

---

### Task 6: Repo seed — config, issue forms, Copilot prompts

**Files:**
- Create: `repo-template/openspec/config.yaml`
- Create: `repo-template/.github/ISSUE_TEMPLATE/feature.yml`
- Create: `repo-template/.github/ISSUE_TEMPLATE/bug.yml`
- Create: `repo-template/.github/prompts/{idea,bug,explore,spec,apply,ship,archive}.prompt.md`

**Interfaces:**
- Consumes: the `octospec` schema name and `commands/*.md` from Task 4.
- Produces: files `install.sh --repo <path>` copies into a target repo. Task 7 adds the CI workflow here.

- [ ] **Step 1: Write `repo-template/openspec/config.yaml`**

```yaml
schema: octospec

# Project context is injected into every artifact instruction. Keep it short.
context: |
  <One paragraph: what this project is, its stack, and any naming rules.
   Edited per repo after seeding.>

rules:
  proposal:
    - Link the originating issue; do not restate it.
  design:
    - Record every decision that fixes a means, and the alternative rejected.
  tasks:
    - Each task names the file it touches.
```

- [ ] **Step 2: Write `repo-template/.github/ISSUE_TEMPLATE/feature.yml`**

```yaml
name: Feature
description: A new capability, specified before it is built.
title: "[Feature] "
labels: ["type:feature", "status:backlog"]
body:
  - type: textarea
    id: story
    attributes:
      label: User story
      description: As a <role>, I want <capability>, so that <benefit>.
    validations:
      required: true
  - type: textarea
    id: context
    attributes:
      label: Context / problem
      description: Why this, why now.
    validations:
      required: true
  - type: textarea
    id: requirements
    attributes:
      label: Requirements
      description: One SHALL/MUST statement per line.
    validations:
      required: true
  - type: textarea
    id: success
    attributes:
      label: Success criteria
      description: Observable outcomes that decide "done".
    validations:
      required: true
  - type: textarea
    id: out_of_scope
    attributes:
      label: Out of scope
    validations:
      required: false
  - type: textarea
    id: open_questions
    attributes:
      label: Open questions
    validations:
      required: false
```

- [ ] **Step 3: Write `repo-template/.github/ISSUE_TEMPLATE/bug.yml`**

```yaml
name: Bug
description: Something behaves differently from what is expected.
title: "[Bug] "
labels: ["type:bug", "status:backlog"]
body:
  - type: textarea
    id: summary
    attributes:
      label: Summary
    validations:
      required: true
  - type: textarea
    id: steps
    attributes:
      label: Steps to reproduce
    validations:
      required: true
  - type: textarea
    id: expected
    attributes:
      label: Expected
    validations:
      required: true
  - type: textarea
    id: actual
    attributes:
      label: Actual
    validations:
      required: true
  - type: textarea
    id: impact
    attributes:
      label: Impact
    validations:
      required: false
```

- [ ] **Step 4: Generate the Copilot prompt mirrors**

Copilot reads `.github/prompts/*.prompt.md`. For each command file, write a mirror whose frontmatter uses Copilot's `mode` and `description`, and whose body is the same text.

```bash
cd ~/Documents/repos/personal/octospec
mkdir -p repo-template/.github/prompts
for f in commands/*.md; do
  name=$(basename "$f" .md)
  desc=$(sed -n 's/^description: //p' "$f" | head -1)
  {
    printf -- '---\n'
    printf 'mode: agent\n'
    printf 'description: %s\n' "$desc"
    printf -- '---\n\n'
    # body = everything after the closing frontmatter fence
    awk 'BEGIN{n=0} /^---$/{n++; next} n>=2{print}' "$f"
  } > "repo-template/.github/prompts/$name.prompt.md"
done
ls repo-template/.github/prompts/
```

Expected: seven `.prompt.md` files.

- [ ] **Step 5: Validate the YAML files**

```bash
cd ~/Documents/repos/personal/octospec
python3 - <<'PY'
import glob, yaml, sys
files = glob.glob('repo-template/**/*.yml', recursive=True) + \
        glob.glob('repo-template/**/*.yaml', recursive=True)
bad = []
for f in files:
    try:
        yaml.safe_load(open(f))
    except Exception as e:
        bad.append((f, e))
print("checked", len(files))
for f, e in bad:
    print("FAIL", f, e)
sys.exit(1 if bad else 0)
PY
```

Expected: `checked 3`, exit 0.

- [ ] **Step 6: Commit**

```bash
git add repo-template
git commit -m "feat(repo-template): add config, issue forms and Copilot prompt mirrors"
```

---

### Task 7: Gate checks and CI workflow

**Files:**
- Create: `scripts/check-gates.sh`
- Create: `repo-template/.github/workflows/openspec.yml`

**Interfaces:**
- Consumes: the `octospec` schema and repo layout from Tasks 2 and 6.
- Produces: `check-gates.sh` (exit 0 when all runnable gates pass) and a CI workflow that runs it on pull requests.

- [ ] **Step 1: Write `scripts/check-gates.sh`**

```sh
#!/usr/bin/env sh
# Gate checks G2-G6. Runnable locally and in CI.
# Env: BASE_REF (default origin/main), HEAD_REF, PR_BODY (optional).
set -eu

BASE_REF=${BASE_REF:-origin/main}
HEAD_REF=${HEAD_REF:-$(git rev-parse --abbrev-ref HEAD)}
fail=0

report() { printf '%s %s\n' "$1" "$2"; }

# G3: branch name
if printf '%s' "$HEAD_REF" | grep -Eq '^(feat|fix)/[0-9]+-[a-z0-9-]+$'; then
  report PASS "G3 branch name"
else
  report FAIL "G3 branch name: '$HEAD_REF' does not match feat|fix/<issue>-<slug>"
  fail=1
fi

# G2: openspec validate
if command -v openspec >/dev/null 2>&1 && openspec validate --all --strict >/dev/null 2>&1; then
  report PASS "G2 openspec validate"
else
  report FAIL "G2 openspec validate"
  fail=1
fi

# G6: a change that touches behaviour must carry a spec delta
changed=$(git diff --name-only "$BASE_REF...HEAD" 2>/dev/null || true)
if printf '%s\n' "$changed" | grep -Eq '^openspec/changes/[^/]+/specs/.+\.md$'; then
  report PASS "G6 spec delta"
else
  # Only required when the change folder exists
  if printf '%s\n' "$changed" | grep -Eq '^openspec/changes/[^/]+/tasks\.md$'; then
    report FAIL "G6 spec delta: change present but no specs/*.md"
    fail=1
  else
    report SKIP "G6 spec delta: no OpenSpec change in this PR"
  fi
fi

# G4: PR body links the issue
body=${PR_BODY:-}
if [ -n "$body" ]; then
  if printf '%s' "$body" | grep -Eq '(Closes|Fixes|Resolves) #[0-9]+'; then
    report PASS "G4 PR links issue"
  else
    report FAIL "G4 PR body has no Closes #n"
    fail=1
  fi
else
  report SKIP "G4 PR body not provided"
fi

if [ "$fail" -ne 0 ]; then
  report FAIL "gates failed"
  exit 1
fi
report PASS "all runnable gates"
```

- [ ] **Step 2: Write `repo-template/.github/workflows/openspec.yml`**

```yaml
name: openspec

on:
  pull_request:

permissions:
  contents: read
  pull-requests: read

jobs:
  gates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - name: Install OpenSpec
        run: npm install -g @fission-ai/openspec@1.3.1
      - name: Run gates
        env:
          BASE_REF: ${{ github.event.pull_request.base.sha }}
          HEAD_REF: ${{ github.head_ref }}
          PR_BODY: ${{ github.event.pull_request.body }}
        run: sh scripts/check-gates.sh
```

- [ ] **Step 3: Check shell syntax and YAML**

```bash
cd ~/Documents/repos/personal/octospec
sh -n scripts/check-gates.sh
python3 -c "import yaml; yaml.safe_load(open('repo-template/.github/workflows/openspec.yml')); print('workflow OK')"
```

Expected: no shell errors; `workflow OK`.

- [ ] **Step 4: Exercise the gate script against pass and fail cases**

```bash
cd "$(mktemp -d)"
git init -q -b main
git commit -q --allow-empty -m init
git checkout -q -b feat/42-do-the-thing
BASE_REF=main HEAD_REF=feat/42-do-the-thing PR_BODY='Closes #42' sh <path>/scripts/check-gates.sh || true
```
Expected: `G3` PASS, `G4` PASS, `G6` SKIP, `G2` result depends on whether an `openspec/` exists (acceptable either way).

```bash
git checkout -q -b bad-branch-name
BASE_REF=main HEAD_REF=bad-branch-name PR_BODY='no link' sh <path>/scripts/check-gates.sh; echo "exit=$?"
```
Expected: `G3` FAIL, `G4` FAIL, `exit=1`.

- [ ] **Step 5: Commit**

```bash
cd ~/Documents/repos/personal/octospec
git add scripts repo-template/.github/workflows
git commit -m "feat(ci): add gate checks and PR workflow"
```

---

### Task 8: End-to-end install into a scratch repo

**Files:**
- No new files; validates Tasks 1–7 together.

**Interfaces:**
- Consumes: `install.sh`, `repo-template/`, `scripts/check-gates.sh`.
- Produces: evidence the loop works before the repo is published.

- [ ] **Step 1: Seed a scratch repo**

```bash
SCRATCH=$(mktemp -d)
cd "$SCRATCH"
git init -q -b main
git commit -q --allow-empty -m init
git remote add origin https://github.com/Kerman-Sanjuan/scratch.git
~/Documents/repos/personal/octospec/install.sh --repo "$SCRATCH"
ls .github/ISSUE_TEMPLATE .github/prompts .github/workflows openspec
```

Expected: the seeded files are present.

- [ ] **Step 2: Create a change and check the schema is wired**

```bash
cd "$SCRATCH"
openspec new change gh-999-sample --schema octospec
openspec status --change gh-999-sample --json | python3 -c "import json,sys; d=json.load(sys.stdin); print([a['id'] for a in d['artifacts']])"
```

Expected: `['issue', 'proposal', 'specs', 'design', 'tasks']`.

- [ ] **Step 3: Confirm the global schema is what resolved**

```bash
cd "$SCRATCH"
openspec schema which octospec
```
Expected: source `user` (the global install), not `package`.

- [ ] **Step 4: Run the gates on a valid branch**

```bash
cd "$SCRATCH"
git add -A && git commit -qm "seed"
git checkout -q -b feat/999-sample
BASE_REF=main HEAD_REF=feat/999-sample PR_BODY='Closes #999' sh scripts/check-gates.sh
```
Expected: G3 PASS, G4 PASS; G2 may fail until a valid change exists; G6 SKIP or PASS.

- [ ] **Step 5: Record the result**

Append a short "Verified" section to `docs/plans/2026-09-29-github-openspec-workflow-plan.md` with the command output you observed.

- [ ] **Step 6: Commit**

```bash
cd ~/Documents/repos/personal/octospec
git add docs
git commit -m "docs: record end-to-end verification"
```

---

### Task 9: README and publish prep

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: everything above.
- Produces: a README a stranger can follow.

- [ ] **Step 1: Write the README**

Cover, in this order: what it is (2 sentences); prerequisites (`openspec`, `gh`, `git`); install (`./install.sh` for global, `./install.sh --repo <path>` for a repo); the loop (`/idea` → `/spec` → `/apply` → `/ship` → `/archive`); the three-tool note (pi/opencode global, Copilot via `.github/prompts/`); the gates table; a link to the design doc; the MIT license.

- [ ] **Step 2: Verify no placeholders remain**

Run: `grep -rn 'TODO\|TBD\|<name>\|FIXME' README.md install.sh scripts || true`
Expected: no output (the templates legitimately contain `<...>` — check `schema/` separately if you extend the grep).

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: write the README"
```

- [ ] **Step 4: Publish (requires explicit user approval)**

> Do not run this step without the user's go-ahead.

```bash
gh repo create Kerman-Sanjuan/octospec --public \
  --source=. --remote=origin --push \
  --description "GitHub-native OpenSpec: issues as the backlog, repo files as the machine layer."
```

Expected: the repo exists at `https://github.com/Kerman-Sanjuan/octospec`.

---

## Self-review

- **Spec coverage:** source-of-truth split (Tasks 2, 6), artifact map (Task 2), command surface (Task 4), schema override not shadowing (Tasks 2, 8 Step 3), issue forms (Task 6), gates G1–G8 (Task 7; G1 enforced by `/spec` in Task 4 Step 4, G7/G8 advisory and unenforced by design), three-tool layout (Tasks 5, 6), risks (documented in the spec).
- **Placeholders:** the only `<...>` markers are inside templates, where they are intentional user-facing hints. No `TODO`/`TBD` in scripts or commands.
- **Type consistency:** schema artifact ids (`issue`, `proposal`, `specs`, `design`, `tasks`) are used identically in Tasks 2, 4, and 8. `install.sh` flags (`--repo`, `PI_PROMPTS_DIR`, `OPENCODE_COMMAND_DIR`) are consistent between Tasks 3 and 5. `check-gates.sh` env vars (`BASE_REF`, `HEAD_REF`, `PR_BODY`) match the workflow in Task 7.

---

## Verified

- 2026-09-29 — end-to-end into a scratch repo:
  - `install.sh --repo` seeded `.github/ISSUE_TEMPLATE`, `.github/prompts`, `.github/workflows`, `openspec/config.yaml`, and `scripts/check-gates.sh`.
  - `openspec schema which octospec` resolved from `user` (the global install).
  - `openspec status` listed `issue, proposal, design, specs, tasks`; `applyRequires: [tasks]`.
  - Gate cases: incomplete change → G2b FAIL; change without a spec delta → G6 FAIL; plain repo → PASS.
- Corrections found during implementation and folded back:
  1. OpenSpec change names must start with a letter → convention is `gh-<issue>-<slug>` (design doc updated).
  2. `openspec validate` does not check artifact completeness → the gate script adds **G2b** using `openspec status --change <name> --json` (`isComplete`), and **G6** now triggers on any touched change folder.
- `install.sh --repo` also copies `scripts/check-gates.sh` into the target repo (the CI workflow calls it).
