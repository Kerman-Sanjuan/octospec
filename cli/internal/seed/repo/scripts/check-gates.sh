#!/usr/bin/env sh
# Gate checks G1-G8. Runnable locally and in CI.
# Env: BASE_REF (default origin/main), HEAD_REF, PR_BODY, PR_AUTHOR (optional).
set -eu

BASE_REF=${BASE_REF:-origin/main}
HEAD_REF=${HEAD_REF:-$(git rev-parse --abbrev-ref HEAD)}
PR_AUTHOR=${PR_AUTHOR:-}
fail=0

report() { printf '%s %s\n' "$1" "$2"; }

# G1: issue linked to a touched change has the required sections.
# Only runs when an OpenSpec change is in the diff.
change_touched=$(git diff --name-only "$BASE_REF...HEAD" 2>/dev/null | grep -E '^openspec/changes/[^/]+/' | grep -v '^openspec/changes/archive/' || true)
if [ -n "$change_touched" ]; then
  if command -v gh >/dev/null 2>&1; then
    change_dirs=$(printf '%s\n' "$change_touched" | sed -E 's#^(openspec/changes/[^/]+)/.*#\1#' | sort -u)
    checked=0
    for d in $change_dirs; do
      cfg="$d/.openspec.yaml"
      [ -f "$cfg" ] || continue
      issue=$(grep -A2 '^github:' "$cfg" 2>/dev/null | grep 'issue:' | awk '{print $2}' | tr -d '"' || true)
      [ -n "$issue" ] || continue
      checked=1
      # Distinguish "gh cannot read the issue" (auth, permission, network, or a
      # missing issue) from "the issue is incomplete". An unread body must not
      # fail G1, or a seeded repo without a token reports a false negative.
      if ! body=$(gh issue view "$issue" --json body --jq '.body' 2>/dev/null); then
        report SKIP "G1 issue #$issue body: gh could not read it"
        continue
      fi
      missing=""
      for section in "User story" "Context" "Requirements" "Success criteria"; do
        if ! printf '%s' "$body" | grep -qE "^#{2,4}[[:space:]]+$section"; then
          missing="$missing $section"
        fi
      done
      if [ -n "$missing" ]; then
        report FAIL "G1 issue #$issue missing section(s):$missing"
        fail=1
      else
        report PASS "G1 issue #$issue sections complete"
      fi
    done
    if [ "$checked" -eq 0 ]; then
      report SKIP "G1 issue body: no .openspec.yaml found"
    fi
  else
    report SKIP "G1 issue body: gh not available"
  fi
else
  report SKIP "G1 issue body: no OpenSpec change in this PR"
fi

# G3: branch name. A feat/ or fix/ branch (and a dependabot/ branch) is
# checked. Any other ref (main, a tag such as v1.0.0, or another long-lived ref)
# is skipped, because its name cannot tell a PR branch from a release ref.
case "$HEAD_REF" in
  feat/* | fix/* | dependabot/*)
    if [ "$PR_AUTHOR" = "dependabot[bot]" ]; then
      report SKIP "G3 branch name: automated dependency PR"
    elif printf '%s' "$HEAD_REF" | grep -Eq '^(feat|fix)/[0-9]+-[a-z0-9-]+$'; then
      report PASS "G3 branch name"
    else
      report FAIL "G3 branch name: '$HEAD_REF' does not match feat|fix/<issue>-<slug>"
      fail=1
    fi
    ;;
  *)
    report SKIP "G3 branch name: long-lived ref '$HEAD_REF'"
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
else
  report SKIP "G2b completeness: openspec not available"
fi

# G7 (advisory): every unarchived change has all tasks checked.
if command -v openspec >/dev/null 2>&1; then
  open_tasks=""
  for d in openspec/changes/*/; do
    [ -d "$d" ] || continue
    name=$(basename "$d")
    [ "$name" = "archive" ] && continue
    if grep -qE '^[[:space:]]*- \[ \]' "$d/tasks.md" 2>/dev/null; then
      open_tasks="$open_tasks $name"
    fi
  done
  if [ -n "$open_tasks" ]; then
    report WARN "G7 unchecked tasks:$open_tasks"
  else
    report PASS "G7 tasks checked"
  fi
else
  report SKIP "G7 tasks checked: openspec not available"
fi

