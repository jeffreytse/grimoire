#!/usr/bin/env bash
# Grimoire Skill Standard conformance test suite.
# Tests validate-skill.sh against canonical fixtures in schema/tests/.
# Usage: scripts/test-schema.sh
# Exit: 0 = all conformance tests pass, 1 = ≥1 failure

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR="$REPO_ROOT/scripts/validate-skill.sh"
TESTS_DIR="$REPO_ROOT/schema/tests"

pass=0
fail=0

red()   { printf '\033[0;31m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }

run_fixture() {
  local kind="$1" fixture="$2"
  local name
  name=$(basename "$(dirname "$fixture")")

  if [[ "$kind" == "valid" ]]; then
    if bash "$VALIDATOR" "$fixture" >/dev/null 2>&1; then
      green "[PASS] valid/$name"
      (( pass++ )) || true
    else
      red "[FAIL] valid/$name — expected PASS, got FAIL"
      (( fail++ )) || true
    fi
  else
    if ! bash "$VALIDATOR" "$fixture" >/dev/null 2>&1; then
      green "[PASS] invalid/$name"
      (( pass++ )) || true
    else
      red "[FAIL] invalid/$name — expected FAIL, got PASS"
      (( fail++ )) || true
    fi
  fi
}

echo "grimoire conformance suite"
echo "=========================="

while IFS= read -r -d '' f; do
  run_fixture "valid" "$f"
done < <(find "$TESTS_DIR/valid" -name 'SKILL.md' -print0 | sort -z)

while IFS= read -r -d '' f; do
  run_fixture "invalid" "$f"
done < <(find "$TESTS_DIR/invalid" -name 'SKILL.md' -print0 | sort -z)

echo ""
echo "Passed: $pass  Failed: $fail"
if [[ "$fail" -gt 0 ]]; then
  red "FAIL — $fail conformance test(s) failed"
  exit 1
else
  green "PASS — all conformance tests pass"
fi
