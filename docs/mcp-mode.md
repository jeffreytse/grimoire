# MCP mode — user guide

How to set up grimoire as an MCP server for your AI agent, what it changes, what you'll see, and how to control it. The one-line setup is in the [README](../README.md#️-install).

← [Back to the README](../README.md)

## Set up

```bash
grimoire mcp install                         # Claude Code, for all your projects
grimoire mcp install --target all            # every supported agent found on your PATH
grimoire mcp install --target codex,gemini   # specific agents
grimoire mcp install --target claude-desktop # Claude Desktop (not included in --target all)
grimoire mcp install --scope project         # only this project (commit the files to share it with your team)
grimoire mcp uninstall --target all          # undo
```

Then start a new agent session. In Claude Code, `/mcp` lists grimoire as connected; each agent's own check is in the "Check" column below. Running install again is safe: nothing is duplicated.

For other MCP clients, `grimoire mcp config --target cursor|windsurf|cline` prints the config to paste in.

<a id="setup-per-agent"></a>

### What install changes, per agent

| Agent | Scopes | MCP server registered in | Instruction added to | Automatic routing | Check |
| --- | --- | --- | --- | --- | --- |
| Claude Code | user, project | `claude mcp add` | `~/.claude/CLAUDE.md` / `./CLAUDE.md` | ✓ hooks in `.claude/settings.json` | `/mcp` |
| Codex | user | `codex mcp add` | `~/.codex/AGENTS.md` | ✓ hooks in `~/.codex/hooks.json`; trust them once in `/hooks` | `codex mcp list` |
| Gemini CLI | user, project | `gemini mcp add` | `~/.gemini/GEMINI.md` / `./GEMINI.md` | ✓ hooks in `.gemini/settings.json` | `gemini mcp list` |
| GitHub Copilot CLI | user | `copilot mcp add` | `~/.copilot/copilot-instructions.md` | instruction only | `copilot mcp list` |
| OpenCode | user, project | `opencode.json` | `~/.config/opencode/AGENTS.md` / `./AGENTS.md` | ✓ grimoire plugin | `opencode mcp list` |
| Pi | user, project | `~/.pi/agent/mcp.json` / `.pi/mcp.json` | `~/.pi/agent/AGENTS.md` / `./AGENTS.md` | instruction only | `pi mcp list` |
| Kimi Code | user, project | `~/.kimi-code/mcp.json` / `.kimi-code/mcp.json` | `~/.kimi-code/AGENTS.md` / `./AGENTS.md` | instruction only | `/mcp` |
| OMP | user, project | `~/.omp/agent/mcp.json` / `.omp/mcp.json` | `~/.omp/agent/AGENTS.md` / `.omp/AGENTS.md` | instruction only | `/mcp` |
| Antigravity | user | `~/.gemini/config/mcp_config.json` | (added by its hooks) | ✓ hooks in `~/.gemini/config/hooks.json` | `/mcp` |
| Claude Desktop | user | `claude_desktop_config.json` | — | grimoire's server tells Claude when to use it | reopen the app; grimoire appears in its connectors |

**Automatic routing** means every request goes through grimoire's workflows without you asking. Where it says *instruction only*, the agent is told to use grimoire but decides for itself whether a request needs it, so routing is less consistent. Claude Desktop is written with the full path to grimoire, because desktop apps don't use your shell's PATH.