# G8 (advisory): one unarchived change per issue.
seen_issues=""
dupe_issues=""
for d in openspec/changes/*/; do
  [ -d "$d" ] || continue
  name=$(basename "$d")
  [ "$name" = "archive" ] && continue
  issue=$(grep -A2 '^github:' "$d/.openspec.yaml" 2>/dev/null | grep 'issue:' | awk '{print $2}' | tr -d '"' || true)
  [ -z "$issue" ] && continue
  case " $seen_issues " in
    *" $issue "*) dupe_issues="$dupe_issues $issue" ;;
    *) seen_issues="$seen_issues $issue" ;;
  esac
done
if [ -n "$dupe_issues" ]; then
  report WARN "G8 more than one change per issue:$dupe_issues"
else
  report PASS "G8 one change per issue"
fi

# G6: a PR that touches an OpenSpec change must carry a spec delta.
changed=$(git diff --name-only "$BASE_REF...HEAD" 2>/dev/null || true)
if [ -n "$change_touched" ]; then
  if printf '%s\n' "$changed" | grep -Eq '^openspec/changes/[^/]+/specs/.+\.md$'; then
    report PASS "G6 spec delta"
  else
    report FAIL "G6 spec delta: change touched but no specs/*.md (set skip_specs: true in its .openspec.yaml only for docs/tooling)"
    fail=1
  fi
else
  report SKIP "G6 spec delta: no OpenSpec change in this PR"
fi

# G4: PR body links the issue. Skip for Dependabot, whose body is generated.
body=${PR_BODY:-}
if [ "$PR_AUTHOR" = "dependabot[bot]" ]; then
  report SKIP "G4 PR links issue: automated dependency PR"
elif [ -n "$body" ]; then
  if printf '%s' "$body" | grep -Eq '(Closes|Fixes|Resolves) #[0-9]+'; then
    report PASS "G4 PR links issue"
  else
    report FAIL "G4 PR body has no Closes #n"
    fail=1
  fi
else
  report SKIP "G4 PR body not provided"
fi

# G9 (advisory): every requirement names the check that proves it. A change declares
# openspec/changes/<name>/coverage.md with one line per requirement,
#   - <requirement title> -> <symbol>
# where the symbol is any string that exists in the tree outside openspec/. A requirement with no
# line, or a line pointing at something that is not there, is a requirement nobody can prove.
if [ -n "$change_touched" ]; then
  g9_seen=0
  g9_nocov=""
  g9_missing=""
  g9_unknown=""
  for d in $(printf '%s\n' "$change_touched" | sed -E 's#^(openspec/changes/[^/]+)/.*#\1#' | sort -u); do
    [ -d "$d/specs" ] || continue
    reqs=$(cat "$d"/specs/*/spec.md 2>/dev/null | grep '^### Requirement: ' | sed 's/^### Requirement: //' || true)
    [ -n "$reqs" ] || continue
    g9_seen=1
    cov="$d/coverage.md"
    if [ ! -f "$cov" ]; then
      g9_nocov="$g9_nocov $(basename "$d")"
      continue
    fi
    while IFS= read -r req; do
      [ -n "$req" ] || continue
      if ! grep -qF -- "$req" "$cov"; then
        g9_missing="$g9_missing \"$req\";"
      fi
    done <<EOF
$reqs
EOF
    while IFS= read -r line; do
      case "$line" in
        *' -> '*) ;;
        *) continue ;;
      esac
      sym=$(printf '%s' "$line" | sed 's/^.* -> //' | tr -d '`' | tr -d '\r')
      [ -n "$sym" ] || continue
      if ! grep -rqF --exclude-dir=openspec --exclude-dir=.git -- "$sym" . 2>/dev/null; then
        g9_unknown="$g9_unknown \"$sym\";"
      fi
    done <<EOF
$(grep ' -> ' "$cov" 2>/dev/null || true)
EOF
  done
  if [ "$g9_seen" -eq 0 ]; then
    report SKIP "G9 requirement coverage: no spec delta"
  elif [ -n "$g9_nocov$g9_missing$g9_unknown" ]; then
    report WARN "G9 requirement coverage:${g9_nocov:+ no coverage.md for$g9_nocov}${g9_missing:+ unmapped:$g9_missing}${g9_unknown:+ missing symbol:$g9_unknown}"
  else
    report PASS "G9 requirement coverage"
  fi
else
  report SKIP "G9 requirement coverage: no OpenSpec change in this PR"
fi

if [ "$fail" -ne 0 ]; then
  report FAIL "gates failed"
  exit 1
fi
report PASS "all runnable gates"
