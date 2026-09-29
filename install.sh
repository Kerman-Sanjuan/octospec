#!/usr/bin/env sh
# Install the octospec CLI from a GitHub release, then install the commands.
#
#   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh -s -- --tool copilot
#
# Set OCTOSPEC_BIN_DIR to change where the binary lands (default ~/.local/bin).
set -eu

REPO=Kerman-Sanjuan/octospec
BIN=octospec

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  arm64 | aarch64) arch=arm64 ;;
  x86_64 | amd64) arch=amd64 ;;
  *)
    echo "octospec: unsupported architecture $(uname -m)" >&2
    exit 1
    ;;
esac

dest=${OCTOSPEC_BIN_DIR:-$HOME/.local/bin}
mkdir -p "$dest"
url="https://github.com/$REPO/releases/latest/download/${BIN}_${os}_${arch}"

echo "octospec: downloading $url"
tmp=$(mktemp)
curl -fsSL "$url" -o "$tmp"
chmod +x "$tmp"
mv "$tmp" "$dest/$BIN"
echo "octospec: installed $dest/$BIN"

case ":$PATH:" in
  *":$dest:"*) ;;
  *) echo "octospec: add $dest to your PATH" ;;
esac

# Install the commands into the requested tools (defaults to all).
"$dest/$BIN" install "$@"
