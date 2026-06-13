---
name: install-grimoire
description: Use when the user wants to install or uninstall grimoire skills by domain or individual skill, upgrade grimoire to the latest version, clean up broken symlinks, list what skills are available, or run a health check on the grimoire installation.
source: Package manager UX patterns (Homebrew, npm, apt); grimoire scripts/grimoire
tags: [uninstall, upgrade, clean, doctor, health, diagnostics, skills, setup]
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

```
# Install all skills
./scripts/grimoire                                                        # installs everything

# Install one domain
./scripts/grimoire --domain engineering                                   # only engineering skills

# Install one subdomain
./scripts/grimoire --domain engineering --subdomain development           # one subdomain

# Install one skill
./scripts/grimoire --skill engineering/development/apply-kiss-principle   # one skill
```

For `install`, also determine target agent:

```
./scripts/grimoire --domain engineering --target all      # all detected agents (default)
./scripts/grimoire --domain engineering --target claude   # Claude Code only
./scripts/grimoire --domain engineering --target codex    # Codex only
./scripts/grimoire --domain engineering --target gemini   # Gemini CLI only
```

If scope or target is ambiguous, ask one question before proceeding.

---

### Step 3: Show command and confirm

Construct the `scripts/grimoire` command and show it before running. Use a platform-aware confirm:
- **Claude Code / OpenCode**: `AskUserQuestion` — options: "Continue (Recommended)" and "Cancel"
- **Gemini CLI**: `ask_user` — `type: "confirm"`
- **Other**: show the block below and wait for `[y/n]`

Install example (other platforms):
```
Will run:
  ./scripts/grimoire --domain engineering --target claude

This will install all engineering skills (~101 skills) to Claude Code.
Continue? [y/n]
```

For `uninstall`, flag it explicitly (same platform-aware confirm — Claude Code/OpenCode/Gemini use their native tool):
```
Will run:
  ./scripts/grimoire --domain business --uninstall --target all

This will REMOVE all business skills from all agents. Cannot be undone without re-installing.
Continue? [y/n]
```

For `upgrade` (same pattern):
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

**Partial-failure handling:** After running the install command, verify each requested component installed successfully. If some domains installed and others failed:
1. List what succeeded and what failed with error reason
2. Do not report 'installation complete' if any component failed
3. Offer retry for failed components: 'Retry failed installs? [y/n]'
4. If a domain fails due to network/permission error vs. not-found error, distinguish them — not-found means the domain name is wrong; network/permission means retry may work.

**Terminal conditions:**
- Max retries: 2 per failed component
- After 2 failures: mark component as FAILED, continue with remaining components
- Error type routing:
  - `404 / not found / unknown domain`: wrong name — do NOT retry; ask user to verify the domain name
  - `network timeout / 503`: transient — retry up to 2×
  - `403 / permission denied`: stop retrying; output manual install command: `./scripts/grimoire --skill [failed-component]`
  - Any other error: retry once; if still failing, mark FAILED and continue

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

## When NOT to Use

- **Changing preferences or settings**: use `configure-grimoire` or `pin-best-practice-preference`.
- **Resolving skill conflicts**: use `resolve-best-practice-conflict`.
- **Writing a new skill**: use `write-best-practice-skill`.

## Common Mistakes

**Uninstalling instead of upgrading**: if the user says "update my skills," that almost always means `--upgrade` (refresh to latest), not `--uninstall`. Confirm the intended operation before running.

**Wrong target agent**: defaulting to `--target all` is usually correct, but if the user only uses one agent, installing to all wastes disk space with unused symlinks. Ask if ambiguous.

**Skipping confirmation on uninstall**: always show the uninstall command and warn it's destructive before running. Never auto-run uninstall without explicit user confirmation.
