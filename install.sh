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
for f in "$REPO_DIR"/commands/*.md; do
  [ -e "$f" ] || continue
  cp "$f" "$PI_DIR/"
done
printf 'pi       -> %s\n' "$PI_DIR"

# 3. opencode commands
OC_DIR="${OPENCODE_COMMAND_DIR:-$HOME/.config/opencode/command}"
mkdir -p "$OC_DIR"
for f in "$REPO_DIR"/commands/*.md; do
  [ -e "$f" ] || continue
  cp "$f" "$OC_DIR/"
done
printf 'opencode -> %s\n' "$OC_DIR"

# 4. Optional repo seed
if [ "${1:-}" = "--repo" ]; then
  TARGET=${2:-}
  if [ -z "$TARGET" ] || [ ! -d "$TARGET" ]; then
    printf 'usage: install.sh --repo <existing-repo-path>\n' >&2
    exit 2
  fi
  mkdir -p "$TARGET/.github/ISSUE_TEMPLATE" "$TARGET/.github/prompts" \
           "$TARGET/.github/workflows" "$TARGET/openspec" "$TARGET/scripts"
  cp -R "$REPO_DIR/repo-template/." "$TARGET/"
  cp "$REPO_DIR/scripts/check-gates.sh" "$TARGET/scripts/"
  chmod +x "$TARGET/scripts/check-gates.sh"
  printf 'repo     -> %s\n' "$TARGET"
fi
