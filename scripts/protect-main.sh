#!/usr/bin/env sh
# Apply the main branch protection. Reproducible form of the settings.
#
#   scripts/protect-main.sh [owner/repo]
#
# Requires pull requests and the cli and gates checks, but lets the owner push
# maintenance commits such as an archive (enforce_admins: false).
set -eu

REPO="${1:-$(gh repo view --json nameWithOwner --jq .nameWithOwner)}"

gh api --method PUT "repos/$REPO/branches/main/protection" \
  -H "Accept: application/vnd.github+json" --input - <<'JSON' >/dev/null
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["cli", "gates"]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 0
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
JSON

printf 'protected %s main: requires cli and gates, owner may push maintenance\n' "$REPO"
