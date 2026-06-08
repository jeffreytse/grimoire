# Agent Integration Guide

grimoire works across six AI agents. Each has its own plugin format and installation method. This guide covers installation, usage, and tool mapping for each.

---

## Supported agents

| Agent | Plugin file | Install method |
|-------|------------|----------------|
| Claude Code | `.claude-plugin/plugin.json` | `/plugin marketplace add jeffreytse/grimoire` then `/plugin install grimoire@grimoire` |
| GitHub Copilot CLI | `.claude-plugin/plugin.json` (shared) | `copilot plugin marketplace add jeffreytse/grimoire` then `copilot plugin install grimoire@grimoire` |
| Gemini CLI | `gemini-extension.json` + `GEMINI.md` | `gemini extensions install https://github.com/jeffreytse/grimoire` |
| OpenCode | `.opencode/plugins/grimoire.js` | `.opencode/plugins/` auto-load or `opencode.json` plugin array |
| Codex CLI | `AGENTS.md` (auto-loaded) | `install.sh --target codex` for skills |
| Cursor | `AGENTS.md` (context injection) | `install.sh --target cursor` |

---

## Claude Code

**Install:**
```
/plugin marketplace add jeffreytse/grimoire
/plugin install grimoire@grimoire
```

**Plugin file:** `.claude-plugin/plugin.json`

Loads all skill paths listed in the `skills` array. Skills are namespaced by plugin name, e.g. `/grimoire-engineering:propose-conventional-commit`.

**Using skills:**
```
/suggest-best-practice
/write-best-practice-skill
/review-best-practice-skill
```

Invoke any meta skill by name with a leading slash. For domain skills, describe your situation and `suggest-best-practice` routes automatically.

**Tool names:** Use native Claude Code tool names (Read, Write, Edit, Bash, Skill, TodoWrite).

---

## GitHub Copilot CLI

**Install:**
```
copilot plugin marketplace add jeffreytse/grimoire
copilot plugin install grimoire@grimoire
```

**Plugin files:** `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json` (Copilot CLI checks `.claude-plugin/` automatically — shared with Claude Code)

Skills are namespaced by plugin name, e.g. `grimoire-engineering:propose-conventional-commit`.

**Using skills:**
Describe your task. Copilot CLI activates the matching skill based on its `description` triggering conditions.

---

## Codex CLI

**Install:**
```bash
./scripts/install.sh --target codex
```

Codex CLI supports `AGENTS.md` natively — grimoire's `AGENTS.md` is auto-loaded at session start, providing skill path conventions and domain routing. For the full skill library, run `install.sh` to copy skills to the Codex skills directory.

**Using skills:**
Codex CLI has an interactive plugin browser: run `codex` then type `/plugins` to browse and install marketplace plugins. For grimoire skills installed via `install.sh`, describe your task and Codex routes to the matching skill via `AGENTS.md` context.

---

## Cursor

**Install:**
```bash
./scripts/install.sh --target cursor
```

No official Cursor CLI plugin install command exists. Skills installed via `install.sh` are available through AGENTS.md context injection. Cursor reads `AGENTS.md` from the project root automatically in Agent mode.

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

> **Note (June 2026):** Gemini CLI is being replaced by Antigravity CLI on June 18, 2026 for unpaid/Google One users. Extension format may change — check the [migration announcement](https://developers.googleblog.com/an-important-update-transitioning-gemini-cli-to-antigravity-cli) for updates.

**Install:**
```
gemini extensions install https://github.com/jeffreytse/grimoire
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

## Other agents

Any agent that reads `AGENTS.md` from the project root gets grimoire's bootstrap context automatically — skill path convention, domain list, and routing instructions.

**Install:**
```bash
./scripts/install.sh --target all
```

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
| `.claude-plugin/plugin.json` | Add subdomain path to `skills` array |
| `.claude-plugin/marketplace.json` | Add domain and subdomain entries |
| `GEMINI.md` | Add domain row to domains table |
| `AGENTS.md` | Same update as GEMINI.md |

`.opencode/plugins/grimoire.js` discovers skills from `./skills/` recursively — no update needed when adding domains.
