#!/usr/bin/env bash
set -euo pipefail

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)        ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *)             echo "error: unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

DEST="${GRIMOIRE_INSTALL_DIR:-/usr/local/bin}/grimoire"
BIN="grimoire-${OS}-${ARCH}"
URL="https://github.com/jeffreytse/grimoire/releases/latest/download/${BIN}"

echo "Downloading grimoire (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "$DEST"
chmod +x "$DEST"
echo "Installed: $DEST"
echo ""
echo "Next:"
echo "  grimoire update    # clone skill library → ~/.grimoire"
echo "  grimoire install   # install to all detected AI agents"
