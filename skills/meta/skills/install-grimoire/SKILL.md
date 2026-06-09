---
name: install-grimoire
description: Use when the user wants to install or uninstall grimoire skills by domain or individual skill, upgrade grimoire to the latest version, clean up broken symlinks, list what skills are available, or run a health check on the grimoire installation.
source: Package manager UX patterns (Homebrew, npm, apt); grimoire scripts/grimoire
tags: [install, uninstall, upgrade, clean, list, doctor, health, diagnostics, grimoire, skills, domains, setup]
---

# Install Grimoire

Install, uninstall, or upgrade grimoire skills — a guided interface over `scripts/grimoire`.

## Why This Is Best Practice

**Adopted by:** Every major package manager (Homebrew, npm, apt, pip) provides a guided CLI with confirmation prompts before destructive operations and clear output of what was changed. The pattern of "show command → confirm → execute → report" is universal because it prevents silent failures and makes installs auditable.
**Impact:** Silent install failures (wrong path, broken symlink, permission error) are the primary cause of "the skill isn't working" confusion. Explicit confirmation before uninstall and clear post-install reporting (what was installed, where, how many) eliminates the ambiguity. Homebrew's `brew install` output format — listing each installed path — is the gold standard for install UX.
**Why best:** Running `grimoire` directly requires knowing the exact flags. A guided skill extracts intent from natural language ("install the engineering skills for Claude"), constructs the correct command, confirms before running, and surfaces the result in plain language.

Sources: Homebrew documentation; npm CLI documentation; grimoire `scripts/grimoire`

## Steps

### Step 1: Detect operation

| User signal | Operation |
|-------------|-----------|
| "install", "add", "set up" | `install` |
| "uninstall", "remove", "delete" | `uninstall` |
| "upgrade", "update", "get latest" | `upgrade` |
| "clean", "fix broken links" | `clean` |
| "list", "what's available", "show domains" | `list` |
| "doctor", "health check", "status", "is everything working" | `doctor` |

---

### Step 2: Gather parameters

For `install` and `uninstall`, determine scope:

| Scope | Flag |
|-------|------|
| All skills | *(no domain flag)* |
| One domain (e.g. "engineering") | `--domain engineering` |
| One subdomain | `--domain engineering --subdomain development` |
| One skill | `--skill engineering/development/apply-kiss-principle` |

For `install`, also determine target agent:

| Target | Flag |
|--------|------|
| All detected agents | `--target all` *(default)* |
| Claude Code only | `--target claude` |
| Codex only | `--target codex` |
| Gemini CLI only | `--target gemini` |

If scope or target is ambiguous, ask one question before proceeding.

---

### Step 3: Show command and confirm

Construct the `scripts/grimoire` command and show it before running:

```
Will run:
  ./scripts/grimoire --domain engineering --target claude

This will install all engineering skills (~101 skills) to Claude Code.
Continue? [y/n]
```

For `uninstall`, flag it explicitly:

```
Will run:
  ./scripts/grimoire --domain business --uninstall --target all

This will REMOVE all business skills from all agents. Cannot be undone without re-installing.
Continue? [y/n]
```

For `upgrade`:

```
Will run:
  ./scripts/grimoire --upgrade

This will pull the latest grimoire from GitHub and refresh all symlinks.
Continue? [y/n]
```

---

### Step 4: Execute

Run the confirmed command via Bash. Stream or capture output.

---

### Step 5: Report result

```
✅ Installed 101 skills from engineering domain to Claude Code.

Installed at: ~/.claude/skills/

Domains installed:
  engineering/development        (18 skills)
  engineering/architecture       (12 skills)
  engineering/testing            (9 skills)
  ... (and 8 more sub-domains)
```

For `upgrade`:

```
✅ Grimoire upgraded to latest.

  Previous: commit abc1234 (2026-05-01)
  Current:  commit def5678 (2026-06-09)

  New skills: 14
  Updated skills: 7
  Symlinks refreshed: 202
```

For `list`:

Show available domains and skill counts. For subdomain or skill scope, show names.

For `doctor`:

Run `./scripts/grimoire --doctor` directly (read-only, no confirmation needed). Output shows 3 sections:

```
Grimoire health check

  Source
    ✅  git repo:    /path/to/grimoire (commit abc1234, 2026-06-09)
    ✅  grimoire:  executable

  Installed skills
    ✅  claude:   312 skills, 0 broken symlinks
    ⚠️   codex:    88 skills, 3 broken symlinks  → run: --clean --target codex
    ⬜  gemini:  (not detected — skipped)

  Config
    ✅  project personal (.grimoire/settings.local.toml) — present, valid TOML
    ⬜  project shared (.grimoire/settings.toml) — not found
    ⬜  global (~/.config/grimoire/settings.toml) — not found

  Summary: 1 warning.
```

---

## grimoire flag reference

| Flag | Effect |
|------|--------|
| `--domain <name>` | Scope to one domain (e.g. `engineering`) |
| `--subdomain <name>` | Scope to one subdomain — requires `--domain` |
| `--skill <path>` | One skill (e.g. `engineering/development/apply-kiss-principle`) |
| `--target <agent>` | `claude`, `codex`, `gemini`, or `all` |
| `--uninstall` | Remove instead of install |
| `--upgrade` | `git pull` latest + refresh all symlinks |
| `--clean` | Remove broken symlinks from all agent skill dirs |
| `--list` | List available domains, subdomains, and skills |
| `--doctor` | Read-only health check: git repo, symlinks per agent, config files |
| `--copy` | Copy files instead of symlinking (for environments where symlinks don't work) |

## When NOT to Use

- **Changing preferences or settings**: use `configure-grimoire` or `pin-best-practice-preference`.
- **Resolving skill conflicts**: use `resolve-best-practice-conflict`.
- **Writing a new skill**: use `write-best-practice-skill`.

## Common Mistakes

**Uninstalling instead of upgrading**: if the user says "update my skills," that almost always means `--upgrade` (refresh to latest), not `--uninstall`. Confirm the intended operation before running.

**Wrong target agent**: defaulting to `--target all` is usually correct, but if the user only uses one agent, installing to all wastes disk space with unused symlinks. Ask if ambiguous.

**Skipping confirmation on uninstall**: always show the uninstall command and warn it's destructive before running. Never auto-run uninstall without explicit user confirmation.
