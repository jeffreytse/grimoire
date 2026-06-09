#!/usr/bin/env python3
"""Fix unquoted YAML frontmatter values containing colon-space sequences."""
import re
import sys
from pathlib import Path

SKILL_ROOT = Path(__file__).parent.parent / "skills"
FIELDS = {"description", "source"}

def needs_quoting(value: str) -> bool:
    return ": " in value and not (value.startswith('"') or value.startswith("'"))

def quote_value(value: str) -> str:
    if '"' in value:
        return f"'{value}'"
    return f'"{value}"'

def fix_file(path: Path) -> bool:
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines(keepends=True)
    in_front = False
    changed = False
    result = []
    for line in lines:
        if line.strip() == "---":
            in_front = not in_front
            result.append(line)
            continue
        if in_front:
            m = re.match(r'^(\w+): (.+)\n?$', line)
            if m and m.group(1) in FIELDS:
                value = m.group(2)
                if needs_quoting(value):
                    line = f"{m.group(1)}: {quote_value(value)}\n"
                    changed = True
        result.append(line)
    if changed:
        path.write_text("".join(result), encoding="utf-8")
    return changed

def main():
    fixed = 0
    for skill_md in SKILL_ROOT.rglob("SKILL.md"):
        if fix_file(skill_md):
            print(f"  fixed: {skill_md.relative_to(SKILL_ROOT.parent)}")
            fixed += 1
    print(f"\n{fixed} files fixed.")

if __name__ == "__main__":
    main()
