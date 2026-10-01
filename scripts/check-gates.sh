#!/usr/bin/env sh
# Gate checks G2-G6. Runnable locally and in CI.
# Env: BASE_REF (default origin/main), HEAD_REF, PR_BODY (optional).
set -eu

BASE_REF=${BASE_REF:-origin/main}
HEAD_REF=${HEAD_REF:-$(git rev-parse --abbrev-ref HEAD)}
fail=0

report() { printf '%s %s\n' "$1" "$2"; }

# G3: branch name. PR-only: skip on a long-lived branch.
case "$HEAD_REF" in
  main | master)
    report SKIP "G3 branch name: long-lived branch '$HEAD_REF'"
    ;;
  *)
    if printf '%s' "$HEAD_REF" | grep -Eq '^(feat|fix)/[0-9]+-[a-z0-9-]+$'; then
      report PASS "G3 branch name"
    else
      report FAIL "G3 branch name: '$HEAD_REF' does not match feat|fix/<issue>-<slug>"
      fail=1
    fi
    ;;
esac

# G2a: openspec validate (structural validity of specs and changes)
if command -v openspec >/dev/null 2>&1 && openspec validate --all --strict >/dev/null 2>&1; then
  report PASS "G2a openspec validate"
else
  report FAIL "G2a openspec validate"
  fail=1
fi

# G2b: every unarchived change is complete (all apply-required artifacts done).
# `openspec validate` does NOT check completeness; `openspec status` does.
if command -v openspec >/dev/null 2>&1; then
  incomplete=""
  for d in openspec/changes/*/; do
    [ -d "$d" ] || continue
    name=$(basename "$d")
    if [ "$name" = "archive" ]; then
      continue
    fi
    out=$(openspec status --change "$name" --json 2>/dev/null || true)
    if [ -z "$out" ]; then
      continue
    fi
    if ! printf '%s' "$out" | grep -Eq '"isComplete"[[:space:]]*:[[:space:]]*true'; then
      incomplete="$incomplete $name"
    fi
  done
  if [ -n "$incomplete" ]; then
    report FAIL "G2b incomplete change(s):$incomplete"
    fail=1
  else
    report PASS "G2b changes complete"
  fi
fi

# G6: a PR that touches an OpenSpec change must carry a spec delta.
changed=$(git diff --name-only "$BASE_REF...HEAD" 2>/dev/null || true)
change_touched=$(printf '%s\n' "$changed" | grep -E '^openspec/changes/[^/]+/' | grep -v '^openspec/changes/archive/' || true)
if [ -n "$change_touched" ]; then
  if printf '%s\n' "$changed" | grep -Eq '^openspec/changes/[^/]+/specs/.+\.md$'; then
    report PASS "G6 spec delta"
  else
    report FAIL "G6 spec delta: change touched but no specs/*.md (use --skip-specs only for docs/tooling)"
    fail=1
  fi
else
  report SKIP "G6 spec delta: no OpenSpec change in this PR"
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
