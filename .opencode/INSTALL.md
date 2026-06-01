# Installing grimoire for OpenCode

## Installation

Add grimoire to the `plugin` array in your `opencode.json` (global or project-level):

```json
{
  "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire.git"]
}
```

Restart OpenCode. The plugin registers all grimoire skills automatically.

## Usage

Use OpenCode's native `skill` tool:

```
use skill tool to list skills
use skill tool to load engineering/development/propose-conventional-commit
```

## Updating

To pin a specific version:

```json
{
  "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire.git#v1.0.0"]
}
```

## Troubleshooting

### Skills not found

Check that the plugin is loading:
```bash
opencode run --print-logs "hello" 2>&1 | grep -i grimoire
```

### Tool mapping

When skills reference Claude Code tools:
- `Skill` tool → OpenCode's native `skill` tool
- `Read`, `Write`, `Edit`, `Bash` → your native tools
- `TodoWrite` → `todowrite`
