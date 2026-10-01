#!/usr/bin/env sh
# End-to-end test: build the CLI, seed and install into a temporary repository,
# and run the gates. Hermetic: it does not call GitHub.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

if ! command -v openspec >/dev/null 2>&1; then
  echo "SKIP: openspec is not on PATH"
  exit 0
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# 1. Build the CLI.
( cd cli && go build -o "$tmp/octospec" ./cmd/octospec )

repo="$tmp/repo"
mkdir -p "$repo"
(
  cd "$repo"
  git init -q
  git config user.email "e2e@octospec.test"
  git config user.name "e2e"
)

# 2. Seed and install.
"$tmp/octospec" seed --no-labels --repo "$repo" >/dev/null
"$tmp/octospec" install --tool claude --repo "$repo" >/dev/null 2>&1 || true

# 3. Assert the layout.
test -f "$repo/scripts/check-gates.sh"
test -f "$repo/.github/workflows/openspec.yml"
test -f "$repo/.github/ISSUE_TEMPLATE/feature.yml"
test -f "$repo/openspec/config.yaml"
test "$(ls "$repo/.claude/commands" | wc -l | tr -d ' ')" = "6"
test "$(ls "$repo/.claude/agents" | wc -l | tr -d ' ')" = "5"
test -f "$repo/.claude/agents/apply.md"
grep -q "tools:" "$repo/.claude/agents/apply.md"

# 4. The gates run in the seeded repository.
out=$( cd "$repo" && HEAD_REF=main sh scripts/check-gates.sh 2>&1 || true )
case "$out" in
  *"PASS G3"* | *"SKIP G3"*) ;;
  *)
    echo "FAIL: the gate script did not report G3" >&2
    printf '%s\n' "$out" >&2
    exit 1
    ;;
esac

echo "e2e tests passed"
