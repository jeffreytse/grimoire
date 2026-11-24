# Configuration

`grimoire.toml` is where you pin preferences, activate profiles, and set compliance thresholds. Grimoire merges three file layers — local, global, and system — so teams can share a project baseline while individuals override with their own global config.

← [Back to the README](../README.md)

> Older versions used `settings.toml` (project: `.grimoire/settings.toml`, user: `~/.config/grimoire/settings.toml`). grimoire no longer reads those files: rename them to `grimoire.toml` at the project root and `~/.config/grimoire/grimoire.toml`. The format is the same.

---

## Override Hierarchy

Grimoire checks these sources in order. First match per key wins.

| Priority | Level | Source | Scope | Flag |
|----------|-------|--------|-------|------|
| 1st | Env vars | `GRIMOIRE_HOME`, `GRIMOIRE_SOURCE` | current process | — |
| 2nd | Local | `grimoire.toml` | this project, all team members (committed) | `--local` |
| 3rd | Global | `~/.config/grimoire/grimoire.toml` | all projects, this user | `--global` |
| 4th | System | `/etc/grimoire/grimoire.toml` (Linux/macOS) · `%PROGRAMDATA%\grimoire\grimoire.toml` (Windows) | all users on this machine | `--system` |

Each key resolves independently. Local can set `[standards.engineering].practices` while global sets `[standards].profiles` — both apply without conflict.

> **Note:** `[core]` keys (`core.home`) are user-level only. Grimoire strips them from project-level settings so teams cannot override machine-wide infrastructure config through a committed file.

---

## Environment Variables

Two env vars override package resolution — useful for CI/CD where writing config files is impractical:

| Variable | Description |
|----------|-------------|
| `GRIMOIRE_HOME` | Local path to the grimoire skills installation. Overrides the official package's home directory. |
| `GRIMOIRE_SOURCE` | Git URL for the official skills source. Overrides `[[package]]` and the built-in default. |

```bash
# CI: use a corporate mirror for skills
GRIMOIRE_SOURCE=https://internal.corp/grimoire-core.git grimoire update

# Point to a non-default install location
GRIMOIRE_HOME=/opt/grimoire grimoire doctor
```

When set, `grimoire config show` shows the source as `$GRIMOIRE_HOME` or `$GRIMOIRE_SOURCE` rather than a file path.

---

## Key Reference

