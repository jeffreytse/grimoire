#!/usr/bin/env bash
# Validates SKILL.md files against STANDARD.md requirements.
# Usage: validate-skill.sh <path/to/SKILL.md> [...]
#        validate-skill.sh --all          (validate every SKILL.md in the repo)
# Exit: 0 = all pass, 1 = one or more failures

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PASS=0
FAIL=1
errors=0
checked=0

red()   { printf '\033[0;31m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }
warn()  { printf '\033[0;33m%s\033[0m\n' "$*"; }

check_skill() {
  local file="$1"
  local file_errors=0

  if [[ ! -f "$file" ]]; then
    red "MISSING: $file"
    (( errors++ )) || true
    return
  fi

  local lines
  lines=$(wc -l < "$file")

  # Extract frontmatter block (between first two ---)
  local frontmatter
  frontmatter=$(awk '/^---/{found++; if(found==2) exit} found==1{print}' "$file")

  local name description source tags

  name=$(echo "$frontmatter"     | grep -E '^name:' | head -1 | sed 's/^name:[[:space:]]*//')
  description=$(echo "$frontmatter" | grep -E '^description:' | head -1 | sed 's/^description:[[:space:]]*//')
  source=$(echo "$frontmatter"   | grep -E '^source:' | head -1 | sed 's/^source:[[:space:]]*//')
  tags=$(echo "$frontmatter"     | grep -E '^tags:' | head -1 | sed 's/^tags:[[:space:]]*//')

  printf '  Checking: %s\n' "$file"

  # --- Frontmatter checks ---

  if [[ -z "$name" ]]; then
    red "    [FAIL] name: field missing"
    (( file_errors++ )) || true
  else
    local verb
    verb=$(echo "$name" | cut -d'-' -f1)
    local approved_verbs="propose write review audit design calculate diagnose optimize suggest deprecate plan negotiate apply run"
    local rejected_verbs="do handle manage improve get use help"
    if echo "$rejected_verbs" | grep -qw "$verb"; then
      red "    [FAIL] name: verb '$verb' is a rejected verb (too vague) — see STANDARD.md"
      (( file_errors++ )) || true
    elif ! echo "$approved_verbs" | grep -qw "$verb"; then
      warn "    [WARN] name: verb '$verb' not in approved list — confirm intentional"
    fi
    if [[ ${#name} -gt 50 ]]; then
      red "    [FAIL] name: '$name' exceeds 50 characters (${#name})"
      (( file_errors++ )) || true
    fi
  fi

  if [[ -z "$description" ]]; then
    red "    [FAIL] description: field missing"
    (( file_errors++ )) || true
  elif [[ "$description" != Use\ when* ]]; then
    red "    [FAIL] description: must start with 'Use when' (got: ${description:0:40}...)"
    (( file_errors++ )) || true
  elif [[ ${#description} -gt 500 ]]; then
    red "    [FAIL] description: exceeds 500 characters"
    (( file_errors++ )) || true
  fi

  if [[ -z "$source" ]]; then
    red "    [FAIL] source: field missing"
    (( file_errors++ )) || true
  fi

  if [[ -z "$tags" ]]; then
    red "    [FAIL] tags: field missing"
    (( file_errors++ )) || true
  else
    local tag_count
    tag_count=$(echo "$tags" | tr ',' '\n' | grep -c '\S' || true)
    if [[ "$tag_count" -lt 3 ]]; then
      red "    [FAIL] tags: need at least 3 tags, found $tag_count"
      (( file_errors++ )) || true
    elif [[ "$tag_count" -gt 8 ]]; then
      warn "    [WARN] tags: $tag_count tags — consider trimming to ≤8"
    fi
  fi

  # --- Required sections ---

  if ! grep -q '^## Why This Is Best Practice' "$file"; then
    red "    [FAIL] missing '## Why This Is Best Practice' section"
    (( file_errors++ )) || true
  fi

  if ! grep -qE '^## (Steps|Core Pattern)' "$file"; then
    red "    [FAIL] missing '## Steps' or '## Core Pattern' section"
    (( file_errors++ )) || true
  fi

  # --- Emerging skill check ---

  local emerging
  emerging=$(echo "$frontmatter" | grep -E '^emerging:' | head -1 | sed 's/^emerging:[[:space:]]*//' || true)
  if [[ "$emerging" == "true" ]]; then
    if ! grep -q '\*\*Status:\*\* Emerging' "$file"; then
      red "    [FAIL] emerging: true but missing '**Status:** Emerging' line in Why section"
      (( file_errors++ )) || true
    fi
  fi

  # --- Lifecycle checks ---

  local stable deprecated deprecated_by practitioner
  stable=$(echo "$frontmatter"        | grep -E '^stable:'        | head -1 | sed 's/^stable:[[:space:]]*//'        || true)
  deprecated=$(echo "$frontmatter"    | grep -E '^deprecated:'    | head -1 | sed 's/^deprecated:[[:space:]]*//'    || true)
  deprecated_by=$(echo "$frontmatter" | grep -E '^deprecated_by:' | head -1 | sed 's/^deprecated_by:[[:space:]]*//' || true)
  practitioner=$(echo "$frontmatter"  | grep -E '^practitioner:'  | head -1 | sed 's/^practitioner:[[:space:]]*//'  || true)

  if [[ "$practitioner" == "true" ]] && [[ "$stable" == "true" ]]; then
    red "    [FAIL] practitioner: and stable: cannot both be true — stable requires broad cross-industry adoption"
    (( file_errors++ )) || true
  fi

  if [[ "$deprecated" == "true" ]] && [[ -z "$deprecated_by" ]]; then
    red "    [FAIL] deprecated: true but missing deprecated_by: field (use skill-name or 'none')"
    (( file_errors++ )) || true
  fi

  if [[ -n "$deprecated_by" ]] && [[ "$deprecated" != "true" ]]; then
    warn "    [WARN] deprecated_by: set but deprecated: true not set"
  fi

  if [[ "$stable" == "true" ]] && [[ "$emerging" == "true" ]]; then
    red "    [FAIL] stable: and emerging: cannot both be true"
    (( file_errors++ )) || true
  fi

  if [[ "$deprecated" == "true" ]] && [[ "$emerging" == "true" ]]; then
    red "    [FAIL] deprecated: and emerging: cannot both be true"
    (( file_errors++ )) || true
  fi

  # --- Why section content checks ---

  if grep -q '^## Why This Is Best Practice' "$file"; then
    if ! grep -q '\*\*Adopted by:\*\*' "$file"; then
      red "    [FAIL] '## Why This Is Best Practice' missing **Adopted by:** line"
      (( file_errors++ )) || true
    fi
    if ! grep -q '\*\*Impact:\*\*' "$file"; then
      red "    [FAIL] '## Why This Is Best Practice' missing **Impact:** line"
      (( file_errors++ )) || true
    fi
  fi

  # --- Size check ---

  if [[ "$lines" -lt 50 ]]; then
    red "    [FAIL] only $lines lines — minimum is 50 (add steps, examples, or edge cases)"
    (( file_errors++ )) || true
  elif [[ "$lines" -gt 300 ]]; then
    red "    [FAIL] $lines lines exceeds 300-line limit — consider splitting"
    (( file_errors++ )) || true
  fi

  if [[ "$file_errors" -eq 0 ]]; then
    green "    [PASS] $name"
  else
    red "    [FAIL] $file_errors error(s)"
    (( errors += file_errors )) || true
  fi

  (( checked++ )) || true
}

# --- Entry point ---

files=()

if [[ "${1:-}" == "--all" ]]; then
  while IFS= read -r -d '' f; do
    files+=("$f")
  done < <(find "$REPO_ROOT/skills" -name 'SKILL.md' -print0 | sort -z)
else
  files=("$@")
fi

if [[ ${#files[@]} -eq 0 ]]; then
  echo "Usage: validate-skill.sh <SKILL.md> [...]"
  echo "       validate-skill.sh --all"
  exit 1
fi

echo "grimoire skill validator"
echo "========================"

for f in "${files[@]}"; do
  check_skill "$f"
done

echo ""
echo "Checked: $checked  Errors: $errors"

if [[ "$errors" -gt 0 ]]; then
  red "FAIL — $errors error(s) found"
  exit 1
else
  green "PASS — all skills valid"
  exit 0
fi
