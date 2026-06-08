---
name: intercept-best-practice
description: Use when starting any user task, action, or request — before taking any other action. Proactively checks whether an applicable best practice exists and applies or offers it before proceeding.
source: Toyota Production System poka-yoke (Shingo, 1986) — error prevention before the action, not after
tags: [proactive, interception, pre-action, best-practice-discovery, guardian, poka-yoke, quality-gate]
---

# Intercept Best Practice

Before any task begins, silently check whether an applicable best practice exists — apply it proactively if confident, offer it if plausible, proceed silently if no match.

## Why This Is Best Practice

**Adopted by:** Toyota's poka-yoke principle (Shingo, 1986) is the foundation of error prevention in manufacturing — checking conditions *before* an action begins, not after defects appear. This principle is applied across aviation (pre-flight checklists), surgery (WHO Surgical Safety Checklist), and software (pre-commit hooks, CI gates).
**Impact:** The WHO Surgical Safety Checklist reduced surgical complications by 36% and deaths by 47% — not by changing technique, but by systematically intercepting the moment before each critical step (Haynes et al., 2009, NEJM). Pre-commit hooks catch an estimated 30–40% of common defects before they enter the codebase (GitHub internal data).
**Why best:** Reactive systems wait for the user to ask. Proactive interception closes the "I didn't know a best practice existed" gap — the most common reason practices aren't applied. The cost is minimal (a silent confidence check); the benefit is consistent practice application without requiring the user to know what to ask.

Sources: Shingo (1986) "Zero Quality Control: Source Inspection and the Poka-Yoke System"; Haynes et al. (2009) NEJM; GitHub Engineering blog

## Steps

### Step 0: Check preferences (silent)

Resolution order — first match wins:
1. Session memory — pinned this session only; not written to disk (highest precedence)
2. `<project-root>/.grimoire/preferences.md` — project-level
3. `~/.config/grimoire/preferences.md` OR `~/.grimoire/preferences.md` — global-level
4. `CLAUDE.md` `## Grimoire Preferences` section — legacy fallback

For the relevant domain, check if a practice is already pinned:
- **Pinned match (file)** → apply the pinned practice directly; skip scoring entirely. No further action needed — already persisted.
- **Pinned match (session)** → apply the pinned practice directly; skip scoring. After applying, offer once per session per domain:
  `"[practice] is pinned for this session only. Save it for future sessions? [y/n]"`
  If yes, invoke `pin-best-practice-preference` to write it to project or global file.
- **Pinned conflict** → warn before suggesting an alternative:
  `"You have [X] pinned for [domain]. Suggest changing it? [y/n]"`
- **No pin** → proceed to Step 1.

### Step 1: Extract intent (silent)

From the user's request, identify:

| Signal | Extract |
|--------|---------|
| **Task type** | What operation? (create, review, fix, design, write, calculate, plan…) |
| **Domain** | Which field? (engineering, law, finance, health, cooking…) |
| **Subject** | What specifically? (unit test, contract, retirement plan, training program…) |
| **Maturity** | Starting fresh, or improving existing work? |

Do not ask the user anything — infer from the request.

### Step 2: Score candidates (silent)

Score installed grimoire skills using the same model as `suggest-best-practice`:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Cap at top 3 matches. Compute silently — do not announce scoring is happening.

### Step 3: Route by confidence

| Condition | Action |
|-----------|--------|
| 1 skill ≥ 0.7, second < 0.4 | Announce + apply, then complete original task |
| 2+ skills ≥ 0.4 (any scores) | Brief ranked offer with ★ recommendation, wait for choice |
| 1 skill 0.4–0.69, no others ≥ 0.4 | Brief offer: apply or skip |
| All skills < 0.4 | Proceed silently — no announcement, no interruption |

**Sole clear match (1 skill ≥ 0.7, second < 0.4):**
```
Best practice detected: [skill-name] ([domain/subdomain])
Applying before proceeding...
```

**Multiple matches (2+ skills ≥ 0.4):**
```
Multiple best practices apply to this task:
  ★ [top-skill] ([domain]) — recommended
    [second-skill] ([domain])
    [third-skill] ([domain])

Apply ★ [top-skill] first? [y / n / 2 / 3 — or just continue]
```

Keep this offer concise — this is an interception, not a discovery session. One line per option.

**Single medium-confidence match (0.4–0.69):**
```
A best practice exists for this: [skill-name] ([domain/subdomain]).
Apply it first? [y/n — or just continue and I'll proceed with your request]
```

**Low confidence (< 0.4):** No output. Proceed with original request.

### Step 4: Complete original task

After the matched skill applies (or if no match), always complete the user's original request. The skill output becomes context for fulfilling the request — never abandon the task.

If the skill redirects significantly, state:
```
Best practice applied. Proceeding with your original request using the above framework.
```

## Rules

- Never block or announce when confidence < 0.4 — proceed silently
- Never apply two skills back-to-back without user confirmation between each
- Always complete the original request after skill application — don't abandon the task
- No browse mode — this is interception only, not discovery
- No clarifying questions — route on what's available, or pass through silently
- If user explicitly invoked a skill by name, don't intercept — they already know what they want

## Key Differences from `suggest-best-practice`

| | `suggest-best-practice` | `intercept-best-practice` |
|---|---|---|
| Trigger | User describes a problem | User starts any action |
| Low confidence | Asks clarifying question | Proceeds silently |
| After skill | Routes to skill | Applies skill, THEN completes original task |
| Browse mode | Supported | Not supported |
| Interrupts flow | Yes (it IS the flow) | Only at ≥ 0.4 confidence |

## Common Mistakes

**Blocking low-confidence tasks**: if score < 0.4, say nothing and proceed. Interrupting the user for uncertain matches destroys trust.

**Abandoning the original task**: always complete what the user asked for after applying a skill.

**Intercepting explicit skill invocations**: if the user said `/write-unit-test`, they already know what they want. Don't re-intercept.

**Announcing scoring**: "Checking for best practices..." on every action is disruptive. Only speak when there is a match.