Install options:
- `--no-hook`: skip automatic routing; the instruction is still added.
- `--no-configure`: register the server only, with no instruction or hooks.
- `--no-allow`: don't pre-approve grimoire's read-only tools (see [Approvals](#approvals)).
- `--no-statusline` / `--statusline`: skip the Claude Code status bar badge, or add it after an existing status line.

Gemini keeps MCP servers disabled in folders you haven't trusted. Open Gemini in the folder and trust it, or install with `--scope user`.

## How it works in a session

With automatic routing on, every request goes through grimoire. It picks the workflow that fits: suggest one practice, plan several in order, break down an unclear problem, teach a practice, and so on. It asks a question when it isn't sure, then applies the specific practice. Small talk ("ok", "thanks") and slash commands are left alone. You don't need to know any skill names.

To run a specific skill, write `grimoire:<skill-name> [task]` in any agent, e.g. `grimoire:propose-conventional-commit`, or use `/grimoire:skill <skill-name>` in Claude Code. A mistyped name comes back with "did you mean" suggestions.

## What you see

grimoire shows its decisions so you can follow and redirect them. Each request gets one `grimoire ▸` line saying what it's doing and why. A question, a choice or "Apply this plan?" is always the last thing in a reply: nothing runs until you answer.

<pre>
❯ We're launching our SaaS in 48 hours. I'm terrified something will break.

grimoire ▸ planning several skills — launch readiness needs them in order

<b>Plan</b> — 2 steps
⎿  Step 1  apply-premortem — find the likeliest failures first
⎿  Step 2  plan-incident-response — decide who acts when one happens
Apply this plan? Steps that change files or run commands ask first.
</pre>

- **Choices:** when several skills fit, grimoire shows at most 4, numbered, with one `★ recommended`. Reply with a number, "all" or "skip".
- **Plans:** you confirm a plan once. A step where several skills fit shows lettered choices (`2a`, `2b`); reply "yes" to take every ★, or name another ("2b"). While it works, grimoire shows where it is (`— step 2 of 4`). It asks again only before a step that changes files, runs commands or acts outside the conversation, or after every step if you ask it to go step by step.
- **No match:** when no installed skill fits, grimoire says so, names what it looked for, answers from general knowledge, and tells you how to add a package that covers it.
- **Names:** skills are always shown by their id, the name you can type back as `grimoire:<id>`.

### How much it shows

Set the level with `grimoire config set core.trace <level>`, or `GRIMOIRE_TRACE=<level>` for one run:

| Level | Shows |
| --- | --- |
| `brief` (default) | the decisions: what grimoire is doing, the plan, each skill as it's applied, the outcome |
| `full` | also the evidence: which workflow, each skill's domain, the candidates and why one was chosen, risks with severity (🔴 high, 🟡 medium, ⚪ low), and the applied skill's source |
| `hidden` | no notes; grimoire still routes and applies practices silently (`off` also works) |

### What you can say

These work in any session, with no setup:

| Say | Effect |
| --- | --- |
| "grimoire help" | lists these controls |
| "show grimoire details" | switch to `full` |
| "quiet grimoire" | switch to `hidden`; routing continues |
| "plain grimoire" | words instead of symbols, for screen readers |
| "skip grimoire" | skip grimoire for this request |
| "grimoire off" / "grimoire on" | stop or resume routing in this session |

The first time grimoire acts after install, a one-time tip points you to "grimoire help".

**Accessibility:** `grimoire config set core.display-style plain` uses words instead of symbols everywhere (`Step 1:`, `Risk (high):`, `Done:`). It's on automatically when `NO_COLOR` is set; `GRIMOIRE_DISPLAY_STYLE=plain` sets it for one run, and "plain grimoire" for one session.

### Notes from grimoire

Some agents show short notes from grimoire itself: Claude Code, Codex and Gemini CLI (Codex shows them as warnings), and OpenCode as pop-ups. To avoid noise, grimoire says nothing when all is well at `brief`. It speaks up:
- **when something needs fixing**, e.g. `grimoire ▸ this project needs <package> — run grimoire update`
- **when you name a skill** (`grimoire ▸ loading skill <name>`, or "unknown skill")
- **at `full`**, also a "ready" note and the skills it matched by keyword

### Status bar badge (Claude Code)

Install adds a small ` GRIMOIRE ` badge to Claude Code's status bar, readable on light and dark themes. It shows the workflow and skill this session used last, and `[GRIMOIRE off]` while grimoire is off. With `NO_COLOR` it's plain text. If another tool already uses the status line, install leaves it alone and tells you; `--statusline` adds the badge after the other tool's output, and uninstall restores it.

## Turn grimoire on or off

Active (the default), grimoire routes every request. Off, it stays out of the way, but a skill you name with `grimoire:<name>` still loads, and nothing is uninstalled.

| Scope | Off | On |
| --- | --- | --- |
| All sessions | `grimoire disable` | `grimoire enable` |
| This session | say "grimoire off" | say "grimoire on" |
| One run (CI, scripts) | `GRIMOIRE_ACTIVE=0` | `GRIMOIRE_ACTIVE=1` |

`GRIMOIRE_ACTIVE` wins over the session setting, which wins over `grimoire enable`/`disable`. Turning it off in one session doesn't affect others, and the assistant confirms it in its reply. While off, a new session reminds you at most once a day. The in-session phrases need automatic routing; agents with *instruction only* follow `grimoire disable`.

Turning grimoire off is different from hiding its notes: `core.trace hidden` (or "quiet grimoire") keeps grimoire working but silent.

## Your project's grimoire.toml

If your project's `grimoire.toml` declares practices (`[dependencies]`, `[standards] profiles`, or `[standards.<domain>] practices`), grimoire **prefers** them; it doesn't restrict you to them:
- The whole library stays available for solving problems.
- Declared practices come first among comparable matches, and are marked with where they were declared: the project, your user config, or your organization's system config (the nearest wins).
- They get the ★ when they fit, and a recommended skill outside them is labeled "not a declared standard".
- BPDD and compliance checks treat only the declared practices as binding.

Teams that want a hard boundary set `[standards] scope = "strict"`. grimoire then offers only the declared practices; a skill you name with `grimoire:<name>` still loads, with a note.

grimoire reads `grimoire.toml` from the project your agent is working in. It never downloads anything on its own: run `grimoire update` (or `grimoire mcp install`) in the project to fetch packages it depends on. Until then, grimoire tells you which package is missing. More: [configuration](./config.md).

## Approvals

In Claude Code and Antigravity, install pre-approves grimoire's **read-only** tools (searching, loading and browsing skills, and reading context, status, config and profiles), so your assistant can look things up without asking you each time. Tools that change config or packages still ask. `--no-allow` turns this off, and uninstall removes exactly those approvals.

## MCP tools

What your assistant can use:

| Job | Tools |
| --- | --- |
| **Find and load skills** | `grimoire_search_skills` (keyword search across skill names, tags, descriptions and full text; several phrasings at once; filter by domain, subdomain or tags; declared practices first), `grimoire_get_skill` (one skill, or up to 20 by name), `grimoire_list_skills` (browse by domain) |
| **Context and checks** | `grimoire_context`, `grimoire_check` (compliance check), `grimoire_run_rules` (structural checks), `grimoire_get_config`, `grimoire_doctor`, `grimoire_status`, `grimoire_info`, `grimoire_version` |
| **Configuration** | `grimoire_config_get`, `grimoire_config_set`, `grimoire_config_unset`, `grimoire_config_add`, `grimoire_config_remove` |
| **Packages** | `grimoire_package_list`, `grimoire_package_add`, `grimoire_package_remove`, `grimoire_package_enable`, `grimoire_package_disable`, `grimoire_package_set`, `grimoire_package_reset`, `grimoire_package_update`, `grimoire_package_validate`, `grimoire_search` |
| **Profiles** | `grimoire_profile_list`, `grimoire_profile_show`, `grimoire_profile_init` |
| **Install and maintenance** | `grimoire_init`, `grimoire_install`, `grimoire_uninstall`, `grimoire_update`, `grimoire_clean`, `grimoire_self_update` |

The server also offers one prompt, `/grimoire:skill <name>`.

## Troubleshooting

| Symptom | Fix |
| --- | --- |
| The agent doesn't list grimoire | Start a new session after installing; run the agent's check command from the table above. |
| grimoire never routes in Codex | Trust grimoire's hooks once in Codex's `/hooks`. |
| Gemini shows grimoire as disabled | Trust the folder in Gemini, or install with `--scope user`. |
| `grimoire ▸ no skills package …` | Run `grimoire update` to download the skill library. |
| `grimoire ▸ this project needs <package>` | Run `grimoire update` in the project. |
| Anything else | `grimoire doctor` checks your setup and says what to fix. |
