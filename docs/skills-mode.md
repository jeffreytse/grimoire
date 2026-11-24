# Skills mode — linking skills into an agent

For agents without MCP support, or a small curated subset. Agents load only a limited listing of skill descriptions (about 8,000–16,000 characters in Claude Code and Codex), so for the whole library use [MCP mode](./mcp-mode.md).

Skills mode is fully supported, but new capabilities arrive in MCP mode first. Examples: project standards (preferred or strict), organization, user and project layers, and search-based skill selection.

In Claude Code, Codex and Gemini CLI, skills mode shows the same display as MCP mode — the `grimoire ▸` lines, the controls ("grimoire help", "show grimoire details", …) and the status bar badge. See [what you see](./mcp-mode.md#what-you-see). `grimoire install` accepts the same `--no-hook`, `--statusline` and `--no-statusline` options as `grimoire mcp install`.

← [Back to the README](../README.md)

## Install skills into the agent

**Skills mode — link skills into the agent** (agents without MCP, or a subset):

```bash
grimoire update                              # fetch the official skill library
grimoire install --domain engineering --subdomain testing   # link a small subset to all detected AI agents
grimoire install                             # link everything declared in grimoire.toml
```

Skills mode auto-detects every supported agent on your `PATH` (see [Agent Support](../README.md#-agent-support)). When the skills you link exceed an agent's listing budget, `grimoire install` says how many the agent will actually see and suggests MCP mode or a smaller `--domain`/`--subdomain`; `grimoire doctor` reports the same for existing installs.

**Package management:**

```bash
grimoire install jeffreytse/grimoire-core   # add + install official package
grimoire install myorg/my-skills            # add + install any git repo as a package
grimoire uninstall myorg/my-skills          # remove + unlink a package
grimoire update                             # update all packages to latest
grimoire update --force                     # update, discarding any local package modifications
grimoire list                               # list installed packages and skill counts
```

**Granular installs (within a package):**

```bash
grimoire install                                            # install all packages
grimoire install --domain engineering
grimoire install --domain engineering --subdomain development
grimoire install --skill engineering/development/propose-conventional-commit
grimoire install --target pi --global                       # one agent, user-wide
grimoire install --target all                               # all agents, even if not detected
grimoire doctor                    # health check: git repo, symlinks, config
grimoire version                   # version info with commit and date
```

## Native plugins and the grimoire CLI

**Recommended: MCP mode** (`grimoire mcp install`) installs nothing into the agent and serves the whole library on demand — see [Install](../README.md#️-install) for the agents it supports and why. The two ways below link skills into the agent (skills mode): its **native plugin** command, or the **grimoire CLI**, which works for every agent below. Agents load only a small listing of linked skills, so use them for agents without MCP support or for a small subset.

| Agent | Native plugin install | grimoire CLI |
| --- | --- | --- |
| Claude Code | `/plugin marketplace add jeffreytse/grimoire-core`<br>`/plugin install grimoire@grimoire` | `--target claude` |
| GitHub Copilot CLI | `copilot plugin marketplace add jeffreytse/grimoire-core`<br>`copilot plugin install grimoire@grimoire` | `--target copilot` |
| Codex | `codex plugin marketplace add jeffreytse/grimoire-core`<br>`codex plugin add grimoire@grimoire` | `--target codex` |
| OpenCode | `opencode plugin 'github:jeffreytse/grimoire-core#main'` | `--target opencode` |
| Pi | `pi install git:github.com/jeffreytse/grimoire-core` | `--target pi` |
| Kimi Code | `/plugins install https://github.com/jeffreytse/grimoire-core` | `--target kimi` |
| OMP (oh-my-pi) | `omp plugin marketplace add jeffreytse/grimoire-core`<br>`omp plugin install --scope user grimoire@grimoire` | `--target omp` |
| Gemini CLI | `gemini extensions install https://github.com/jeffreytse/grimoire-core` (context only) | `--target gemini` |
| Antigravity (agy) | not supported — use the CLI | `--target antigravity` |
| OpenClaw | see [`.openclaw/INSTALL.md`](https://github.com/jeffreytse/grimoire-core/blob/main/.openclaw/INSTALL.md) | `--target openclaw` |
| DeepSeek Harness | — | `--target dsh` |
| Grok Build | — | `--target grok --global` |
| Muse | — | `--target muse --global` |
| MiniMax Code | — | `--target minimax --global` |
| Any other agent | — | `--target agent` (`~/.agents/skills`) |

Every CLI row is `grimoire install --target <name>`. Add `--global` for a user-wide install; Grok, Muse, and MiniMax are global-only because they load skills from a plugin bundle, which grimoire builds and registers for you.

**Domain plugins:** native marketplaces also offer one plugin per domain, e.g. `grimoire-engineering@grimoire`. Skills are namespaced (`/grimoire-engineering:propose-conventional-commit`). For subdomain-level installs, use the grimoire CLI.

**Verification status:** native installs for Claude Code, Copilot, Codex, and OpenCode are tested end to end. Pi, Kimi, and OMP native installs follow each agent's documented format but are untested. Gemini's extension loads `GEMINI.md` context only; use the CLI for its skills. Antigravity's plugin installer ignores declared skill paths, so it cannot find grimoire's nested skills.
