#!/usr/bin/env sh
# Test that the gate script skips the PR-only branch check off a PR ref.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

out=$(HEAD_REF=main sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "SKIP G3"; then
  echo "FAIL: G3 is not skipped on main" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

out=$(HEAD_REF=feat/6-gates-non-pr-ref sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "PASS G3"; then
  echo "FAIL: G3 is not checked on a PR branch" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

echo "gate tests passed"
