---
name: resolve-best-practice-conflict
description: Use when two or more installed grimoire skills give contradictory guidance for the same situation, or when the user wants to scan their preferences file proactively to find and rank conflicting skills before a conflict arises.
source: Preference resolution hierarchy pattern (XDG Base Directory Specification, VS Code settings precedence); priority-based conflict resolution (RFC 2119 precedence model)
tags: [conflict-resolution, preferences, priority, skill-ranking, meta, user-choice, multi-skill, context-qualifier]
---

# Resolve Best Practice Conflict

When two active skills disagree, record which wins — once — so the same conflict never surfaces again.

## Why This Is Best Practice

**Adopted by:** Every major developer tool that supports multiple rule sources implements an explicit precedence model: ESLint (config cascade with override), VS Code (workspace → user → default settings), Git (repo → global → system config), TypeScript (tsconfig extends). Conflict-without-precedence is universally rejected at scale because it forces the same decision repeatedly.
**Impact:** Without explicit priority, users re-resolve the same skill conflict every session. The XDG Base Directory Specification (freedesktop.org) documents that absent a resolution order, tools produce inconsistent behavior across invocations — the primary cause of "why did it do that?" confusion in multi-config environments. RFC 2119 precedence modeling (IETF) is the canonical pattern for deterministic rule application in multi-source systems.
**Why best:** The alternative — applying whichever skill the AI happens to retrieve first — produces non-deterministic behavior. Explicit user-defined priority is deterministic, auditable, and revocable. Context qualifiers narrow precedence to the exact situations where the user actually wants it, avoiding over-application of a blanket rule.

Sources: XDG Base Directory Specification (freedesktop.org); RFC 2119 (IETF); VS Code settings precedence documentation; ESLint configuration cascade documentation

## Steps

### Step 1: Identify the conflict

**Mode A — explicit:** User names two skills that conflict (e.g., "SOLID and KISS both apply here — which wins?").

**Mode B — proactive scan:** User asks to scan their preferences file for conflicts, OR you are about to apply two skills and notice they give contradictory guidance.

**Mode B trigger:** This skill also applies when a user observes two practices being applied inconsistently in a project (e.g., 'sometimes we use X, sometimes Y'). The conflict need not be explicitly stated — if the user describes inconsistency or asks 'should we use X or Y?', treat it as a conflict to resolve.

For Mode B, load the active settings files and merge in resolution order:

```
<project>/.grimoire/settings.local.toml  (highest — project personal)
  > <project>/.grimoire/settings.toml   (project shared)
    > ~/.config/grimoire/settings.toml  (lowest — global)
```

Within any file, cascade by specificity: `specific entry in [domain.subdomain] > [domain.subdomain] > [domain] > [global]`

Merge from least to most specific — each level overrides the level below for the same key. `[engineering.architecture]` overrides `[engineering]` which overrides `[global]`.

Before scanning each domain:
- **Resolve active profile** — if user has activated a named profile for this domain, load `[domain.subdomain.profiles.name]` instead of the untagged `[domain.subdomain]` default
- **Check `lock:`** — if set, ignore all less-specific sections (`## domain`, `## global`) for this subdomain; it's self-contained
- **Check `remind:`** — if today ≥ remind date, say "Reminder: review [domain] preferences (set on [date])." Don't block; continue scanning.
- **Check `expires:`** — if today ≥ expires date, warn: "Preferences for [domain] expired on [date] — review before scanning?"
- **Check `skip-if:`** — if current file/context matches any skip-if pattern, skip this domain entirely; don't scan or apply
- **Skip `disabled:` skills** — don't load or analyze them; they're intentionally off
- **Flag `require:` contradictions** — if a skill appears in both `require:` and `disabled:`, flag as contradiction before continuing

**Contradiction rubric:** Two practices are in conflict (not complementary) when applying both simultaneously would require contradictory actions:
- Same trigger, opposite prescriptions (e.g., 'always add types' vs 'prefer implicit types for brevity')
- Mutually exclusive architectural decisions (e.g., 'prefer composition' vs 'prefer inheritance for domain hierarchy')
- Contradictory tool or process requirements for the same step

**Complementary, not conflicting:** practices that apply to different parts of the same domain (e.g., 'write unit tests' + 'write integration tests') or at different lifecycle stages. Do not raise as a conflict.

For each domain (non-disabled skills), compare the "Steps" and "When NOT to Use" sections of co-listed skills. Flag pairs where one skill says "always do X" and another says "avoid X" or "do Y instead." Remember: all listed skills apply — only direct contradictions need resolution.

---

### Step 2: Show the conflict

Load both SKILL.md files. Summarize the specific contradiction in one line each:

```
Conflict: engineering/architecture

  apply-solid-principles  → "Each class should have one responsibility; split
                             classes that do more than one thing."
  apply-kiss-principle    → "Don't split a function when it's still explainable
                             in one sentence — splitting adds indirection."

These conflict when deciding whether to extract a single-responsibility class
that is currently simple.
```

If multiple conflicts found in Mode B, present all at once, then resolve in sequence.

---

### Step 3: Show current priorities (if any)

If the relevant domain section already exists in `settings.toml`, display it:

```
Current priorities for engineering/architecture:
  1. Google's engineering practices
  2. SOLID principles
  (no fallback set — default: higher-ranked skill governs silently)
```

