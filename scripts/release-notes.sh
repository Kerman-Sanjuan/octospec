#!/usr/bin/env sh
# Print the CHANGELOG.md section for a version, for use as the release notes.
# Usage: release-notes.sh <version>   (for example v1.0.0 or 1.0.0)
#
# The version is matched against the "## <version>" heading. A leading "v" is
# stripped. When the section is missing, a one-line fallback is printed so a
# release never fails on it.
set -eu

version="${1:-}"
if [ -z "$version" ]; then
  echo "usage: release-notes.sh <version>" >&2
  exit 1
fi
version="${version#v}"

file="${CHANGELOG:-CHANGELOG.md}"
if [ ! -f "$file" ]; then
  printf 'Release %s\n' "$version"
  exit 0
fi

section=$(awk -v ver="$version" '
  $0 ~ "^## " ver "([[:space:]]|$)" { insec=1; next }
  insec && /^## / { exit }
  insec { print }
' "$file")

# Drop leading blank lines.
section=$(printf '%s\n' "$section" | awk 'NF { found=1 } found')

if [ -z "$section" ]; then
  printf 'Release %s\n' "$version"
else
  printf '%s\n' "$section"
fi
