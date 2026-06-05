#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_FILE="$REPO_ROOT/VERSION"

NEW_VERSION="${1:-}"

if [[ -z "$NEW_VERSION" ]]; then
  echo "Usage: $0 <new-version>   e.g. $0 1.1.0"
  echo "Current version: $(cat "$VERSION_FILE")"
  exit 1
fi

if ! [[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: version must be semver (x.y.z), got: $NEW_VERSION"
  exit 1
fi

OLD_VERSION="$(cat "$VERSION_FILE")"

if [[ "$OLD_VERSION" == "$NEW_VERSION" ]]; then
  echo "Already at $NEW_VERSION — nothing to do."
  exit 0
fi

echo "$NEW_VERSION" > "$VERSION_FILE"

find "$REPO_ROOT" \( \
  -path '*/.claude-plugin/plugin.json' -o \
  -path '*/.codex-plugin/plugin.json'  -o \
  -path '*/.cursor-plugin/plugin.json' -o \
  -name 'gemini-extension.json' \
\) | while read -r f; do
  sed -i '' "s/\"version\": \"${OLD_VERSION}\"/\"version\": \"${NEW_VERSION}\"/" "$f"
done

COUNT=$(find "$REPO_ROOT" \( \
  -path '*/.claude-plugin/plugin.json' -o \
  -path '*/.codex-plugin/plugin.json'  -o \
  -path '*/.cursor-plugin/plugin.json' -o \
  -name 'gemini-extension.json' \
\) | wc -l | tr -d ' ')

echo "Bumped $OLD_VERSION → $NEW_VERSION in VERSION + $COUNT manifest files."