If no priorities exist yet for that domain, say so.

---

### Step 4: Ask which wins

If existing priority already found in Step 3, use it silently — skip this step.

If user hasn't already stated a preference:

```
Which should take priority?
  [1] apply-solid-principles — when?  (all contexts / specific context)
  [2] apply-kiss-principle   — when?  (all contexts / specific context)
  [3] Both apply — I want the AI to note the tension and apply both
  [4] Ask me each time

You can narrow either choice to a context (e.g., "SOLID for production code, KISS for scripts").
```

If user already stated a preference earlier in the conversation, use it directly — don't re-ask.

---

### Step 5: Update settings.toml

**Duplicate-write check:** Before writing the resolution, check if the exact same domain+practice pair already has a pinned preference. If yes, show the existing pin and confirm: 'This conflict was already resolved: [existing-pin]. Overwrite with new decision? [y/n]'. Do not silently re-pin.

Write the priority into the correct domain section using TOML. Default to `<project>/.grimoire/settings.toml` (shared) — ask if user wants `settings.local.toml` (personal) or global (`~/.config/grimoire/settings.toml`).

```toml
[engineering.architecture]
practices = [
  "SOLID principles: production code",  # index 0 = higher priority; context qualifier after ':'
  "KISS: prototypes, scripts"           # index 1 = lower priority
]
fallback = "ask"                        # optional; only needed if user chose [4]
```

Rules:
- **All listed skills apply.** Array index = priority for resolving direct contradictions. When two skills conflict on a specific point, lower-index skill wins. Non-conflicting guidance from all skills applies regardless of index.
- Context qualifier as `: qualifier` suffix in the practices string. Omit for all-context preference.
- `fallback = "ask"` — AI prompts when skills directly contradict and no priority is set.
- `fallback = "both"` — AI applies both and surfaces the tension.
- Omit `fallback` for default (lower-index skill wins conflict silently).
- `disabled` — array of skill names the AI must never apply in this domain even if installed.
- `ask-before` — array of skill names AI must confirm before applying. Use for high-stakes skills (security, DB migrations, destructive operations).
- `expires` — ISO date string. After this date AI warns before applying. Use for temporary overrides.
- `remind` — ISO date string. Soft nudge on/after date; doesn't block.
- `shared = true` — commit `settings.toml` to repo as team standard.
- `lock = true` — less-specific sections (`[domain]`, `[global]`) do not contribute to this section.
- `note` — string. AI reads when applying domain skills. Cite when user asks why.
- `require` — array. Skills that MUST apply unconditionally. Win over all conflict resolution. Flag if a skill appears in both `require` and `disabled` — contradiction.
- `author` — string. Who set this preference. Cite when user asks "why do we prefer X?" Include in Step 6 if present.
- `skip-if` — array of file/context patterns. Skip entire domain when matched.
- `[domain.profiles.name]` — named profile variant. Inactive until user activates. On activation replaces the untagged default for that domain. Ask whether to create as default or named profile when writing.
- If the domain section already exists, merge conflicting skills into `practices` at the correct index — don't replace other keys.

---

### Step 6: Confirm

Show the updated domain section and the file path:

```
Updated: .grimoire/settings.toml

[engineering.architecture]
practices = [
  "SOLID principles: production code",
  "KISS: prototypes, scripts"
]
fallback = "ask"

Next time both skills apply, SOLID governs in production contexts.
KISS governs in prototype/script contexts. Tied contexts → I'll ask.
```

If any optional keys were written, confirm them explicitly:

```
Also recorded:
  require: apply-input-validation     → always applied, unconditional
  ask-before: apply-solid-principles  → I'll confirm before applying SOLID
  skip-if: test files                 → skipping this domain for test files
  author: @backend-team               → attributed to @backend-team
  expires: 2026-09-01                 → preferences expire on 2026-09-01
  remind: 2026-07-01                  → review reminder set for 2026-07-01
  lock: true                          → less-specific sections won't contribute
  note: legacy codebase — avoid new OOP patterns
```

---

### Proactive scan summary (Mode B)

After scanning and resolving all conflicts:

```
Scanned 12 skills across 4 domains. Found 3 conflicts:

  ✅ engineering/architecture  → SOLID > KISS (resolved above)
  ✅ engineering/testing       → TDD > test-after (resolved above)
  ⚠️  engineering/development  → YAGNI vs DRY — 1 conflict found, skipped (add to queue?)

Updated: .grimoire/settings.toml
```

## When NOT to Use

- **Same-domain, non-conflicting skills**: SOLID and Law of Demeter are complementary, not conflicting. Don't rank them unless an actual contradiction arises.
- **Skills in different domains**: KISS (engineering) and KISS-equivalent (writing/simplicity) operate independently — no cross-domain resolution needed.
- **Only one skill applies**: if context clearly selects one skill over the other without conflict, apply it directly.

## Common Mistakes

**Over-resolving**: not every pair of skills that sounds related actually conflicts. Load the SKILL.md files and identify a specific contradiction before treating it as a conflict.

**Global-scoping everything**: default to project-level resolution (`settings.toml`). Only write to `~/.config/grimoire/settings.toml` when the user explicitly wants the preference across all projects.

**Forgetting context qualifiers**: "SOLID wins" as a blanket rule overrides KISS even in script/prototype contexts where KISS would be correct. Always ask whether the priority is conditional.
