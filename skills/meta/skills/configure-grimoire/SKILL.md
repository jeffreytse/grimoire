---
name: configure-grimoire
description: Use when the user wants to view, edit, remove, or validate their grimoire settings — including reading current preferences, changing or deleting a setting, switching a named profile, or checking settings.toml for contradictions and expired entries.
source: Settings management patterns (VS Code settings UI, Git config get/set/unset); XDG Base Directory Specification (freedesktop.org)
tags: [settings, configuration, preferences, profile, toml, view, edit, validate]
---

# Configure Grimoire

Read, update, or validate `settings.toml` — the single source of truth for all grimoire preferences and configuration.

## Why This Is Best Practice

**Adopted by:** Every major developer tool provides a guided settings interface alongside raw file editing — VS Code `settings.json` + Settings UI, Git `git config get/set/unset`, npm `.npmrc` + `npm config`. Guided operations prevent syntax errors and silent overwrites that raw editing can introduce.
**Impact:** Configuration drift — where users forget what they set, or stale settings silently govern behavior — is one of the primary causes of "why is it doing that?" confusion. VS Code's settings validation (invalid key detection) reduced user-reported configuration bugs by catching errors at write time rather than at apply time. Explicit confirmation before destructive operations (delete, overwrite) is standard across all major CLI tools (`git config --unset`, `npm config delete`).
**Why best:** Manual TOML editing is error-prone (syntax, wrong key names, stale dates). A guided interface validates before writing, confirms before deleting, and surfaces the full settings cascade so the user knows which file actually governs a domain.

Sources: XDG Base Directory Specification (freedesktop.org); VS Code settings documentation; Git config documentation (`git help config`)

## Steps

### Step 1: Detect operation

| User signal | Operation |
|-------------|-----------|
| "show", "view", "what are my settings for X", "list my preferences" | `show` |
| "edit", "update", "change", "set X to Y" | `edit` |
| "remove", "delete", "unset", "clear" | `remove` |
| "switch to X profile", "use X profile", "activate profile" | `switch-profile` |
| "validate", "check settings", "any issues" | `validate` |

---

### Step 2: Load settings

Resolve the active settings by merging in order (highest wins):

```
<project>/.grimoire/settings.local.toml  (project personal — gitignored)
  > <project>/.grimoire/settings.toml   (project shared — committed)
    > ~/.config/grimoire/settings.toml  (global — XDG primary)
```

Fall back to `~/.grimoire/settings.toml` if `~/.config/grimoire/` does not exist.

Within any file, cascade by specificity: `[domain.subdomain] > [domain] > [global]`

---

### Step 3a: Show

Display the merged effective settings for the requested domain (or all domains):

```
Effective settings for engineering/architecture:
  Source: .grimoire/settings.toml (project)

  practices: ["SOLID principles: production code", "KISS: scripts, prototypes"]
  fallback:  "ask"
  author:    "@backend-team"
  note:      "agreed in ADR-014"
  expires:   "2026-09-01"

  Active profile: (none — using default)
  Overrides from: (no global entry for this domain)
```

Show which file each key came from if multiple files contribute.

---

### Step 3b: Edit

Parse the requested change from user input. Show what will change, then ask which file to write to:

```
Change:

  [engineering.architecture]
  - fallback = "ask"
  + fallback = "both"

Write to:
  [1] This project (shared)    → .grimoire/settings.toml       (committed to repo)
  [1b] This project (personal) → .grimoire/settings.local.toml (gitignored)
  [2] All projects (global)    → ~/.config/grimoire/settings.toml
```

After file selection, confirm and write. Never write invalid TOML.

---

### Step 3c: Remove

Ask which file to remove from, then confirm:

```
Remove from:
  [1] This project (shared)    → .grimoire/settings.toml       (committed to repo)
  [1b] This project (personal) → .grimoire/settings.local.toml (gitignored)
  [2] All projects (global)    → ~/.config/grimoire/settings.toml
```

After file selection, show what will be removed and confirm:

```
Remove from .grimoire/settings.toml:

  [engineering.architecture]
  expires = "2026-09-01"   ← remove this key

Or remove the entire [engineering.architecture] section? [key / section / cancel]
```

If removing from the project file would expose a global default the user didn't intend, show the effective value after removal before confirming.

After confirmation, remove the key or section. If the domain section becomes empty, ask whether to remove the section header too.

---

### Step 3d: Switch profile

Store the active profile for a domain in `settings.local.toml` (personal, never committed):

```toml
# settings.local.toml
[engineering.architecture]
active-profile = "prototype"
```

Confirm the switch:

```
Switched engineering/architecture to profile: prototype
Stored in: .grimoire/settings.local.toml (personal, not committed)

Active practices:
  1. KISS
  2. YAGNI
```

To reset to default: remove `active-profile` from `settings.local.toml`.

---

### Step 3e: Validate

Check the merged settings for issues:

| Check | Issue |
|-------|-------|
| `require` + `disabled` overlap | Contradiction — skill can't be both required and disabled |
| `expires` date < today | Preference has expired — flag for review or removal |
| `remind` date ≤ today | Reminder due — surface to user |
| Unknown key names | Typo in key name — flag |
| `[domain.profiles.X]` referenced but no matching profile | Profile defined but never activated |

Output:

```
Validating settings files...

  .grimoire/settings.local.toml (project personal)
    ✅  engineering/architecture: OK

  .grimoire/settings.toml (project shared)
    ⚠️  engineering/architecture: expires "2026-09-01" is past — still apply? [keep / remove]
    ❌  engineering/testing: "apply-tdd" is in both require and disabled — contradiction
    ✅  engineering/development: OK

  ~/.config/grimoire/settings.toml (global)
    ✅  global: OK

1 error, 1 warning across 3 files. Fix errors before conflicts can be resolved correctly.
```

## When NOT to Use

- **Adding a new preference from scratch**: use `pin-best-practice-preference` — it guides the full pin flow with scope selection.
- **Resolving conflicts between skills**: use `resolve-best-practice-conflict` — it loads both SKILL.md files, identifies the specific contradiction, and updates priorities.
- **Installing or upgrading grimoire**: use `install-grimoire`.

## Common Mistakes

**Editing the wrong file**: always confirm which file (project shared, project personal, global) before writing. A setting in `settings.toml` commits to repo; a setting in `settings.local.toml` stays personal.

**Removing without checking cascade**: deleting a key from the project file may expose a global default the user didn't intend. Show the effective value after removal before confirming.

**Profile confusion**: `active-profile` is always personal (`settings.local.toml`). Never write it to shared `settings.toml` — it would commit a personal profile choice for the whole team.
