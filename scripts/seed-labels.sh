#!/usr/bin/env sh
# Provision the octospec workflow labels into a repository.
#   scripts/seed-labels.sh [<repo-path>]   # default: current directory
#
# Idempotent: uses `gh label create --force`, so re-running updates the
# label colour/description and never fails on an existing label.
set -eu

TARGET=${1:-.}

if ! command -v gh >/dev/null 2>&1; then
  printf 'seed-labels: gh not found; skipping label provisioning\n' >&2
  exit 0
fi

if ! git -C "$TARGET" rev-parse --git-dir >/dev/null 2>&1; then
  printf 'seed-labels: %s is not a git repository; skipping\n' "$TARGET" >&2
  exit 0
fi

printf 'seed-labels: provisioning workflow labels in %s\n' "$TARGET"

label() {
  # name colour description
  ( cd "$TARGET" && gh label create "$1" --color "$2" --description "$3" --force >/dev/null )
  printf '  %s\n' "$1"
}

label "type:feature"       "a2eeef" "New capability"
label "type:bug"           "d73a4a" "Something behaves differently than expected"
label "status:backlog"     "d4c5f9" "Queued for planning"
label "status:spec-ready"  "0e8a16" "Spec artifacts created and validated"
label "status:in-progress" "fbca04" "Work in progress"
label "status:in-review"   "1d76db" "Under review"

printf 'seed-labels: done\n'
