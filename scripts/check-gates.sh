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
