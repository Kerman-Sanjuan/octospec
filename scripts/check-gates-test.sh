#!/usr/bin/env sh
# Test the gate script: G3 skips off a PR ref, G1 checks the linked issue.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

out=$(BASE_REF=HEAD HEAD_REF=main sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "SKIP G3"; then
  echo "FAIL: G3 is not skipped on main" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

if ! printf '%s\n' "$out" | grep -q "SKIP G1"; then
  echo "FAIL: G1 is not skipped when no change is in the diff" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

out=$(BASE_REF=HEAD HEAD_REF=feat/6-gates-non-pr-ref sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "PASS G3"; then
  echo "FAIL: G3 is not checked on a PR branch" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

# G3/G4: a Dependabot pull request is exempt from both. A non-bot pull request on
# the same branch and body still fails. No gh or network call is involved.
dep_branch='dependabot/go_modules/cli/x'
dep_body='Bumps x from 1 to 2.'

out=$(BASE_REF=HEAD HEAD_REF="$dep_branch" PR_AUTHOR='dependabot[bot]' PR_BODY="$dep_body" sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "SKIP G3"; then
  echo "FAIL: G3 is not skipped for a dependabot PR" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
if ! printf '%s\n' "$out" | grep -q "SKIP G4"; then
  echo "FAIL: G4 is not skipped for a dependabot PR" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

out=$(BASE_REF=HEAD HEAD_REF="$dep_branch" PR_BODY="$dep_body" sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "FAIL G3"; then
  echo "FAIL: G3 is not enforced for a non-bot PR" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
if ! printf '%s\n' "$out" | grep -q "FAIL G4"; then
  echo "FAIL: G4 is not enforced for a non-bot PR" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

# G1: a change whose linked issue is missing a section fails; a complete body
# passes. The gh call is faked, so the test stays hermetic.
g1=$(mktemp -d)
trap 'rm -rf "$g1"' EXIT
(
  cd "$g1"
  git init -q
  git config user.email "gates@octospec.test"
  git config user.name gates
  mkdir -p openspec/changes/gh-1-x/specs/x
  printf 'github:\n  issue: "1"\n' > openspec/changes/gh-1-x/.openspec.yaml
  printf '## ADDED Requirements\n### Requirement: x\nA SHALL x.\n#### Scenario: y\n- **WHEN** a\n- **THEN** b.\n' > openspec/changes/gh-1-x/specs/x/spec.md
  git add -A
  git commit -qm base
  git branch -m main
  git checkout -qb feat/1-x
  printf 'more\n' >> openspec/changes/gh-1-x/proposal.md
  git add -A
  git commit -qm "feat: change"

  mkdir -p bin
  incomplete() { printf '#!/bin/sh\nprintf "### User story\\n\\nA.\\n\\n### Context\\n\\nB.\\n\\n### Requirements\\n\\nC.\\n"\n'; }
  complete() { printf '#!/bin/sh\nprintf "### User story\\n\\nA.\\n\\n### Context\\n\\nB.\\n\\n### Requirements\\n\\nC.\\n\\n### Success criteria\\n\\nD.\\n"\n'; }

  incomplete > bin/gh
  chmod +x bin/gh
  out=$(PATH="$PWD/bin:$PATH" BASE_REF=main HEAD_REF=feat/1-x sh "$root/scripts/check-gates.sh" 2>&1 || true)
  case "$out" in
    *"FAIL G1 issue #1"*) ;;
    *) echo "FAIL: G1 did not flag a missing section" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
  esac

  complete > bin/gh
  chmod +x bin/gh
  out=$(PATH="$PWD/bin:$PATH" BASE_REF=main HEAD_REF=feat/1-x sh "$root/scripts/check-gates.sh" 2>&1 || true)
  case "$out" in
    *"PASS G1 issue #1"*) ;;
    *) echo "FAIL: G1 did not pass a complete issue" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
  esac
)

# G3: a tag ref (as the release-acceptance workflow passes on a tag push) must
# skip, not fail. This is the regression that reddened main after v1.0.0.
out=$(BASE_REF=HEAD HEAD_REF=v1.0.0 sh scripts/check-gates.sh 2>&1 || true)
if ! printf '%s\n' "$out" | grep -q "SKIP G3"; then
  echo "FAIL: G3 is not skipped on a tag ref" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
if printf '%s\n' "$out" | grep -q "FAIL G3"; then
  echo "FAIL: G3 failed on a tag ref" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

# G2b and G7 must report a SKIP (not pass silently) when openspec is missing.
# Hide openspec by running with a PATH that contains only the essentials.
(
  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT
  mkdir -p "$tmp/bin"
  for tool in git grep sed awk tr basename printf mktemp sh; do
    p=$(command -v "$tool" 2>/dev/null || true)
    [ -n "$p" ] && ln -s "$p" "$tmp/bin/$tool"
  done
  out=$(PATH="$tmp/bin" BASE_REF=HEAD HEAD_REF=main sh "$root/scripts/check-gates.sh" 2>&1 || true)
  case "$out" in
    *"SKIP G2b"*) ;;
    *) echo "FAIL: G2b did not report a SKIP without openspec" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
  esac
  case "$out" in
    *"SKIP G7"*) ;;
    *) echo "FAIL: G7 did not report a SKIP without openspec" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
  esac
)

echo "gate tests passed"