| Key | Where | Type | Description |
|-----|-------|------|-------------|
| `home` | `[core]` | string | Override grimoire home path. User-level only. Prefer `GRIMOIRE_HOME` env var in CI. |
| `check-agents` | `[core]` | `string[]` | Local AI CLI resolution order for `grimoire check`. Defaults to `["claude","gemini","codex","copilot","opencode","openclaw"]`. Override to prefer a specific agent or limit which are tried. |
| `check-provider` | `[core.check-provider]` | section | API provider used when no local CLI is found. See [`[core.check-provider]`](#corecheck-provider) below. |
| `"<package-url>" = "<version>"` | `[dependencies]` | one key per dependency | Declares which skills to install and from where. Each key is a package URL (see [Package URL Format](#package-url-format)); the value is a version (`"*"` for any). Drives `grimoire install` with no args. The older `skills = [...]` list is not read. |
| `profiles` | `[standards]` | `string[]` | Activate skill bundles. Each entry is a profile name matched using depth-anywhere glob semantics (see [Profile Name Resolution](#profile-name-resolution)). A single name can expand to multiple profiles — `"engineering"` activates every `.toml` matching `**/engineering` in any installed package. |
| `practices` | `[standards.domain]` or `[standards.domain.subdomain]` | `string[]` | Ordered skill/practice hints for that domain — first entry wins conflicts |
| `disabled` | `[standards.domain]` or `[standards.domain.subdomain]` | `string[]` | Skills to suppress from the effective spec regardless of profiles — overrides profile inclusions |
| `fallback` | `[standards.domain]` or `[standards.domain.subdomain]` | `"ask"` \| `"skip"` | Behavior when no matching skill is found: `"ask"` prompts for clarification, `"skip"` silently moves on |
| `compliance-threshold` | `[standards.domain]` or `[standards.domain.subdomain]` | number | Minimum overall criteria coverage % to pass `check-best-practice-compliance` |
| `compliance-threshold-error` | `[standards.domain]` or `[standards.domain.subdomain]` | number | Maximum allowed error-severity violations — `0` means none allowed |

Domain sections go under `[standards]`: `[standards.engineering]` or subdomain `[standards.engineering.architecture]`. More specific section wins — `[standards.engineering.architecture]` overrides `[standards.engineering]` for architecture work.

---

## Package URL Format

Package URLs encode both the **package source** and a **path/glob selector** in a single string:

```
[host/][owner/repo][@version][:glob_path]
```

| Part | Optional | Default |
|------|----------|---------|
| `host/` | yes | `github.com` |
| `owner/repo` | yes | official package |
| `@version` | yes | latest |
| `:glob_path` | yes | all skills or profiles |

`:` is the explicit delimiter between the package locator (`host/owner/repo@version`) and the path selector (glob pattern). No `:` = all skills from that package. A single segment with no `/` and no `:` = shorthand for a path within the official package.

### Examples

```
"engineering"                                              → official package, glob="engineering"
"/local/path/to/skills"                                    → local Unix path, all skills
"C:\\path\\to\\skills:engineering/bdd"                    → local Windows path, glob="engineering/bdd"
"gitlab.com/google/standards"                              → full package, all skills
"jeffreytse/grimoire-core:engineering/bdd"                 → specific package + path
"acmecorp/standards@v0.1.0:engineering/tdd"               → pinned version + path
"acmecorp/standards@v0.1.0:engineering/**/apply-*"        → version + glob pattern
```

### Glob path semantics

Glob paths match against paths **relative to the package (repo) root** — no `skills/` or `profiles/` prefix is assumed. Different packages may use any directory layout.

Grimoire uses **Standard Glob with globstar**, always root-anchored:

| Wildcard | Matches |
|----------|---------|
| `*` | any characters within one path segment (does not cross `/`) |
| `**` | zero or more path segments (crosses `/`; must be explicit) |
| `?` | single character (does not cross `/`) |
| `{a,b}` | alternation — matches any of the comma-separated alternatives |
| `[abc]` / `[a-z]` | character class / range |
| `[^abc]` / `[!abc]` | negated character class |
| `[[:alpha:]]` | POSIX named class — see below |

All patterns are **root-anchored** — `"engineering"` only matches `<root>/engineering/`, not `<root>/any/depth/engineering/`. Use `"**/engineering"` to match at any depth.

**Supported POSIX classes:** `[[:alnum:]]`, `[[:alpha:]]`, `[[:digit:]]`, `[[:lower:]]`, `[[:upper:]]`, `[[:xdigit:]]`, `[[:space:]]`, `[[:blank:]]`, `[[:word:]]`, `[[:graph:]]`, `[[:print:]]`, `[[:punct:]]`

**Skill matching:** walk package root for `SKILL.md` markers; match glob against the containing directory path relative to repo root.

---

## Profile Name Resolution

Profile names in `[standards] profiles` use **depth-anywhere** glob semantics — unlike skills (root-anchored), `"engineering"` automatically finds profiles at any depth in the package tree.

For a name `p`, grimoire searches project `.grimoire/` first, then each installed package in priority order. Within each root, it matches `.toml` files (suffix always inferred) if:
1. **File match:** `**/p` matches the file path — finds `engineering.toml` anywhere
2. **Directory match:** `**/p` matches a parent directory — activates all `.toml` under any `engineering/` dir anywhere

One name can expand to multiple profiles. `"engineering"` activates every matching profile in the tree, not just the first one.

Patterns with `**` are used as-is; patterns without `**` get `**/` prepended automatically:

```
"engineering"        → **/engineering.toml  OR  all .toml under **/engineering/
"engineering/tdd"    → **/engineering/tdd.toml  OR  all .toml under **/engineering/tdd/
"**/tdd"             → used as-is (already has **)
```

Three ref forms, disambiguated by `:` presence:

| Form | Example | Behavior |
|------|---------|----------|
| Simple name | `"engineering"` | Depth-anywhere glob across all installed packages |
| Qualified name | `"acmecorp/engineering"` | Package-name prefix match → search only that package |
| Package URL | `"acmecorp/standards@v0.1.0:engineering"` | Parse package + version + glob path |

```toml
[standards]
profiles = [
  # simple name — depth-anywhere across all packages
  "clean-architecture",
  # any tdd.toml anywhere, or all .toml under any tdd/
  "tdd",
  # qualified — search only within the acmecorp package
  "acmecorp/engineering",
  # package URL — pinned version, engineering glob
  "acmecorp/standards@v0.1.0:engineering",
  # package URL — explicit host, latest version
  "github.com/acmecorp/standards:engineering",
]
```

---

## Dependencies

`[dependencies]` declares which skills the project uses. Each key is a package URL. In MCP mode the server prefers these practices (or, with `[standards] scope = "strict"`, serves only them), and `grimoire update` clones any missing packages. In skills mode, `grimoire install` with no args resolves the list, clones missing packages and links the matched skills into each agent.

```toml
# grimoire.toml
[dependencies]
# official package, skills under engineering/
"engineering" = "*"
# local path, all skills
"/local/path/to/skills" = "*"
# all skills from that package
"gitlab.com/google/standards" = "*"
# exact skill
"acmecorp/standards@v0.1.0:engineering/tdd/apply-internal-convention" = "*"
# skills under engineering/tdd/
"acmecorp/new-standards@v0.2.1:engineering/tdd" = "*"
# skills under engineering/bdd/
"jeffreytse/grimoire-core:engineering/bdd" = "*"
```

**Which skills an entry selects:** each entry picks every skill under its path (or all of the package's skills when there's no path). Packages you don't have yet are fetched by `grimoire update` or `grimoire install`. If two entries select the same skill, it is used once — the first entry wins. Active `[standards] profiles` narrow the selection to the skills those profiles include.

`grimoire install <package-ref>` adds the ref to `[dependencies]` and, in skills mode, links it into the agents immediately:

```bash
grimoire install jeffreytse/grimoire-core:engineering/bdd         # project-scoped (default)
grimoire install acmecorp/standards@v0.1.0:engineering --global  # global-scoped
```

`grimoire uninstall <package-ref>` removes matching skills and the entry from settings:

```bash
grimoire uninstall jeffreytse/grimoire-core:engineering/bdd
grimoire uninstall acmecorp/standards:engineering --global
```

---

## Packages

Grimoire uses a package model similar to `apt` or `pacman`: one **official** package provides the canonical skills source, and any number of **user packages** extend it with domain-specific or team-specific practices.

Packages are declared as a `[[package]]` table array in global settings:

```toml
# ~/.config/grimoire/grimoire.toml

[[package]]
name     = "official"
url      = "jeffreytse/grimoire-core"  # default — omit to use built-in
official = true
priority = 100

[[package]]
name     = "acmecorp"
url      = "acmecorp/standards"
priority = 50
```

### `[[package]]` fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | required | Package identifier. Used in qualified profile refs (`"acmecorp/engineering"`). |
| `url` | string | required | `owner/repo[@version]`, full Git URL, or absolute local path. |
| `official` | bool | `false` | Marks this as the official skills source. Only the first `official = true` entry is used for resolution — `grimoire doctor` warns when multiple exist. |
| `priority` | int | 100 (official) / 50 (user) | Search order — higher wins. First match per skill name is used. |
| `enabled` | bool | `true` | Set `false` to temporarily disable without removing the entry. |

### Official package

The official package is the canonical source for skills and profiles. There is one active official package at a time (first `official = true` entry wins). If no `[[package]]` is configured, grimoire falls back to the built-in default (`jeffreytse/grimoire-core`).

Override to use a corporate mirror or fork:

```bash
grimoire package set acmecorp/grimoire-core@v2.1.0
grimoire package reset   # revert to built-in default
```

### User packages

User packages extend the official source without replacing it — like adding a PPA in apt. When you add a package, grimoire automatically loads its `grimoire.toml` as a base configuration layer, so its profiles, practices, and compliance thresholds apply to your projects. Your own settings always win over package-provided defaults.

```bash
grimoire package add acmecorp acmecorp/standards
grimoire package add myteam  github.com/myorg/practices@v1
```

After adding, `grimoire package update` clones and pulls all enabled packages. Their profiles and practices are then available under the package name.

Skills in user packages do not need to follow STANDARD.md — any file format is accepted. User package authors may optionally follow STANDARD.md to enable `review-best-practice-skill` interoperability and CI validation badges.

---

## `[core.check-provider]`

Configures the API provider used by `grimoire check` when no local AI CLI is detected. Built-in named providers have sensible defaults — only override the fields that differ.

| Name | Base URL | API Key Env | Default model |
|------|----------|-------------|---------------|
| `anthropic` | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` | `claude-haiku-4-5-20251001` |
| `openai` | `https://api.openai.com/v1` | `OPENAI_API_KEY` | `gpt-4o-mini` |
| `openrouter` | `https://openrouter.ai/api/v1` | `OPENROUTER_API_KEY` | `anthropic/claude-3.5-haiku` |
| `grok` | `https://api.x.ai/v1` | `XAI_API_KEY` | `grok-3-mini` |
| `groq` | `https://api.groq.com/openai/v1` | `GROQ_API_KEY` | `llama-3.1-8b-instant` |
| `ollama` | `http://localhost:11434/v1` | (none — local process) | `llama3.2` |

### `[core.check-provider]` fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Built-in provider name (table above) or `"custom"` for arbitrary endpoints. |
| `base-url` | string | Override the base URL. Empty = use built-in default for named provider. |
| `api-key-env` | string | Env var holding the API key. Empty = use built-in default. |
| `model` | string | Model identifier. Empty = use built-in default. |
| `format` | string | `"anthropic"` or `"openai"`. Empty = derived from `name`. |

### Resolution order for `grimoire check`

1. `--via <agent>` flag — force a specific local CLI
2. Local CLIs — tried in `[core] check-agents` order (default: `["claude","gemini","codex","copilot","opencode","openclaw"]`)
3. `[core.check-provider]` — configured API provider
4. Auto-detect — scans `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `OPENROUTER_API_KEY`, `XAI_API_KEY`, `GROQ_API_KEY` in that order

### Examples

```toml
# ~/.config/grimoire/grimoire.toml

[core]
check-agents = ["gemini", "claude", "codex"]  # prefer gemini first

[core.check-provider]
name  = "openrouter"
model = "meta-llama/llama-3.1-70b-instruct"   # override model only; URL + key use built-in defaults
```

Custom internal endpoint:

```toml
[core.check-provider]
name        = "custom"
base-url    = "https://llm.internal.corp/v1"
api-key-env = "INTERNAL_LLM_KEY"
model       = "corp-llm-v2"
format      = "openai"
```

---

## Examples

### Global defaults

Applies to every project for this user. Good for personal style preferences and company-wide packages.

```toml
# ~/.config/grimoire/grimoire.toml

[[package]]
name     = "official"
# use company mirror for all projects
url      = "acmecorp/grimoire-core"
official = true

[[package]]
name     = "acmecorp"
# team-wide practices
url      = "acmecorp/standards"

[dependencies]
# official package engineering skills
"engineering" = "*"
# acmecorp engineering skills
"acmecorp/standards:engineering" = "*"

[standards]
# OOP skills active everywhere by default
profiles = ["oop"]

[standards.engineering]
practices = ["Google Engineering Practices"]

[standards.engineering.testing]
practices = ["apply-test-driven-development"]

[standards.finance]
practices = ["CFA Institute standards"]

[standards.law]
practices = ["ABA Model Rules"]
# always ask before applying a law skill — don't guess jurisdiction
fallback = "ask"
```

---

### Project config

Committed to the repo. Every team member gets the same baseline.

```toml
# grimoire.toml
[dependencies]
# official engineering skills
"engineering" = "*"
# pinned team TDD practices
"acmecorp/standards@v1.2.0:engineering/tdd" = "*"

[standards]
# skills from both; priority field resolves per-skill conflicts
profiles  = ["clean-architecture", "tdd"]

[standards.engineering]
practices = ["Google Engineering Practices"]
# 75% overall coverage to pass CI
compliance-threshold = 75

[standards.engineering.architecture]
practices = [
  "apply-solid-principles",
  "apply-law-of-demeter",
  "apply-domain-driven-design"
]
# architecture held to a higher bar
compliance-threshold = 80
# zero tolerance for error-severity findings
compliance-threshold-error = 0

[standards.engineering.testing]
practices = ["apply-test-driven-development"]
# if no exact test skill matches, skip silently
fallback = "skip"
```

---

### Multi-package project

Pulling skills from the official package and a team package, with glob path selectors:

```toml
# grimoire.toml
[dependencies]
# all apply-* skills anywhere under engineering/
"engineering/**/apply-*" = "*"
# tdd skills from team package
"acmecorp/standards:engineering/tdd" = "*"
# specific path from grimoire-core
"jeffreytse/grimoire-core:security/owasp" = "*"

[standards]
profiles = [
  # depth-anywhere: any clean-architecture.toml or clean-architecture/ dir
  "clean-architecture",
  # package URL — search engineering glob in acmecorp/standards
  "acmecorp/standards:engineering",
]
```

---

### Multi-domain project (startup)

A project that spans engineering, legal, and finance — each domain has its own standard:

```toml
# grimoire.toml
[dependencies]
"engineering" = "*"
"law" = "*"
"finance" = "*"

[standards]
profiles = ["lean-startup"]

[standards.engineering]
practices = ["Google Engineering Practices"]
compliance-threshold = 70       # early-stage: 70% is the bar

[standards.law.contracts]
practices = ["review-saas-contract", "audit-gdpr-compliance"]
compliance-threshold-error = 0  # legal errors are always blocking

[standards.finance.personal-finance]
practices = ["CFA Institute standards"]

[standards.writing.copywriting]
practices = ["apply-aida", "write-value-proposition"]
```

---

### Enforcement-only project (legacy codebase)

Existing codebase with tech debt — set a lower threshold now, raise it over time:

```toml
# grimoire.toml
[dependencies]
"engineering" = "*"

[standards.engineering]
practices = [
  "apply-solid-principles",
  "apply-kiss-principle"
]
compliance-threshold = 40       # current baseline — raise 5% per sprint
compliance-threshold-error = 0  # errors always block regardless of overall score
```

---

### CI/check provider config

Team settings that configure which AI checks compliance and how CI consumes the results:

```toml
# ~/.config/grimoire/grimoire.toml
[core]
check-agents = ["claude", "gemini"]  # try claude first, fall back to gemini

[core.check-provider]
name = "openrouter"                  # API fallback if no local CLI is found
model = "anthropic/claude-3.5-haiku"
```

```yaml
# .github/workflows/compliance.yml
steps:
  - name: Run compliance check
    run: grimoire check --from-report --ci --junit .grimoire/reports/junit.xml

  - name: Upload JUnit report
    uses: actions/upload-artifact@v4
    with:
      name: junit-compliance
      path: .grimoire/reports/junit.xml
```

`--ci` emits GitHub Actions `::error` / `::warning` annotations per diagnostic and auto-enables when `GITHUB_ACTIONS=true`. `--junit` writes a JUnit XML file for CI test reporters.

---

### Session pin (via AI — no TOML edit needed)

Grimoire can pin a preference for the current session without writing any file:

```
User: For this session, use Facebook's engineering practices instead of Google's.

Claude: Pinned for this session: Facebook Engineering Standards → [engineering].
        Reverts when session ends. Save permanently? [y] global  [p] project
```

The session pin overrides all files. Nothing is written to disk until you confirm.

---

## How Layers Compose

Given this setup:

```
~/.config/grimoire/grimoire.toml   [[package]] name="acmecorp" url="acmecorp/standards"
                                   [standards] profiles = ["oop"]
                                   [standards.engineering] practices = ["Google"]

grimoire.toml            [standards.engineering.architecture] practices = ["SOLID", "LoD"]
                                   [standards.engineering.architecture] compliance-threshold = 80
```

The resolved spec for `engineering.architecture` work:

| Key | Value | Source |
|-----|-------|--------|
| `[standards].profiles` | `["oop"]` | global |
| `[standards.engineering].practices` | `["Google"]` | global (no local override for this key) |
| `[standards.engineering.architecture].practices` | `["SOLID", "LoD"]` | local |
| `[standards.engineering.architecture].compliance-threshold` | `80` | local |

`[standards.engineering.architecture]` practices do not fall back to `[standards.engineering]` practices — they replace them for architecture-scoped work. Other keys not set at the subdomain level do inherit from the domain level.

Standards contributed by user packages (their `grimoire.toml`) are loaded as the lowest-priority layer — project and global settings always win over package-provided defaults.

---

## CLI Management

Read and write settings without editing TOML manually:

```bash
# Read
grimoire config show                     # show effective settings with per-key source
grimoire config get standards.profiles   # show resolved value and source file

# Write (default level: standards.* → --local, core.* → --global)
grimoire config set standards.profiles "clean-architecture,tdd"
grimoire config set standards.engineering.compliance-threshold 75
grimoire config set standards.engineering.testing.practices "apply-test-driven-development"
grimoire config set standards.engineering.fallback skip

# Write to a specific level
grimoire config set standards.profiles "oop" --global
grimoire config set core.home /opt/grimoire --global
grimoire config set standards.engineering.practices "corp-standards" --system

# Unset
grimoire config unset standards.engineering.compliance-threshold
```

Array values are comma-separated. `--local` writes to `grimoire.toml`, `--global` to `~/.config/grimoire/grimoire.toml`, `--system` to `/etc/grimoire/grimoire.toml` (requires admin).

Dependency management:

```bash
grimoire install                                        # install all skills in [dependencies] (skills mode)
grimoire install <package-ref>                          # add ref to settings, clone, install
grimoire install jeffreytse/grimoire-core:engineering   # project-scoped (default)
grimoire install acmecorp/standards --global           # global-scoped
grimoire uninstall <package-ref>                       # remove matching skills + settings entry
grimoire uninstall jeffreytse/grimoire-core:engineering/bdd
grimoire uninstall acmecorp/standards:engineering --global
```

Package management:

```bash
grimoire package list                         # show all configured packages
grimoire package add <name> <url>             # add a user package
grimoire package remove <name>                # remove a package by name
grimoire package set owner/repo@v1            # set the official package source
grimoire package reset                        # revert official package to built-in default
grimoire package update                       # pull latest for all enabled packages
grimoire package update acmecorp              # pull a single package
grimoire package publish                      # publishing checklist + open github.com/new
```

Compliance check:

```bash
grimoire check                                        # independent AI check (auto-selects local CLI or API)
grimoire check --via claude                           # force a specific local agent
grimoire check --from-report                          # read pre-generated compliance report
grimoire check --ci                                   # GitHub Actions annotations (auto on GITHUB_ACTIONS=true)
grimoire check --junit .grimoire/reports/junit.xml    # write JUnit XML report
```

Real-time watch:

```bash
grimoire watch                         # re-run check on every file save
grimoire watch --via gemini            # use specific agent for watch mode
```

Project health:

```bash
grimoire status                        # profile, skills count, last check, staleness
```

Init:

```bash
grimoire init                        # TUI profile picker; auto-detects project type
grimoire init --profile engineering  # skip TUI, apply named profile directly
grimoire init --yes                  # non-interactive; auto-detected profile, no prompts
grimoire init --auto                 # full auto-detect + install
```

---

## Guided Management

Use skills instead of editing TOML directly:

- `/configure-grimoire` — view the resolved spec for any domain, edit and validate settings interactively
- `/pin-best-practice-preference` — pin a preference when two practices conflict; choose session, project, or global scope
- `/resolve-best-practice-conflict` — resolve contradictions between two installed skills and record the priority automatically

Or inspect the resolved state with full source attribution:

```bash
grimoire config show                   # effective settings with per-key source file
grimoire context                       # full ground truth: settings + resolved profiles + agents + compliance
grimoire profile show <name>           # resolve a single profile with composition chain
```

---

## Related

- [profiles.md](./profiles.md) — how `profiles` resolves, custom profile files, `profiles` vs `practices`
- [bpdd.md](./bpdd.md) — compliance thresholds in the BPDD cycle, CI gating, incremental mode


## Overview: layers, precedence and guided configuration

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

| Priority | Level   | Stored in                          | Scope                                           |
| -------- | ------- | ---------------------------------- | ----------------------------------------------- |
| 1st      | Session | In-memory                          | Current session only — resets when session ends |
| 2nd      | Project | `grimoire.toml`                    | Current project, committed to repo              |
| 3rd      | Global  | `~/.config/grimoire/grimoire.toml` | All projects on this machine                    |
| 4th      | System  | `/etc/grimoire/grimoire.toml`      | All users on this machine                       |

**What each layer can declare.** Libraries, then the people and teams that use them: an organization lists its libraries and standards in the system file, each user adds their own, and each project adds its standards on top:

| | System (org) | Global (user) | Project |
| --- | --- | --- | --- |
| Libraries (`[[package]]`) | ✓ | ✓ (wins per name, e.g. to disable an org library) | — via `[dependencies]` only |
| `[dependencies]` | ✓ | ✓ | ✓ |
| Pinned practices `[standards.<domain>] practices` | ✓ | ✓ | ✓ (nearest wins per domain) |
| `[standards] profiles`, `extends` | ✓ | ✓ | ✓ |
| `[standards] scope = "strict"` | — | — | ✓ |
| `[core]` settings | ✓ | ✓ | ✓ except `home` |

A project can't add a skill source with `[[package]]`: a cloned repository must not decide what your agent loads. It names the libraries it needs in `[dependencies]`, which `grimoire update` or `grimoire install` fetch when you run them. Declared practices show which layer they came from: `standard: "project"`, `"user"` or `"org"`. `GRIMOIRE_SYSTEM_CONFIG` moves the system file, for containers and managed machines.

**Configure manually** — edit the settings files directly without going through the AI:

```toml
# ~/.config/grimoire/grimoire.toml  (global — applies everywhere)
[standards]
profiles = ["oop"]  # activate all skills tagged "oop"
# profiles = ["clean-architecture", "tdd"]  # multiple — first entry wins conflicts

[standards.engineering]
practices = ["Google Engineering Practices"]

[standards.finance]
practices = ["CFA Institute standards"]
```

```toml
# grimoire.toml  (project — overrides global for this repo)
[standards.engineering.architecture]
practices = [
  "SOLID principles: production code",
  "KISS: prototypes, scripts"
]
fallback = "ask"
note = "agreed in ADR-014"  # optional reason, shown when the pin is applied
```

Or from the shell, without editing TOML: `grimoire config add standards.engineering.practices "Google Engineering Practices" -g` (drop `-g` for the project file), and `grimoire config get <key>` to see the resolved value and where it came from.

Every skill that reads or saves a preference — `pin-best-practice-preference`, `suggest-best-practice`, `audit-applied-best-practices`, … — uses these same files and keys, in skills mode and MCP mode alike.

Project configuration overrides global. Session pins override both. Teams can share a global standard while individual projects deviate where needed.

**`practices = ["OOP"]` vs `profiles = ["oop"]`** — both signal OOP intent, but differently. `practices = ["OOP"]` in a `[standards.<domain>]` section is a loose hint — the AI leans toward OOP conventions from its training. `profiles = ["oop"]` under `[standards]` activates specific installed skills (exact steps, validated sources). Use `profiles` for precision; `practices` for domain-level style preference. → [Full comparison](./profiles.md#profiles-vs-practices)

**Guided configuration:**

- `/configure-grimoire` — view, edit, or validate config without touching TOML directly
- `/apply-best-practice-profile` — activate a full paradigm (OOP, TDD, clean architecture) in one command
- `/resolve-best-practice-conflict` — resolve contradictions between two installed skills and record the priority automatically
- `/apply-best-practice-driven-development` — run the full BPDD cycle ([BPDD guide](./bpdd.md))

→ [Full configuration reference](./config.md) — all keys, override hierarchy, TOML examples
