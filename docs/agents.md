# Agent Integration Guide

grimoire works across six AI agents. Each has its own plugin format and installation method. This guide covers installation, usage, and tool mapping for each.

---

## Supported agents

| Agent | Plugin file | Install method |
|-------|------------|----------------|
| Claude Code | `.claude-plugin/plugin.json` | `/plugins add github:jeffreytse/grimoire` |
| Codex | `.codex-plugin/plugin.json` | `codex plugin add github:jeffreytse/grimoire` |
| Cursor | `.cursor-plugin/plugin.json` | Cursor plugin marketplace |
| OpenCode | `.opencode/plugins/grimoire.js` | `opencode.json` plugin array |
| Gemini CLI | `gemini-extension.json` + `GEMINI.md` | `gemini extension install github:jeffreytse/grimoire` |
| Agents CLI | `AGENTS.md` | `agents install github:jeffreytse/grimoire` |

---

## Claude Code

**Install:**
```
/plugins add github:jeffreytse/grimoire
```

**Plugin file:** `.claude-plugin/plugin.json`

Loads all domains listed in the `domains` array. Claude Code reads each domain's `plugin.json` to discover skills.

**Using skills:**
```
/suggest-best-practice
/write-best-practice-skill
/review-best-practice-skill
```

Invoke any meta skill by name with a leading slash. For domain skills, describe your situation and `suggest-best-practice` routes automatically.

**Tool names:** Use native Claude Code tool names (Read, Write, Edit, Bash, Skill, TodoWrite).

---

## Codex

**Install:**
```
codex plugin add github:jeffreytse/grimoire
```

**Plugin file:** `.codex-plugin/plugin.json`

Uses the `skills` field pointing to `./skills/` — Codex discovers all SKILL.md files recursively from that path.

**Using skills:**
Invoke skills by name. Codex loads skill content on demand when a skill matches the current task.

**Tool mapping:** See `references/codex-tools.md` in the skills directory for Codex-specific tool name equivalents.

---

## Cursor

**Install:**
Install via the Cursor plugin marketplace, or add manually to your Cursor configuration pointing to `github:jeffreytse/grimoire`.

**Plugin file:** `.cursor-plugin/plugin.json`

Uses the `skills` field pointing to `./skills/` — Cursor discovers skills from that path.

**Using skills:**
Describe your situation in the Cursor chat or Composer panel. Skills activate when the task matches a skill's triggering conditions.

---

## OpenCode

**Install:**
Add to your `opencode.json` (global: `~/.config/opencode/opencode.json`, or project-level):

```json
{
  "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire.git"]
}
```

Restart OpenCode. The plugin registers all grimoire skills automatically.

**Plugin file:** `.opencode/plugins/grimoire.js` (ESM module)

The plugin:
1. Scans `./skills/` recursively for SKILL.md files
2. Injects `AGENTS.md` as bootstrap context
3. Maps Claude Code tool references to OpenCode native equivalents

**Using skills:**
```
use skill tool to load engineering/development/propose-conventional-commit
use skill tool to list skills
```

**Tool mapping:**

| Claude Code tool | OpenCode equivalent |
|-----------------|-------------------|
| `Skill` | `skill` (native) |
| `Read`, `Write`, `Edit`, `Bash` | native tools |
| `TodoWrite` | `todowrite` |

**Pinning a version:**
```json
{
  "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire.git#v1.0.0"]
}
```

**Troubleshooting — skills not found:**
```bash
opencode run --print-logs "hello" 2>&1 | grep -i grimoire
```

---

## Gemini CLI

**Install:**
```
gemini extension install github:jeffreytse/grimoire
```

**Plugin files:** `gemini-extension.json` (extension metadata) + `GEMINI.md` (context injected at session start)

`gemini-extension.json` points to `GEMINI.md` via the `contextFileName` field. Gemini CLI loads `GEMINI.md` into context at startup, giving the agent the skill path convention and domain list.

**Using skills:**
Gemini CLI loads skill metadata at session start. Activate a skill by describing your task — the agent matches to the closest skill. For explicit activation:
```
activate skill: engineering/development/propose-conventional-commit
```

**Key difference from Claude Code:** Gemini CLI uses `activate_skill` rather than slash commands. `GEMINI.md` provides the routing bootstrap; individual skills are activated on demand from their SKILL.md files.

**Tool mapping:** See `GEMINI.md` in the repo root for Gemini-specific tool references.

---

## Agents CLI

**Install:**
```
agents install github:jeffreytse/grimoire
```

**Plugin file:** `AGENTS.md` (at repo root)

`AGENTS.md` provides the same bootstrap context as `GEMINI.md` — skill path convention, domain list, and routing instructions. Agents CLI reads it at session start.

**Using skills:**
Same as Gemini CLI — describe your task and the agent routes to the best skill. `AGENTS.md` contains the context needed for auto-routing.

---

## Adding grimoire to a new agent

If you're integrating grimoire with an agent not listed here:

1. **Check for a plugin format.** Most agents use either a JSON manifest (`plugin.json`) or a context file (markdown loaded at startup).

2. **Point to `./skills/`** — all SKILL.md files live under this path in a consistent structure (`<domain>/<subdomain>/skills/<name>/SKILL.md`).

3. **Inject bootstrap context** — copy `AGENTS.md` or `GEMINI.md` content into whatever context injection mechanism your agent supports. This gives the agent the skill path convention and domain list.

4. **Map tool names** — SKILL.md files reference Claude Code tool names (Read, Write, Edit, Bash, Skill, TodoWrite). Document the equivalents for your agent so skills can be adapted.

5. **Open a PR** — add the new agent's plugin file(s) to the repo and update the supported agents table in this document.

---

## Keeping plugins in sync

When a new domain is added to grimoire (via `design-best-practice-domain`), the following files must be updated:

| File | What to update |
|------|----------------|
| `.claude-plugin/plugin.json` | Add domain path to `domains` array |
| `.claude-plugin/marketplace.json` | Add domain and subdomain entries |
| `GEMINI.md` | Add domain row to domains table |
| `AGENTS.md` | Same update as GEMINI.md |

`.codex-plugin/plugin.json`, `.cursor-plugin/plugin.json`, and `.opencode/plugins/grimoire.js` discover skills from `./skills/` recursively — no update needed when adding domains.
