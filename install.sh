#!/usr/bin/env sh
# Install the octospec CLI from a GitHub release, then install the commands.
#
#   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh -s -- --tool copilot
#   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh -s -- --version v1.1.0
#
# The installer verifies the release checksum before it installs the binary.
# Set OCTOSPEC_BIN_DIR to change where the binary lands (default ~/.local/bin).
# Set OCTOSPEC_BASE_URL to override the releases base (used by the tests).
set -eu

REPO=Kerman-Sanjuan/octospec
BIN=octospec

usage() {
  cat <<'EOF'
usage: install.sh [--version <tag>] [--tool <name>]... [--repo <path>]

  --version <tag>  install a specific release instead of the latest
  other flags are passed to `octospec install`
EOF
}

version=""
passthrough=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      [ "$#" -ge 2 ] || { echo "octospec: --version needs a value" >&2; exit 1; }
      version="$2"
      shift 2
      ;;
    --version=*)
      version="${1#--version=}"
      shift
      ;;
    -h | --help)
      usage
      exit 0
      ;;
    *)
      escaped=$(printf '%s' "$1" | sed "s/'/'\\\\''/g")
      passthrough="$passthrough '$escaped'"
      shift
      ;;
  esac
done
eval "set -- $passthrough"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  arm64 | aarch64) arch=arm64 ;;
  x86_64 | amd64) arch=amd64 ;;
  *)
    echo "octospec: unsupported architecture $(uname -m)" >&2
    exit 1
    ;;
esac

asset="${BIN}_${os}_${arch}"
rel=${OCTOSPEC_BASE_URL:-https://github.com/$REPO/releases}
if [ -n "$version" ]; then
  base="$rel/download/$version"
else
  base="$rel/latest/download"
fi

filehash() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

echo "octospec: downloading $base/$asset"
curl -fsSL "$base/$asset" -o "$work/$asset"
curl -fsSL "$base/checksums.txt" -o "$work/checksums.txt"

expected=$(awk -v f="$asset" '$2 == f { print $1 }' "$work/checksums.txt")
if [ -z "$expected" ]; then
  echo "octospec: no checksum for $asset in checksums.txt" >&2
  exit 1
fi
actual=$(filehash "$work/$asset")
if [ "$actual" != "$expected" ]; then
  echo "octospec: checksum mismatch for $asset" >&2
  echo "  expected $expected" >&2
  echo "  got      $actual" >&2
  exit 1
fi

dest=${OCTOSPEC_BIN_DIR:-$HOME/.local/bin}
mkdir -p "$dest"
chmod +x "$work/$asset"
mv "$work/$asset" "$dest/$BIN"
echo "octospec: installed $dest/$BIN"

case ":$PATH:" in
  *":$dest:"*) ;;
  *) echo "octospec: add $dest to your PATH" ;;
esac

# Install the commands into the requested tools (defaults to all).
"$dest/$BIN" install "$@"
