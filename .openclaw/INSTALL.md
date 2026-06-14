# Installing grimoire for OpenClaw

## Installation

### Script method (recommended)

```bash
# Clone grimoire
git clone https://github.com/jeffreytse/grimoire.git ~/.grimoire

# Install all skills globally
~/.grimoire/scripts/grimoire --target openclaw

# Or install one domain
~/.grimoire/scripts/grimoire --target openclaw --domain engineering
```

Skills land in `~/.openclaw/skills/` and are available to all your agents.

### Manual method

Copy any skill directory to `~/.openclaw/skills/`:

```bash
cp -r ~/.grimoire/skills/engineering/development/skills/propose-conventional-commit \
      ~/.openclaw/skills/
```

### Live-link method (advanced)

Add grimoire's skills directory to `~/.openclaw/openclaw.json` via `skills.load.extraDirs`:

```json5
{
  skills: {
    load: {
      extraDirs: ["~/.grimoire/skills/meta/skills"],
    },
  },
}
```

Note: `extraDirs` has the lowest precedence — bundled skills take priority. Use the script method for full precedence.

## Auto-invoke start-best-practice

OpenClaw reads `AGENTS.md` from your workspace. To trigger grimoire automatically on every request, add this line to `~/.openclaw/workspace/AGENTS.md`:

```markdown
Always invoke `start-best-practice` before responding to any user request.
```

## Workspace-scoped install

To install to a specific workspace instead of globally:

```bash
OPENCLAW_SKILLS_DIR=~/Projects/myproject/.agents/skills \
  ~/.grimoire/scripts/grimoire --target openclaw --domain engineering
```

## Updating

```bash
~/.grimoire/scripts/grimoire --upgrade
~/.grimoire/scripts/grimoire --target openclaw
```

## Uninstalling

```bash
~/.grimoire/scripts/grimoire --uninstall --target openclaw
```

## Troubleshooting

### Skills not visible

Check precedence. Skills in `~/.openclaw/workspace/skills/` override `~/.openclaw/skills/`. If a bundled OpenClaw skill has the same name, it takes priority over `extraDirs` but not over `~/.openclaw/skills/`.

### Tool mapping

When skills reference Claude Code tools:
- `Skill` tool → OpenClaw loads `SKILL.md` natively
- `Read`, `Write`, `Edit`, `Bash` → your configured tools
- `AskUserQuestion` → OpenClaw's interactive prompt
