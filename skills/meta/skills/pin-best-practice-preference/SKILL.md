---
name: pin-best-practice-preference
description: Use when the user wants to remember, save, or pin a specific best practice preference — either a new choice ("remember I prefer X for Y", "always use X", "pin X") or to promote existing session-level pins to project or global storage for future sessions ("save my session preferences", "persist my choices").
source: User preference persistence patterns (browser localStorage, IDE settings sync, dotfile conventions); XDG Base Directory Specification (freedesktop.org)
tags: [preference, pin, remember, session, persist, user-choice, settings, quick-pin]
---

# Pin Best Practice Preference

Let the user pin a best practice preference in one step — either a new intentional choice or existing session pins promoted to disk for future sessions.

## Why This Is Best Practice

**Adopted by:** Every major developer tool persists user preferences across sessions — VS Code settings.json, Git config, npm .npmrc, SSH config. The XDG Base Directory Specification (freedesktop.org) formalizes the hierarchy: project-level config overrides user-level config, which overrides system-level config. Browser vendors (Chrome, Firefox) implement the same pattern: in-session choices can be persisted to profile storage.
**Impact:** Session-only preferences require the user to re-specify choices every session — the primary driver of preference drift, where users accept suboptimal defaults rather than repeat corrections. Persistent preferences reduce cognitive load by eliminating re-specification: a study on IDE preference systems (Ko et al., 2004, CHI) found that users who could persist tool preferences spent 40% less time on configuration tasks per session.
**Why best:** The alternative — manual editing of config files — requires users to know the file location, format, and valid values. A guided pin flow eliminates all three barriers while producing the same persistent artifact. Session-first (option 0) respects users who want to try a preference before committing it to disk.

Sources: XDG Base Directory Specification (freedesktop.org); Ko et al. (2004) "Designing the Whyline: A Debugging Interface for Asking Questions about Program Behavior," CHI

## Steps

### Step 1: Detect mode (silent)

| Signal | Mode |
|--------|------|
| User names a practice or domain ("remember X for Y", "always use X", "pin X", "I prefer X") | Quick pin |
| User asks to save/promote session choices ("save my session preferences", "promote", "persist my choices", "remember everything from this session") | Session promotion |

### Step 2a: Quick pin

Extract from user input:
- **Practice**: skill name or natural description (e.g., "Jest", "conventional commits", "4% withdrawal rate")
- **Domain**: explicit or inferred from context (e.g., "engineering/testing", "finance/personal-finance")

If domain is ambiguous, ask ONE question:
```
Which domain is this for? (e.g. engineering/testing, finance/personal-finance, health/fitness)
```

Confirm before writing:
```
Pin "[practice]" as your preference for [domain]?
[y] yes  [n] cancel  [e] edit details first
```

If user selects `[e]`, ask: "Describe the preference in more detail (e.g. tool name, version, key parameters):"

After confirmation, proceed to Step 3.

### Step 2b: Session promotion

List all session-level pins accumulated this session:

```
Session preferences to save:

  engineering/testing       → write-unit-test (Jest, co-located)
  engineering/development   → propose-conventional-commit
  finance/personal-finance  → calculate-fire-number (3.5% rate)

Save all to:
  [1] This project (shared)   → <project-root>/.grimoire/settings.toml   (committed)
  [1b] This project (personal) → <project-root>/.grimoire/settings.local.toml  (gitignored)
  [2] All projects             → ~/.config/grimoire/settings.toml
  [3] Both (project + global)
  [4] Choose per preference
```

If user picks `[4]`, ask per preference: "Save '[practice]' for [domain]? [1] project/shared [1b] project/personal [2] global [3] both [n] skip"

If no session pins exist:
```
No preferences pinned this session.
Use "remember I prefer X for Y" to pin one now.
```

After collecting choices, proceed to Step 4 for each selected preference.

### Step 3: Choose persistence level (Quick pin only)

```
Save to:
  [0] This session only          → in memory; resets when session ends
  [1] This project (shared)      → <project-root>/.grimoire/settings.toml      (committed to repo)
  [1b] This project (personal)   → <project-root>/.grimoire/settings.local.toml (gitignored)
  [2] All my projects (global)   → ~/.config/grimoire/settings.toml
                                   ($XDG_CONFIG_HOME/grimoire/settings.toml if XDG_CONFIG_HOME set)
  [3] Both (project shared + global)
```

### Step 4: Write

Write to selected location(s) using TOML format. Domain/subdomain path uses dots (`[engineering.architecture]`). Array order = priority (index 0 = highest).

