<div align="center">
  <a href="https://github.com/jeffreytse/grimoire">
    <img alt="grimoire" src="./assets/banner.svg" width="700">
  </a>

  <p>The world's knowledge is in your AI. The world's practice is not.<br>A modern skills manager for AI agents — install, manage, and enforce best practices the way npm manages packages.</p>

<br><h1>📖 Grimoire 📖</h1>

</div>

<h4 align="center">
  The skills package manager for AI agents — with a built-in practices linter.<br>
  Declare the practices you require. Install them in one command. Enforce them automatically,<br>
  in CI, on every save, across every domain.
</h4>

<p align="center">
  <a href="https://github.com/sponsors/jeffreytse">
    <img src="https://img.shields.io/static/v1?label=sponsor&message=%E2%9D%A4&logo=GitHub&link=&color=greygreen"
      alt="Donate (GitHub Sponsor)" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire/releases">
    <img src="https://img.shields.io/github/v/release/jeffreytse/grimoire?color=brightgreen"
      alt="Release Version" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire/graphs/contributors">
    <img src="https://img.shields.io/github/contributors/jeffreytse/grimoire?color=brightgreen"
      alt="Contributors" />
  </a>

  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-greygreen.svg"
      alt="License: MIT" />
  </a>

  <a href="https://liberapay.com/jeffreytse">
    <img src="http://img.shields.io/liberapay/goal/jeffreytse.svg?logo=liberapay"
      alt="Donate (Liberapay)" />
  </a>

  <a href="https://patreon.com/jeffreytse">
    <img src="https://img.shields.io/badge/support-patreon-F96854.svg?style=flat-square"
      alt="Donate (Patreon)" />
  </a>

  <a href="https://ko-fi.com/jeffreytse">
    <img height="20" src="https://ko-fi.com/img/githubbutton_sm.svg"
      alt="Donate (Ko-fi)" />
  </a>

  <a href="#-agent-support">
    <img src="https://img.shields.io/badge/works%20with-Claude%20%C2%B7%20Codex%20%C2%B7%20Cursor%20%C2%B7%20Gemini%20%C2%B7%20OpenCode%20%C2%B7%20OpenClaw-blue"
      alt="Works with" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire-core">
    <img src="https://img.shields.io/badge/skills-1000%2B-blue"
      alt="1000+ Skills" />
  </a>
</p>

<div align="center">
  <h4>
    <a href="#-why-grimoire">Why</a> |
    <a href="#-skills-in-action">Features</a> |
    <a href="#%EF%B8%8F-install">Install</a> |
    <a href="#-quick-start">Quick Start</a> |
    <a href="#%EF%B8%8F-domains">Domains</a> |
    <a href="#-contributing">Contributing</a> |
    <a href="https://github.com/jeffreytse/grimoire/releases">Changelog</a> |
    <a href="#-license">License</a>
  </h4>
</div>

<div align="center">
  <sub>Built with ❤︎ by
  <a href="https://jeffreytse.net">jeffreytse</a> and
  <a href="https://github.com/jeffreytse/grimoire/graphs/contributors">contributors</a>
  </sub>
</div>
<br>

## 🎬 Demo

> "We're launching our SaaS in 48 hours. I'm terrified something will break. What do we do?"

![grimoire demo — pre-launch protocol: apply-premortem → design-slo → plan-incident-response → run-game-day](./assets/demo.gif)

