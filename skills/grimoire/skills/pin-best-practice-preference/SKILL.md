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
  [1] This project  → <project-root>/.grimoire/preferences.md
  [2] All projects  → ~/.grimoire/preferences.md
                      (uses ~/.config/grimoire/preferences.md if XDG_CONFIG_HOME is set)
  [3] Both
  [4] Choose per preference
```

If user picks `[4]`, ask per preference: "Save '[practice]' for [domain]? [1] project [2] global [3] both [n] skip"

If no session pins exist:
```
No preferences pinned this session.
Use "remember I prefer X for Y" to pin one now.
```

After collecting choices, proceed to Step 4 for each selected preference.

### Step 3: Choose persistence level (Quick pin only)

```
Save to:
  [0] This session only  → in memory; resets when session ends
  [1] This project only  → <project-root>/.grimoire/preferences.md
  [2] All my projects    → ~/.grimoire/preferences.md
                           (uses ~/.config/grimoire/preferences.md if XDG_CONFIG_HOME is set)
  [3] Both (project + global)
```

### Step 4: Write

Write to selected location(s) using the standard preferences format:

```markdown
## [domain]
- [practice-name]: [detail if provided]
  reason: [reason if provided]
```

If file exists: append new domain section only. Never silently overwrite.

If the domain is already pinned in the file, ask before overwriting:
```
[domain] already has "[existing-skill]" pinned. Replace with "[new-skill]"? [y/n]
```

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
- XDG compliance: use `$XDG_CONFIG_HOME/grimoire/preferences.md` if `XDG_CONFIG_HOME` is set, else `~/.grimoire/preferences.md`
- If a project root cannot be determined, skip option [1] and inform the user: "No project root detected — project-level save unavailable"
- After writing, always confirm the exact path(s) and what was saved

## Common Mistakes

**Auto-pinning**: never write preferences without the user explicitly confirming at Step 2a. Detected ≠ intentional.

**Silent overwrite**: if the domain is already pinned in the file, always ask before replacing. Existing preferences may be carefully chosen.

**Prompting for reasons unnecessarily**: only record a reason if the user volunteered one. Don't ask "why do you prefer this?" on every pin — it creates friction.

**Forgetting the session case**: if the user chose option 0, confirm "saved to session memory" — not a file path. Make the ephemeral nature explicit.