```toml
# Full example — all keys (all optional except practices)

profiles = ["clean-architecture", "tdd"]  # array order = conflict priority; resolves file → default.toml → tag query (see docs/profiles.md)

[global]
practices = ["simplicity / KISS"]

[engineering]
practices = ["Google's engineering practices"]

[engineering.architecture]
author = "@backend-team"
require = ["apply-input-validation"]
skip-if = ["test files", "config files"]
practices = [
  "SOLID principles: production code",  # index 0 = highest priority; qualifier after ':'
  "KISS: scripts, prototypes"
]
fallback = "ask"
disabled = ["apply-law-of-demeter"]
ask-before = ["apply-solid-principles"]
expires = "2026-09-01"
remind = "2026-07-01"
shared = true
lock = true
note = "agreed in ADR-014 — contact @backend-team for changes"

[engineering.architecture.profiles.prototype]
practices = ["KISS", "YAGNI"]
note = "speed over structure"
```

**Resolution order (most specific wins):**
```
specific entry in [domain.subdomain] > [domain.subdomain] > [domain] > [global]
```
Each level overrides the level below for the same preference. `[engineering.architecture]` overrides `[engineering]` which overrides `[global]`. Use broad sections for defaults; narrow sections for overrides.

- `practices` array — **all listed skills apply**. Array order = conflict-resolution priority (index 0 highest). When two skills contradict on a specific point, the lower-index skill wins that conflict. Non-conflicting guidance from all skills applies regardless of index.
- Context qualifier as `: qualifier` suffix in the practices string (e.g., `"SOLID principles: production code"`). Omit for all-context preference.
- `fallback` — `"ask"` = AI prompts when skills directly contradict. `"both"` = AI surfaces the tension. Omit = higher-index skill wins the conflict silently.
- `disabled` — array of skill names the AI must never apply in this domain, even if installed globally.
- `ask-before` — array of skill names where AI must confirm before applying. Use for high-stakes skills (security, DB migrations).
- `expires` — ISO date string. On/after this date AI warns "preference has expired — still apply?"
- `remind` — ISO date string. Soft nudge only — AI says "time to review [domain] preferences." No warning/block.
- `shared = true` — marks as team-level (commit `settings.toml` to repo). Omit for personal preferences (`settings.local.toml`).
- `lock = true` — less-specific sections (`[domain]`, `[global]`) do not contribute to this section.
- `note` — free text the AI reads when applying skills here. Cited when asked why a preference exists.
- `require` — array of skill names that MUST apply unconditionally. Win over all conflict resolution. Flag if a skill appears in both `require` and `disabled` — that is a contradiction.
- `author` — who set this preference. AI cites it when asked "why do we prefer X here?"
- `skip-if` — array of file/context patterns. AI skips this entire section when current file/context matches.
- `[domain.profiles.name]` — named profile variant. Default (untagged) section is always active. Named profiles inactive until user says "use [name] profile" — then that variant replaces the default for the session.

If file exists: append new domain section only. Never silently overwrite.

If the domain is already pinned in the file, ask before overwriting:
```
[domain] already has "[existing-skill]" pinned. Replace with "[new-skill]"? [y/n]
```

To add rankings or resolve conflicts between existing preferences, use `resolve-best-practice-conflict`.

For session-level (option 0): store in session memory only — do not write any file.

### Step 5: Confirm

```
Pinned: [practice] for [domain]
Saved to: [path(s) or "session memory (resets when session ends)"]
```

## Rules

- Never auto-pin without explicit user confirmation at Step 2a
- Never prompt for a reason — only record it if the user provided one or used `[e]` / `[r]`
- Session-level pins (option 0) are never written to disk under any circumstances
- Project-level overrides global for the same domain — if pinning to "both", write identical content to both files
- XDG compliance: global config is `$XDG_CONFIG_HOME/grimoire/settings.toml` (defaults to `~/.config/grimoire/settings.toml`). Fall back to `~/.grimoire/settings.toml` only if `~/.config/grimoire/` does not exist.
- If a project root cannot be determined, skip option [1] and inform the user: "No project root detected — project-level save unavailable"
- After writing, always confirm the exact path(s) and what was saved

## Common Mistakes

**Auto-pinning**: never write preferences without the user explicitly confirming at Step 2a. Detected ≠ intentional.

**Silent overwrite**: if the domain is already pinned in the file, always ask before replacing. Existing preferences may be carefully chosen.

**Prompting for reasons unnecessarily**: only record a reason if the user volunteered one. Don't ask "why do you prefer this?" on every pin — it creates friction.

**Forgetting the session case**: if the user chose option 0, confirm "saved to session memory" — not a file path. Make the ephemeral nature explicit.
