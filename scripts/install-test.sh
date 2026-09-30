#!/usr/bin/env sh
# Test the installer checksum verification and version pinning without GitHub.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
installer="$root/install.sh"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  arm64 | aarch64) arch=arm64 ;;
  *) arch=amd64 ;;
esac
asset="octospec_${os}_${arch}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

filehash() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

make_release() {
  # $1 is the version subpath, for example latest/download or download/v9.9.9.
  dir="$tmp/releases/$1"
  mkdir -p "$dir"
  printf '#!/bin/sh\nexit 0\n' > "$dir/$asset"
  printf '%s  %s\n' "$(filehash "$dir/$asset")" "$asset" > "$dir/checksums.txt"
}

run_installer() {
  OCTOSPEC_BASE_URL="file://$tmp/releases" \
  OCTOSPEC_BIN_DIR="$tmp/bin" \
    sh "$installer" "$@"
}

# 1. The latest release installs and passes the checksum.
make_release "latest/download"
rm -rf "$tmp/bin"
run_installer --tool claude >/dev/null
test -x "$tmp/bin/octospec"

# 2. A pinned version installs.
make_release "download/v9.9.9"
rm -rf "$tmp/bin"
run_installer --version v9.9.9 --tool claude >/dev/null
test -x "$tmp/bin/octospec"

# 3. A mismatch fails and installs nothing.
dir="$tmp/releases/latest/download"
printf 'tampered\n' > "$dir/$asset"
rm -rf "$tmp/bin"
if run_installer --tool claude >/dev/null 2>&1; then
  echo "FAIL: a bad checksum should fail" >&2
  exit 1
fi
test ! -e "$tmp/bin/octospec"

echo "installer tests passed"