Not an engineer? See [Skills in Action](#-skills-in-action) for a sample, or browse [all 27 domains in grimoire-core](https://github.com/jeffreytse/grimoire-core#-featured-skills).

## 🤔 Why Grimoire?

> Your AI knows everything — Grimoire makes it practice it.

Books gave everyone knowledge. Google gave everyone access. AI gave everyone comprehension. None of them gave everyone _practice_. Grimoire does.

The model knows SOLID, DDD, Google SRE, and the ABA Model Rules. Without explicit guidance, it enforces none of them. Grimoire is the enforcement layer — **skills**: named, citable, executable units of expert practice, one concept, one source, one set of steps.

- 🔍 **You don't know what you don't know.** Grimoire surfaces the governing standard you didn't know applied.
- 🤖 **Knowing is not enough. Grimoire enforces.** Specific steps, verifiable criteria, repeatable results. Not summaries — verdicts.
- 🌍 **The world's best practices belong to everyone.** McKinsey charges $1M. Senior lawyers bill $800/hr. The practices they follow are not proprietary. Grimoire makes them free — as executable steps, cited and verified.
- 🧱 **Every profession. Every AI.** 1000+ skills across 27 domains. Works with Claude, Codex, Cursor, Gemini, OpenCode, and OpenClaw.
- 📦 **Package-managed.** Declare skills in `grimoire.toml`, lock versions in `grimoire.lock`. Reproducible skill sets across machines and teams — like Cargo or npm for best practices.
- 🌐 **Open ecosystem.** Any git repo is a grimoire package. We encourage community packages to follow the [grimoire skill standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md) — but it's not required. Pair `grimoire-core` with company-internal skills, community packages, or your own library; all declared in one `grimoire.toml`.
- 🔬 **Semantic compliance, not syntax checking.** Other tools lint whether your `CLAUDE.md` is valid TOML. Grimoire checks whether your project actually follows the practices it declared — criteria by criteria, domain by domain. `grimoire check` is ESLint for practices. Gate CI with exit codes. Watch for changes with `--live`.

If you've spent 10,000 hours mastering something, your practice belongs here.

→ [Read the full philosophy](https://github.com/jeffreytse/grimoire-core/blob/main/PHILOSOPHY.md)

## 📦 Skills as Packages

Every AI coding tool ships its own config format — `CLAUDE.md`, `.cursorrules`, `AGENTS.md`. Developers copy-paste and drift. Grimoire solves this with a package manager model: declare the practices you need in `grimoire.toml`, pin versions in `grimoire.lock`, install to every agent in one command.

Any git repo is a valid package. grimoire-core is the official one — curated, cited, peer-reviewed. Tools like Tessl distribute skills broadly. Grimoire adds the enforcement layer: curated expert practices, version-locked, with `grimoire check` to close the loop from declaration to compliance.

```toml
# grimoire.toml — commit this to your repo
[package]
name    = "my-project"
version = "0.1.0"

[dependencies]
"jeffreytse/grimoire-core"             = "*"   # official package — all 1000+ skills
"jeffreytse/grimoire-core:engineering" = "*"   # one domain only
"jeffreytse/grimoire-core@0.1.0"       = "*"   # pinned to a specific release
"mycompany/internal-skills"            = "*"   # private company package
```

```bash
grimoire install jeffreytse/grimoire-core                      # all skills from official package
grimoire install "jeffreytse/grimoire-core:engineering"        # engineering domain only
grimoire install "jeffreytse/grimoire-core@0.1.0"              # pinned to a specific release
grimoire install "jeffreytse/grimoire-core@0.1.0:health"       # pinned release + health domain
grimoire install "gitlab.com/mycompany/internal-skills"        # any git repo as a package
grimoire install                                               # (re-)install all from grimoire.toml
grimoire update                                                # update all packages to latest
grimoire list                                                  # show installed packages and skills
grimoire uninstall mycompany/internal-skills                   # remove a package
```

**Ref syntax:** `owner/repo[@tag][:path]` — the `:path` suffix is a Standard Glob pattern (doublestar) matched against each skill's domain path. `**` matches any depth; `*` matches within one segment. Examples: `engineering/**` (all engineering skills), `health/sleep/**` (one subdomain only), `**/development/**` (development subdomain in any domain).

| Package                                                        | Type                | Quality gate                     |
| -------------------------------------------------------------- | ------------------- | -------------------------------- |
| [`grimoire-core`](https://github.com/jeffreytse/grimoire-core) | Official            | STANDARD.md peer review required |
| Any git repo                                                   | Community / private | Owner's discretion               |

grimoire-core is the official package — curated, cited, and reviewed against STANDARD.md. But grimoire manages _any_ package without restriction. Mix and match.

**Community packages** install identically to grimoire-core. We encourage following the [grimoire skill standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md) for quality and interoperability — but any git repo works. Publish as `grimoire-<name>` for discoverability:

| Package example | Focus |
|---|---|
| `yourorg/grimoire-fintech` | fintech-specific practices |
| `yourorg/grimoire-medical` | clinical protocols, FDA workflows |
| `yourorg/grimoire-legal-us` | US jurisdiction guides |
| `yourorg/internal-skills`  | your team's private practices |

The practices in your installed packages define what `grimoire check` validates against — declare the package, get the linter rules for free.

## 🔬 Linting for Practices

Declare which practices your project follows. `grimoire check` enforces them — CI gate, watch mode, or on-demand.

```bash
grimoire check                              # run compliance check against declared practices
grimoire check --live                       # watch mode — re-check on every file save
grimoire check --live --port 8080           # custom port (default 7890)
grimoire check --live --host 127.0.0.1      # localhost-only (default: all interfaces)
grimoire check --junit report.xml           # JUnit XML output for CI systems
grimoire check --scope changed              # check only changed files (incremental)
```

Sample output:

```
  ✓  propose-conventional-commit   100%  all 4 criteria passing
  ✓  apply-solid                    88%  7/8 criteria passing
  ✗  review-pull-request            62%  below threshold (80%)
     └─ missing: security-checklist, breaking-change-annotation

2 passed · 1 failed · exit 1
```

Gate any PR: `grimoire check --junit report.xml` exits non-zero when coverage drops below the configured threshold. No other tool checks whether your project _follows_ its declared practices — only whether the config file is valid.

## ✨ Skills in Action

**Engineering** — propose a commit message from staged diff:

```
User: /propose-conventional-commit

Claude: Inspecting staged changes...

  feat(auth): add OAuth2 PKCE flow for mobile clients

  Replaces implicit grant with PKCE to address RFC 6749 §10.16
  implicit flow vulnerabilities in native apps.

  Approve? [y/n]
```

Every skill is structured: triggering condition, verified source, numbered steps, no freeform guessing.

**Or describe the situation in plain language** — you don't need to know which skill applies:

```
User: My team keeps missing sprint goals and I don't know why.

Claude: Situation matches: plan-retrospective (engineering/project-management)

        Run a blameless retrospective (Amazon format):
        1. Timeline reconstruction — what actually happened
        2. Five whys — surface root causes, not symptoms
        3. Action items with owners and deadlines
        → /plan-retrospective to run the session now
```

`suggest-best-practice` auto-classifies any situation, routes to the matching skill,
or tells you exactly what to install if the skill isn't in your library yet.

→ [See skills across all 27 domains — law, finance, health, cooking, and more](https://github.com/jeffreytse/grimoire-core)

## ⚒️ Install

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.ps1 | iex
```

**Go:**

```bash
go install github.com/jeffreytse/grimoire@latest
```

**Then run the interactive wizard (all platforms):**

```bash
grimoire wizard    # guided setup: registry, domains, agents — no config editing required
```

Or set up manually:

```bash
grimoire update                              # fetch the official skill library
grimoire install                             # install to all detected AI agents
```

Auto-detects Claude Code, Codex, Gemini CLI, OpenClaw, and OpenCode.

**Native plugin shortcuts (Claude Code):**

```bash
# Step 1: add the marketplace
/plugin marketplace add jeffreytse/grimoire-core

# Step 2: install (skills are namespaced, e.g. /grimoire-engineering:propose-conventional-commit)
/plugin install grimoire@grimoire-core                   # all domains (latest)
/plugin install grimoire-engineering@grimoire-core       # one domain

# For subdomain-level installs, use grimoire
```

**Package management:**

```bash
grimoire install jeffreytse/grimoire-core   # add + install official package
grimoire install myorg/my-skills            # add + install any git repo as a package
grimoire uninstall myorg/my-skills          # remove + unlink a package
grimoire update                             # update all packages to latest
grimoire list                               # list installed packages and skill counts
```

**Granular installs (within a package):**

```bash
grimoire install                                            # install all packages
grimoire install --domain engineering
grimoire install --domain engineering --subdomain development
grimoire install --skill engineering/development/propose-conventional-commit
grimoire install --target all      # install to all agents, even if not detected
grimoire doctor                    # health check: git repo, symlinks, config
grimoire version                   # version info with commit and date
```

**Gemini CLI:**

```bash
gemini extensions install https://github.com/jeffreytse/grimoire-core           # latest
gemini extensions install https://github.com/jeffreytse/grimoire-core@v1.0.0    # pin to a release
gemini extensions update grimoire                                               # update later
```

**Cursor:**

```bash
grimoire install --target cursor
```

**OpenCode:**

```bash
grimoire install --target opencode
```

Or via plugin in `opencode.json`:

```json
{ "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire-core"] }
```

**OpenClaw:** see [`.openclaw/INSTALL.md`](./.openclaw/INSTALL.md) or run `grimoire install --target openclaw`.

## 🤖 Agent Support

| Agent              | Plugin install                                                                                                 | Script install                       |
| ------------------ | -------------------------------------------------------------------------------------------------------------- | ------------------------------------ |
| Claude Code        | `/plugin marketplace add jeffreytse/grimoire-core` then `/plugin install grimoire@grimoire-core`               | `grimoire install --target claude`   |
| GitHub Copilot CLI | `copilot plugin marketplace add jeffreytse/grimoire-core` then `copilot plugin install grimoire@grimoire-core` | `grimoire install --target all`      |
| Gemini CLI         | `gemini extensions install https://github.com/jeffreytse/grimoire-core`                                        | `grimoire install --target gemini`   |
| OpenCode           | See [`.opencode/INSTALL.md`](./.opencode/INSTALL.md)                                                           | `grimoire install --target opencode` |
| OpenClaw           | See [`.openclaw/INSTALL.md`](./.openclaw/INSTALL.md)                                                           | `grimoire install --target openclaw` |
| Codex CLI          | `AGENTS.md` auto-loaded; browse `/plugins` in CLI                                                              | `grimoire install --target codex`    |
| Cursor             | `AGENTS.md` context injection                                                                                  | `grimoire install --target cursor`   |

## 🖥️ Editor Integration (LSP)

`grimoire lsp` implements the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/) over stdio. Any LSP-capable editor gets compliance diagnostics in the gutter in real time — no plugin required beyond a one-time language-client config.

**How it works:** on every file save the server runs `grimoire check` in the background and pushes findings as LSP diagnostics to your editor. Pass/hint items are suppressed; only errors, warnings, and info appear.

| Editor | Setup |
|--------|-------|
| Neovim | `lspconfig` custom server (see below) |
| VSCode | grimoire extension, or `tasks.json` with a custom server entry |
| Helix  | `languages.toml` custom language server entry |
| Any LSP client | Point `cmd` at `grimoire lsp` |

**Neovim (nvim-lspconfig):**

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

if not configs.grimoire then
  configs.grimoire = {
    default_config = {
      cmd = { 'grimoire', 'lsp' },
      filetypes = { 'go', 'python', 'javascript', 'typescript', 'rust', 'ruby' },
      root_dir = lspconfig.util.root_pattern('grimoire.toml', '.git'),
      single_file_support = true,
    },
  }
end

lspconfig.grimoire.setup {}
```

**Helix (`languages.toml`):**

```toml
[[language-server]]
name = "grimoire"
command = "grimoire"
args = ["lsp"]

[[language]]
name = "go"
language-servers = ["gopls", "grimoire"]
```

**VSCode (`settings.json` — requires grimoire extension or a custom extension):**

```json
{
  "grimoire.lsp.enable": true,
  "grimoire.lsp.command": "grimoire",
  "grimoire.lsp.args": ["lsp"]
}
```

## 🚀 Quick Start

**New to grimoire? Run the interactive wizard:**

```bash
grimoire wizard
```

The wizard walks you through registry selection, domain choice, and agent linking — no manual config editing. It detects your installed AI agents and installs the right skills automatically.

---

**After install, describe any problem in plain language:**

```
User: I need to raise a Series A but don't know how to pitch investors.

Claude: Situation matches: write-value-proposition + design-go-to-market + apply-pyramid-principle
        Applying suggest-best-practice...
        → Start with your value prop. /write-value-proposition
```

Or invoke a skill directly:

```
/suggest-best-practice     # describe any problem — auto-routes to the right skill
/review-pull-request       # engineering code review
/calculate-fire-number     # how much do I need to retire?
/review-saas-contract      # flag dangerous clauses before signing
/design-training-program   # build a training program
```

For CI enforcement, initialize once per project and gate PRs with the `grimoire` CLI:

```bash
grimoire init    # creates grimoire.toml with auto-detected profile
grimoire check   # exits 0 (pass) or 1 (fail)
```

**New to grimoire?** Start with `/suggest-best-practice`. Describe any professional or life situation — it reads your context and routes you to the matching skill, or tells you exactly what to install if the skill isn't in your library yet.

## 🎯 Workflows

| Your situation                                                                                                | Start here                                |
| ------------------------------------------------------------------------------------------------------------- | ----------------------------------------- |
| Know exactly which skill you need                                                                             | `/skill-name` directly                    |
| Have a problem, unsure which skill                                                                            | `/suggest-best-practice`                  |
| Already have a plan, want gaps checked                                                                        | `/review-best-practice-fit`               |
| Needs 2+ practices coordinated — within one domain or across many — and sub-problems are identifiable upfront | `/plan-best-practice-solution`            |
| Complex problem where sub-problems are opaque and emerge through execution                                    | `/apply-best-practice-tree`               |
| Don't know what practices exist for a topic                                                                   | `/discover-best-practices`                |
| About to start a task — want to catch gaps before you begin                                                   | `/start-best-practice`                    |
| Problem isn't clear yet — need to define it before solving                                                    | `/analyze-best-practice-problem`          |
| Activate a paradigm's best practices (OOP, TDD, etc.)                                                         | `/apply-best-practice-profile`            |
| Align any project or artifact to stated best practice preferences (BPDD)                                      | `/apply-best-practice-driven-development` |
| Check if any artifact aligns with stated best practice preferences                                            | `/check-best-practice-compliance`         |
| Have a specific compliance finding to fix                                                                     | `/fix-best-practice-finding`              |
| Two practices exist — want side-by-side comparison                                                            | `/compare-best-practices`                 |
| Two practices conflict — want to reason through which fits                                                    | `/resolve-best-practice-conflict`         |
| Resolved a conflict — want to save the decision for future sessions                                           | `/pin-best-practice-preference`           |

→ [BPDD guide](./docs/bpdd.md) — cycle, compliance linter, LSP output, suppression, CI integration

### Common paths

**Describe any problem → `suggest-best-practice`**

Describe any situation. `suggest-best-practice` reads it, finds the best-matching skill(s), and routes you there.

```
User: My senior engineer just quit and the team is losing confidence.

Claude: Situation matches: design-onboarding-program + write-leadership-principles
        → /design-onboarding-program to rebuild team structure
        → /write-leadership-principles to establish decision-making clarity
```

**Multiple practices needed → `plan-best-practice-solution`**

When a problem needs 2+ practices coordinated in sequence — within one domain or across many — `plan-best-practice-solution` decomposes it with MECE methodology, sequences skills by dependency, and recursively plans complex sub-problems.

```
User: I want to leave my job and launch a SaaS startup.

Claude: Situation spans 4 domains. Solution plan:
        1. design-business-model    — validate the idea before building
        2. calculate-unit-economics — confirm revenue model is viable
        3. review-saas-contract     — legal structure before hiring
        4. design-go-to-market      — growth strategy before launch

        Apply in order? I'll pause for confirmation after each step.
```

### Learning workflows

| Want to...                                                | Start here               |
| --------------------------------------------------------- | ------------------------ |
| Understand what a skill does and why it exists            | `/explain-best-practice` |
| Adapt a skill to a different context or constraint        | `/adapt-best-practice`   |
| Teach a practice to someone else (structured walkthrough) | `/teach-best-practice`   |


## ⚙️ Configuration

When top-tier orgs disagree (Google vs IBM, startup vs enterprise), pin which approach you prefer once — grimoire applies it automatically from then on.

```
User: I follow Google's engineering practices when there's a conflict.

Claude: Pinning preference: Google Engineering Practices → engineering domain.
        Save to: session / project / global?

User: global

Claude: Saved to ~/.config/grimoire/grimoire.toml.
        Future engineering conflicts resolve toward Google's approach automatically.
```

**Override priority** — grimoire checks preferences in this order, first match wins:

| Priority | Level   | Stored in                                  | Scope                                           |
| -------- | ------- | ------------------------------------------ | ----------------------------------------------- |
| 1st      | Session | In-memory                                  | Current session only — resets when session ends |
| 2nd      | Project | `grimoire.toml`                            | Current project, committed to repo              |
| 3rd      | Global  | `~/.config/grimoire/grimoire.toml`         | All projects on this machine                    |
| 4th      | System  | `/etc/grimoire/grimoire.toml`              | All users on this machine                       |

**Configure manually** — edit the settings files directly without going through the AI:

```toml
# ~/.config/grimoire/grimoire.toml  (global — applies everywhere)
profiles = ["oop"]                 # activate all skills tagged "oop"
# profiles = ["clean-architecture", "tdd"]  # multiple — first entry wins conflicts

[engineering]
practices = ["Google Engineering Practices"]

[finance]
practices = ["CFA Institute standards"]
```

```toml
# grimoire.toml  (project — overrides global for this repo)
[engineering.architecture]
practices = [
  "SOLID principles: production code",
  "KISS: prototypes, scripts"
]
fallback = "ask"
```

Project configuration overrides global. Session pins override both. Teams can share a global standard while individual projects deviate where needed.

**`practices = ["OOP"]` vs `profiles = ["oop"]`** — both signal OOP intent, but differently. `practices = ["OOP"]` in a domain section is a loose hint — the AI leans toward OOP conventions from its training. `profiles = ["oop"]` at the top level activates specific installed skills (exact steps, validated sources). Use `profiles` for precision; `practices` for domain-level style preference. → [Full comparison](./docs/profiles.md#profiles-vs-practices)

**Guided configuration:**

- `/configure-grimoire` — view, edit, or validate config without touching TOML directly
- `/apply-best-practice-profile` — activate a full paradigm (OOP, TDD, clean architecture) in one command
- `/resolve-best-practice-conflict` — resolve contradictions between two installed skills and record the priority automatically
- `/apply-best-practice-driven-development` — run the full BPDD cycle ([BPDD guide](./docs/bpdd.md))

→ [Full configuration reference](./docs/config.md) — all keys, override hierarchy, TOML examples

## 🎭 Profiles

Activate a named set of skills in one line — no list to maintain, no file to create.

```toml
# grimoire.toml
profiles = ["oop"]   # activates every installed skill tagged "oop"
```

Grimoire resolves the name in this order, first match wins:

1. `.grimoire/profiles/<name>.toml` — project-level file
2. `~/.grimoire/profiles/<name>.toml` — user-level file
3. `.grimoire/profiles/default.toml` — project-level fallback
4. `~/.grimoire/profiles/default.toml` — user-level fallback
5. `[profiles.<name>]` in `grimoire.toml` — inline definition (no separate file needed)
6. Tag query — all installed skills where `tags` contains the name

If no file or inline definition exists, the tag query fires automatically. `profiles = ["oop"]` works without creating any file.

**Multiple profiles** — combine paradigms; first entry wins conflicts, duplicates are deduplicated:

```toml
profiles = ["clean-architecture", "tdd"]  # clean-architecture wins if both include the same skill
```

**Custom profile** — curate a team-specific subset when the tag set is too broad:

```toml
# .grimoire/profiles/my-team.toml
name = "my-team"
description = "Our backend team's default practices"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"
```

Commit `.grimoire/profiles/` to share standards across the team. Publish as a gist or repo (`grimoire-profile-<name>`) for the community.

**Inline profile** — define directly in `grimoire.toml`, no separate file:

```toml
# grimoire.toml
profiles = ["my-team"]

[profiles.my-team]
description = "Our backend team's default practices"
extends = ["oop"]

[[profiles.my-team.skills]]
name = "apply-solid-principles"

[[profiles.my-team.skills]]
name = "apply-domain-driven-design"

exclude = ["apply-law-of-demeter"]
```

If a `profiles/my-team.toml` file also exists for the same name, the file wins.

**`profiles` vs `practices`** — `profiles` activates skill bundles globally; `practices` is a domain-scoped explicit list. → [Full comparison](./docs/profiles.md#profiles-vs-practices)

→ [Full profiles guide](./docs/profiles.md) — resolution order, conflict handling, sharing profiles

## 📏 BPDD — Best Practice Driven Development

grimoire is a linter for best practices — same model as ESLint for code style:

```
# Code style linter:             │  # Best-practice linter:
npm install eslint               │  grimoire install jeffreytse/grimoire-core
eslint .                         │  grimoire check
# watch mode:                    │  # watch mode:
eslint --watch                   │  grimoire check --live
```

Declare which practices you require in `grimoire.toml`, install the packages that encode them, then run `grimoire check` against any artifact — a codebase, a legal contract, a business plan, a training program. Same criteria every run. Gaps that survive human review get caught by the check.

**The cycle** — same inversion as TDD: declare what "good" looks like first, then bring the artifact into alignment.

```
1. Red      — run compliance check; identify which practices FAIL or are PARTIAL
2. Green    — invoke the relevant grimoire skill; fix until the check passes
3. Refactor — clean up while keeping the check green
4. Commit   — record progress; repeat for next gap
```

Run `/apply-best-practice-driven-development` to drive the full cycle. Run `/check-best-practice-compliance` for a one-off check.

**Two modes** — choose based on how you want to run the AI:

- **Independent mode** (default) — grimoire runs the AI check itself; no prior AI session needed
- **Report mode** — an AI session skill (`/check-best-practice-compliance`) generates the JSON report; `grimoire check --from-report` reads it and enforces thresholds

```bash
grimoire init                        # one-time project setup
grimoire check                       # independent mode: auto-selects local CLI or API provider
grimoire check --via claude          # force a specific local agent
grimoire check --live                       # file watcher + browser report — like eslint --watch
grimoire check --live --port 8080           # custom port (default 7890)
grimoire check --live --host 127.0.0.1     # localhost-only (default: all interfaces)
grimoire check --ci                  # + GitHub Actions annotations
grimoire check --junit report.xml    # + JUnit XML for CI reporters
grimoire check --from-report         # report mode: reads .grimoire/reports/compliance-latest.json
```

**Output** — always written to `.grimoire/reports/`:

| File                          | Format              | Use                                            |
| ----------------------------- | ------------------- | ---------------------------------------------- |
| `compliance-<timestamp>.json` | LSP-compatible JSON | editors, CI pipelines, LSP servers, dashboards |
| `compliance-latest.json`      | JSON (symlink)      | Always points to most recent run — use in CI   |
| `compliance-<timestamp>.html` | HTML                | browser or CI artifact upload                  |

The JSON follows the LSP Diagnostic schema — `uri` + `range` locate any finding in any text artifact, not just code.

**Coverage thresholds** — set in `grimoire.toml`, enforced on every check:

```toml
[standards.engineering]
compliance-threshold = 80        # fail if overall criteria coverage < 80%
compliance-threshold-error = 0   # fail if any error-severity violations remain
```

Use `/fix-best-practice-finding` to fix one specific compliance finding — targeted, location-aware, verified. Use `/apply-best-practice-driven-development` to fix everything systematically.

→ [Full BPDD guide](./docs/bpdd.md) — cycle, linter, LSP schema, false positive suppression, incremental mode

## 🌟 Featured Skills

| Skill                                                                                                                                                                  | Domain                   | Source methodology                       | Verified |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ | ---------------------------------------- | -------- |
| [`review-saas-contract`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/contracts/skills/review-saas-contract/)                                      | law/contracts            | ABA model SaaS agreements                | ✓        |
| [`calculate-fire-number`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/personal-finance/skills/calculate-fire-number/)                         | finance/personal-finance | Bengen (1994) / Trinity Study            | ✓        |
| [`negotiate-salary`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/personal-finance/skills/negotiate-salary/)                                   | finance/personal-finance | Fisher & Ury "Getting to Yes" / BLS data | ✓        |
| [`design-sleep-protocol`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/sleep/skills/design-sleep-protocol/)                                     | health/sleep             | Matthew Walker "Why We Sleep" / AASM     | ✓        |
| [`apply-mise-en-place`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/techniques/skills/apply-mise-en-place/)                                   | cooking/techniques       | Culinary Institute of America            | ✓        |
| [`apply-five-whys`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/reliability/skills/apply-five-whys/)                                      | engineering/reliability  | Toyota Production System / Google SRE    | ✓        |
| [`design-training-program`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/training/skills/design-training-program/)                               | sports/training          | NSCA CSCS curriculum                     | ✓        |
| [`apply-acceptance-commitment-therapy`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/cognitive/skills/apply-acceptance-commitment-therapy/) | psychology/cognitive     | Hayes / ACBS meta-analyses               | ✓        |
| [`write-value-proposition`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/copywriting/skills/write-value-proposition/)                          | writing/copywriting      | Osterwalder "Value Proposition Design"   | ✓        |
| [`design-training-periodization-plan`](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/training/skills/design-training-periodization-plan/)        | sports/training          | Bompa "Periodization" / NSCA             | ✓        |

→ [Browse all skills by domain](https://github.com/jeffreytse/grimoire-core/blob/main/SKILLS.md)

## 📐 The Grimoire Skill Standard

grimoire-core maintains an open standard for AI agent skill quality — freely adoptable by any skill library. Community packages are encouraged to follow it; grimoire installs any git repo regardless. This quality gate applies to grimoire-core contributions only.

Every skill must pass `review-best-practice-skill` before merge:

| Criterion      | Requirement                           | Rejection example                |
| -------------- | ------------------------------------- | -------------------------------- |
| **Adopted by** | Named organizations or institutions   | "Many top companies"             |
| **Impact**     | Cited study or % number               | "Significantly improves quality" |
| **Steps**      | Immediately executable                | Abstract theory or advice        |
| **Scope**      | One concept per skill                 | "Nutrition and training program" |
| **Source**     | External institution or standard body | Internal opinion                 |

→ [Read the full standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md) · [Adopt this standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md#adopting-this-standard)

See [CONTRIBUTING.md](https://github.com/jeffreytse/grimoire-core/blob/main/CONTRIBUTING.md) to submit a skill.

## 🗺️ Domains

Domains below are from **grimoire-core**, the official package. Community packages may cover any domain structure — grimoire imposes no layout.

📦 **Official package:** [github.com/jeffreytse/grimoire-core](https://github.com/jeffreytse/grimoire-core)

| Domain                                                                                     | Sub-domains                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [grimoire](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/)             | **Setup:** [install-grimoire](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/install-grimoire/) · [configure-grimoire](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/configure-grimoire/) · **Problem analysis:** [analyze-best-practice-problem](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/analyze-best-practice-problem/) · [discover-best-practices](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/discover-best-practices/) · **Routing:** [suggest-best-practice](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/suggest-best-practice/) · [start-best-practice](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/start-best-practice/) · **Solution planning:** [plan-best-practice-solution](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/plan-best-practice-solution/) · [apply-best-practice-tree](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/apply-best-practice-tree/) · **Practice evaluation:** [review-best-practice-fit](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/review-best-practice-fit/) · [compare-best-practices](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/compare-best-practices/) · [audit-applied-best-practices](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/audit-applied-best-practices/) · **Practice understanding:** [explain-best-practice](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/explain-best-practice/) · [adapt-best-practice](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/adapt-best-practice/) · [teach-best-practice](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/teach-best-practice/) · **Preferences:** [pin-best-practice-preference](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/pin-best-practice-preference/) · [resolve-best-practice-conflict](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/resolve-best-practice-conflict/) · [apply-best-practice-profile](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/apply-best-practice-profile/) · [write-best-practice-profile](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/write-best-practice-profile/) · [review-best-practice-profile](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/review-best-practice-profile/) · [share-best-practice-profile](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/share-best-practice-profile/) · **Compliance:** [apply-best-practice-driven-development](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/apply-best-practice-driven-development/) · [check-best-practice-compliance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/check-best-practice-compliance/) · **Contributors:** [write-best-practice-skill](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/write-best-practice-skill/) · [review-best-practice-skill](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/review-best-practice-skill/) · [revise-best-practice-skill](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/revise-best-practice-skill/) · [audit-best-practice-domain](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/audit-best-practice-domain/) · [deprecate-best-practice-skill](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/deprecate-best-practice-skill/) · [design-best-practice-domain](https://github.com/jeffreytse/grimoire-core/tree/main/skills/meta/skills/design-best-practice-domain/) |
| [engineering](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/)   | [development](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/development/skills/), [frontend](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/frontend/skills/), [architecture](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/architecture/skills/), [testing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/testing/skills/), [reliability](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/reliability/skills/), [devops](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/devops/skills/), [cloud](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/cloud/skills/), [networking](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/networking/skills/), [security](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/security/skills/), [data](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/data/skills/), [ai](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/ai/skills/), [hardware](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/hardware/skills/), [mobile](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/mobile/skills/), [performance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/performance/skills/), [project-management](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/project-management/skills/), [product](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/product/skills/), [documentation](https://github.com/jeffreytse/grimoire-core/tree/main/skills/engineering/documentation/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| [writing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/)           | [creative](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/creative/skills/), [technical](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/technical/skills/), [copywriting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/copywriting/skills/), [academic](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/academic/skills/), [journalism](https://github.com/jeffreytse/grimoire-core/tree/main/skills/writing/journalism/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| [design](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/)             | [ui-ux](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/ui-ux/skills/), [graphic](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/graphic/skills/), [branding](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/branding/skills/), [motion](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/motion/skills/), [product](https://github.com/jeffreytse/grimoire-core/tree/main/skills/design/product/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| [business](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/)         | [strategy](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/strategy/skills/), [operations](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/operations/skills/), [leadership](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/leadership/skills/), [entrepreneurship](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/entrepreneurship/skills/), [hr](https://github.com/jeffreytse/grimoire-core/tree/main/skills/business/hr/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| [science](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/)           | [biology](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/biology/skills/), [physics](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/physics/skills/), [chemistry](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/chemistry/skills/), [mathematics](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/mathematics/skills/), [earth-science](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/earth-science/skills/), [astronomy](https://github.com/jeffreytse/grimoire-core/tree/main/skills/science/astronomy/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| [marketing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/)       | [seo](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/seo/skills/), [content](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/content/skills/), [social-media](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/social-media/skills/), [paid-ads](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/paid-ads/skills/), [growth](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/growth/skills/), [analytics](https://github.com/jeffreytse/grimoire-core/tree/main/skills/marketing/analytics/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| [health](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/)             | [fitness](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/fitness/skills/), [nutrition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/nutrition/skills/), [mental-health](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/mental-health/skills/), [sleep](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/sleep/skills/), [medicine](https://github.com/jeffreytse/grimoire-core/tree/main/skills/health/medicine/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| [finance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/)           | [personal-finance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/personal-finance/skills/), [investing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/investing/skills/), [accounting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/accounting/skills/), [real-estate](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/real-estate/skills/), [corporate](https://github.com/jeffreytse/grimoire-core/tree/main/skills/finance/corporate/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| [education](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/)       | [curriculum](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/curriculum/skills/), [teaching](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/teaching/skills/), [e-learning](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/e-learning/skills/), [assessment](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/assessment/skills/), [learning-science](https://github.com/jeffreytse/grimoire-core/tree/main/skills/education/learning-science/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| [film](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/)                 | [cinematography](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/cinematography/skills/), [directing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/directing/skills/), [editing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/editing/skills/), [screenwriting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/screenwriting/skills/), [production](https://github.com/jeffreytse/grimoire-core/tree/main/skills/film/production/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| [law](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/)                   | [contracts](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/contracts/skills/), [ip](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/ip/skills/), [employment](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/employment/skills/), [privacy](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/privacy/skills/), [corporate](https://github.com/jeffreytse/grimoire-core/tree/main/skills/law/corporate/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| [photography](https://github.com/jeffreytse/grimoire-core/tree/main/skills/photography/)   | [composition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/photography/composition/skills/), [lighting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/photography/lighting/skills/), [editing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/photography/editing/skills/), [genres](https://github.com/jeffreytse/grimoire-core/tree/main/skills/photography/genres/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| [music](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/)               | [composition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/composition/skills/), [production](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/production/skills/), [mixing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/mixing/skills/), [theory](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/theory/skills/), [performance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/music/performance/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| [cooking](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/)           | [techniques](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/techniques/skills/), [baking](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/baking/skills/), [flavor](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/flavor/skills/), [nutrition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/nutrition/skills/), [world-cuisine](https://github.com/jeffreytse/grimoire-core/tree/main/skills/cooking/world-cuisine/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| [language](https://github.com/jeffreytse/grimoire-core/tree/main/skills/language/)         | [learning](https://github.com/jeffreytse/grimoire-core/tree/main/skills/language/learning/skills/), [linguistics](https://github.com/jeffreytse/grimoire-core/tree/main/skills/language/linguistics/skills/), [translation](https://github.com/jeffreytse/grimoire-core/tree/main/skills/language/translation/skills/), [communication](https://github.com/jeffreytse/grimoire-core/tree/main/skills/language/communication/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| [art](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/)                   | [drawing](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/drawing/skills/), [painting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/painting/skills/), [digital-art](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/digital-art/skills/), [illustration](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/illustration/skills/), [color-theory](https://github.com/jeffreytse/grimoire-core/tree/main/skills/art/color-theory/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| [sports](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/)             | [training](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/training/skills/), [coaching](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/coaching/skills/), [nutrition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/nutrition/skills/), [tactics](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/tactics/skills/), [recovery](https://github.com/jeffreytse/grimoire-core/tree/main/skills/sports/recovery/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| [productivity](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/) | [time-management](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/time-management/skills/), [habits](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/habits/skills/), [focus](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/focus/skills/), [goals](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/goals/skills/), [tools](https://github.com/jeffreytse/grimoire-core/tree/main/skills/productivity/tools/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| [travel](https://github.com/jeffreytse/grimoire-core/tree/main/skills/travel/)             | [planning](https://github.com/jeffreytse/grimoire-core/tree/main/skills/travel/planning/skills/), [budgeting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/travel/budgeting/skills/), [cultural](https://github.com/jeffreytse/grimoire-core/tree/main/skills/travel/cultural/skills/), [adventure](https://github.com/jeffreytse/grimoire-core/tree/main/skills/travel/adventure/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| [psychology](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/)     | [cognitive](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/cognitive/skills/), [behavioral](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/behavioral/skills/), [social](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/social/skills/), [clinical](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/clinical/skills/), [positive](https://github.com/jeffreytse/grimoire-core/tree/main/skills/psychology/positive/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| [home](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/)                 | [renovation](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/renovation/skills/), [interior-design](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/interior-design/skills/), [gardening](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/gardening/skills/), [organization](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/organization/skills/), [smart-home](https://github.com/jeffreytse/grimoire-core/tree/main/skills/home/smart-home/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| [environment](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/)   | [sustainability](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/sustainability/skills/), [ecology](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/ecology/skills/), [climate](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/climate/skills/), [energy](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/energy/skills/), [policy](https://github.com/jeffreytse/grimoire-core/tree/main/skills/environment/policy/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| [pets](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/)                 | [dogs](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/dogs/skills/), [cats](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/cats/skills/), [training](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/training/skills/), [nutrition](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/nutrition/skills/), [health](https://github.com/jeffreytse/grimoire-core/tree/main/skills/pets/health/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| [fashion](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/)           | [styling](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/styling/skills/), [wardrobe](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/wardrobe/skills/), [design](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/design/skills/), [sustainability](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/sustainability/skills/), [accessories](https://github.com/jeffreytse/grimoire-core/tree/main/skills/fashion/accessories/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| [parenting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/parenting/)       | [infant](https://github.com/jeffreytse/grimoire-core/tree/main/skills/parenting/infant/skills/), [toddler](https://github.com/jeffreytse/grimoire-core/tree/main/skills/parenting/toddler/skills/), [school-age](https://github.com/jeffreytse/grimoire-core/tree/main/skills/parenting/school-age/skills/), [teen](https://github.com/jeffreytse/grimoire-core/tree/main/skills/parenting/teen/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| [automotive](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/)     | [maintenance](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/maintenance/skills/), [troubleshooting](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/troubleshooting/skills/), [buying](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/buying/skills/), [modifications](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/modifications/skills/), [ev](https://github.com/jeffreytse/grimoire-core/tree/main/skills/automotive/ev/skills/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |

## ❓ FAQ

**Isn't this already in the model's training data?**

Yes — and no. Models know _about_ best practices. Skills make models _do_ them, reliably.

The difference:

| Without a skill                            | With a skill                                                          |
| ------------------------------------------ | --------------------------------------------------------------------- |
| Model improvises a version of the practice | Model follows the exact steps from the source institution             |
| Output varies every run                    | Same process, same structure, every time                              |
| Practice applied only if you know to ask   | Skill triggers automatically when the situation matches               |
| Generic advice                             | Specific: the right gate, the right question, the right output format |

For simple tasks (write a test, fix a bug), the skeptic is right — the model doesn't need a skill. For complex, multi-step workflows — an SLO design, a post-mortem, an incident response — skills measurably change what you get. The model knows Google's engineering review process exists. It does not reliably know which question to ask first, what the output format is, or when to stop. That's what a skill encodes.

**These are just textbook practices the model already knows. Why bother?**

Knowing a practice and reliably executing it are different things. Ask any model "I just had a production incident" — you'll get a generic write-up. Run `write-post-mortem` and you get: blameless framing, 5-whys, timeline, contributing factors, action items with owners, and a detection section. The model _knew_ all of that before the skill existed. The skill is what makes it happen consistently, in the right format, every time.

The "textbook" objection gets it backwards. Established practices are _ideal_ for skills precisely because they're falsifiable — you can verify whether the output matches what Google's SRE book, Amazon's mechanisms, or the WHO protocol actually prescribes. If you find a skill that adds nothing over a bare prompt, that's a quality failure. [File an issue.](https://github.com/jeffreytse/grimoire/issues)

**Some frameworks here are universally known — isn't the value already in the model?**

The question isn't whether the model knows the framework name. It's whether the skill encodes what practitioners who already know the name still get wrong.

A framework qualifies when the skill has substantial content beyond the acronym or label — the step most people skip, the failure mode they don't avoid, the discipline that separates expert application from surface-level application. SWOT, for example, is universally known, but most practitioners stop at the 4-quadrant list and never derive the TOWS cross-matrix (SO/WO/ST/WT strategies) — the step that converts a diagnosis into actionable options. The skill encodes that gap.

A framework doesn't qualify when the full implementation reduces to restating the framework name. If a skill's entire content would be "follow the acronym," it adds nothing — the model already knows the letters.

The test: _"What would this skill contain beyond the framework name?"_ If the answer is "the step practitioners skip + the failure mode they don't avoid" — it qualifies. If the answer is "the letters, explained" — it doesn't.

**Does grimoire conflict with my team's existing conventions?**

Skills describe what the world's top institutions do. Your team may do things differently — and be right to. Two ways to handle it:

**Pin your preference.** Tell grimoire which approach to follow when practices conflict:

```
User: We follow Google's engineering practices, not IBM's.
→ Claude pins this via `pin-best-practice-preference` — applies automatically from now on.
```

**Override or fork.** A skill is a starting point, not a mandate. Adapt any skill to your context, or ignore it entirely. The format is plain Markdown and the license is MIT.

## 🤝 Contributing

Two ways to contribute:

**Add a skill** — share expert knowledge across law, finance, engineering, cooking, or any of the 27 domains. First skill takes ~30 minutes.

→ [Contributing guide](https://github.com/jeffreytse/grimoire-core/blob/main/CONTRIBUTING.md) in grimoire-core

**Improve the CLI** — bug fixes, new commands, performance, documentation.

1. Fork and clone this repo
2. `make build` to compile, `make test` to run tests
3. Open a PR — describe the problem it solves

→ [Open issues](https://github.com/jeffreytse/grimoire/issues) · [CONTRIBUTING.md](./CONTRIBUTING.md)

## ❤️ Support

grimoire is free. It replaces $500/hr lawyers, $300 doctor visits, and $1M McKinsey
engagements — at zero cost, forever.

If it saved you time, money, or a bad decision:

- **[⭐ Star this repo](https://github.com/jeffreytse/grimoire)** — takes 2 seconds, helps thousands of people find it
- **[💖 Sponsor on GitHub](https://github.com/sponsors/jeffreytse)** — keeps the maintainer funded to add more skills across more domains
- **[☕ Ko-fi](https://ko-fi.com/jeffreytse)** · **[Patreon](https://patreon.com/jeffreytse)** · **[Liberapay](https://liberapay.com/jeffreytse)** — one-time or recurring

Every star makes grimoire more visible. Every sponsorship funds one more domain.

## 📄 License

This project is licensed under the [MIT license](https://opensource.org/licenses/mit-license.php) © Jeffrey Tse.
