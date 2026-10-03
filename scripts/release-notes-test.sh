#!/usr/bin/env sh
# Test scripts/release-notes.sh: it extracts a version's section, strips a
# leading v, and falls back when the version is missing.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

# 1. The real v1.0.0 section is extracted and is grouped by impact.
out=$(sh scripts/release-notes.sh v1.0.0)
case "$out" in
  *"### Added"*) ;;
  *) echo "FAIL: v1.0.0 notes have no '### Added'" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
esac
case "$out" in
  *"### Fixed"*) ;;
  *) echo "FAIL: v1.0.0 notes have no '### Fixed'" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
esac
# The section must not bleed into a following version heading.
if printf '%s\n' "$out" | grep -q '^## '; then
  echo "FAIL: the section leaked into a following version heading" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

# 2. A leading v is optional.
a=$(sh scripts/release-notes.sh v1.0.0)
b=$(sh scripts/release-notes.sh 1.0.0)
[ "$a" = "$b" ] || { echo "FAIL: 'v1.0.0' and '1.0.0' differ" >&2; exit 1; }

# 3. A missing version falls back instead of failing.
out=$(sh scripts/release-notes.sh v9.9.9)
case "$out" in
  *"9.9.9"*) ;;
  *) echo "FAIL: a missing version has no fallback" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
esac

# 4. A synthetic change to CHANGELOG is honored, so the helper reads the file.
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
cat > "$tmp/CHANGELOG" <<'EOF'
# Changelog

## 2.0.0 - 2030-01-01

### Added
- A thing.

## 1.0.0 - 2026-10-03

### Fixed
- Something.
EOF
out=$(CHANGELOG="$tmp/CHANGELOG" sh scripts/release-notes.sh v2.0.0)
case "$out" in
  *"A thing."*) ;;
  *) echo "FAIL: did not read the synthetic changelog" >&2; printf '%s\n' "$out" >&2; exit 1 ;;
esac
if printf '%s\n' "$out" | grep -q "Something."; then
  echo "FAIL: the 2.0.0 section included the older 1.0.0 entry" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

echo "release-notes tests passed"
