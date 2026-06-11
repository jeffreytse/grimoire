---
name: start-best-practice
description: Use when responding to any user request — maps the intent to the appropriate grimoire workflow, exits silently if no workflow matches, and routes task-starts through skill scoring or other intents to the right entry skill.
source: Toyota Production System poka-yoke (Shingo, 1986) — error prevention before the action, not after
tags: [proactive, interception, pre-action, best-practice-discovery, guardian, poka-yoke, quality-gate]
---

# Intercept Best Practice

Map every user input to the appropriate grimoire workflow — or exit silently if none applies. For task-start inputs, score installed skills and offer the best match before the task begins. For all other workflow intents, route directly to the right skill.

## Why This Is Best Practice

**Adopted by:** Toyota's poka-yoke principle (Shingo, 1986) is the foundation of error prevention in manufacturing — checking conditions *before* an action begins, not after defects appear. This principle is applied across aviation (pre-flight checklists), surgery (WHO Surgical Safety Checklist), and software (pre-commit hooks, CI gates).
**Impact:** The WHO Surgical Safety Checklist reduced surgical complications by 36% and deaths by 47% — not by changing technique, but by systematically intercepting the moment before each critical step (Haynes et al., 2009, NEJM). Pre-commit hooks catch an estimated 30–40% of common defects before they enter the codebase (GitHub internal data).
**Why best:** Reactive systems wait for the user to ask. Proactive interception closes the "I didn't know a best practice existed" gap — the most common reason practices aren't applied. The cost is minimal (a silent confidence check); the benefit is consistent practice application without requiring the user to know what to ask.

Sources: Shingo (1986) "Zero Quality Control: Source Inspection and the Poka-Yoke System"; Haynes et al. (2009) NEJM; GitHub Engineering blog

## Steps

### Step 0: Map intent to workflow (silent)

Silently match the input against the grimoire workflow intent map. This is the only gate — no scoring, no output. If no row matches, exit immediately and proceed with the original request unmodified.

| Intent signal | Route |
|---------------|-------|
| About to start a task — executing, implementing, building, writing, deploying, designing, refactoring, testing, planning a concrete deliverable | **self** — continue to Step 1 |
| Has a problem or situation, unsure which skill fits | invoke `suggest-best-practice` |
| Already has a plan, wants gaps checked | invoke `review-best-practice-fit` |
| Problem spans 3+ independent domains | invoke `plan-best-practice-solution` |
| Complex problem within one domain, wants structured decomposition | invoke `apply-best-practice-tree` |
| Doesn't know what practices exist for a topic | invoke `discover-best-practices` |
| Problem isn't clear yet, needs defining before solving | invoke `analyze-best-practice-problem` |
| Wants to activate a paradigm's practices (OOP, TDD, DDD, etc.) | invoke `apply-best-practice-profile` |
| Wants to align a project/artifact to stated preferences (BPDD) | invoke `apply-best-practice-driven-development` |
| Wants to check if an artifact aligns with stated preferences | invoke `check-best-practice-compliance` |
| Has a specific compliance finding to fix | invoke `fix-best-practice-finding` |
| Two practices conflict | invoke `pin-best-practice-preference` |
| **None of the above** | **exit silently** |

**Non-match signals (always exit silently — do not score, do not route):**
- Primary verb is explanatory: *explain, describe, what is, how does, tell me about, define*
- Conversational acknowledgment: "ok", "thanks", "looks good", "got it", "continue"
- Grimoire meta-queries: "what skills are installed?", "show me available skills"
- Explicit `/skill-name` invocations — user already knows what they want (see Rules)
- Follow-up message in an ongoing execution (not a new intent)

**If ambiguous:** treat as "about to start a task" → continue to Step 1. Missing an intercept is acceptable; mis-routing a real task is not.

**When routing to another skill** (any row except "self" and "exit silently"):
```
This looks like [matched situation]. Routing to [skill-name]...
```
Then invoke that skill and stop — do not continue to Step 1.

---

### Step 1: Check preferences (silent)

Resolution order — first match wins:
1. Session memory — pinned this session only; not written to disk (highest precedence)
2. `<project-root>/.grimoire/preferences.md` — project-level
3. `~/.config/grimoire/preferences.md` OR `~/.grimoire/preferences.md` — global-level
4. `CLAUDE.md` `## Grimoire Preferences` section — legacy fallback

For the relevant domain, check if a practice is already pinned:
- **Pinned match (file)** → apply the pinned practice directly; skip scoring entirely. No further action needed — already persisted.
- **Pinned match (session)** → apply the pinned practice directly; skip scoring. After applying, offer once per session per domain using a platform-aware confirm: "Save [practice] for future sessions?" (Claude Code/OpenCode: `AskUserQuestion`; Gemini CLI: `ask_user type: confirm`; other: `[y/n]`). If yes, invoke `pin-best-practice-preference`.
- **Pinned conflict** → warn before suggesting an alternative using a platform-aware confirm: "You have [X] pinned for [domain]. Suggest changing it?" (Claude Code/OpenCode: `AskUserQuestion`; Gemini CLI: `ask_user type: confirm`; other: `[y/n]`).
- **No pin** → proceed to Step 1.

### Step 2: Extract intent (silent)

From the user's request, identify:

| Signal | Extract |
|--------|---------|
| **Task type** | What operation? (create, review, fix, design, write, calculate, plan…) |
| **Domain** | Which field? (engineering, law, finance, health, cooking…) |
| **Subject** | What specifically? (unit test, contract, retirement plan, training program…) |
| **Maturity** | Starting fresh, or improving existing work? |

Do not ask the user anything — infer from the request.

### Step 3: Score candidates (silent)

Score installed grimoire skills using the same model as `suggest-best-practice`:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Cap at top 3 matches. Compute silently — do not announce scoring is happening.

### Step 4: Route by confidence

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

- **Claude Code**: `AskUserQuestion` — one option per candidate skill, top-ranked marked `(Recommended)`, include "Skip — continue without applying" option, `multiSelect: false`
- **Gemini CLI**: `ask_user` — `type: "select"`, recommended first
- **OpenCode**: `question` — same schema as `AskUserQuestion`
- **All other platforms**: plain text:
  ```
  Multiple best practices apply to this task:
    ★ [top-skill] ([domain]) — recommended
      [second-skill] ([domain])
      [third-skill] ([domain])

  Apply ★ [top-skill] first? [y / n / 2 / 3 — or just continue]
  ```

Keep this offer concise — this is an interception, not a discovery session. One option per candidate skill.

**Single medium-confidence match (0.4–0.69):**

- **Claude Code**: `AskUserQuestion` — two options: "Apply [skill-name] first (Recommended)" and "Skip — continue with my request"
- **Gemini CLI**: `ask_user` — `type: "confirm"`, recommended first
- **OpenCode**: `question` — same schema as `AskUserQuestion`
- **Other**: plain text:
  ```
  A best practice exists for this: [skill-name] ([domain/subdomain]).
  Apply it first? [y/n — or just continue and I'll proceed with your request]
  ```

**Low confidence (< 0.4):** No output. Proceed with original request.

### Step 5: Complete original task

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

| | `suggest-best-practice` | `start-best-practice` |
|---|---|---|
| Trigger | User describes a problem | User starts any action |
| Unrelated inputs | N/A — user invokes directly | Exit silently (Step 0 gate) |
| Low confidence | Asks clarifying question | Proceeds silently |
| After skill | Routes to skill | Applies skill, THEN completes original task |
| Browse mode | Supported | Not supported |
| Interrupts flow | Yes (it IS the flow) | Only at ≥ 0.4 confidence |
| Dispatch role | No | Yes — routes non-task intents to correct workflow skill |

## Common Mistakes

**Blocking low-confidence tasks**: if score < 0.4, say nothing and proceed. Interrupting the user for uncertain matches destroys trust.

**Abandoning the original task**: always complete what the user asked for after applying a skill.

**Intercepting explicit skill invocations**: if the user said `/write-unit-test`, they already know what they want. Don't re-intercept.

**Announcing scoring**: "Checking for best practices..." on every action is disruptive. Only speak when there is a match.
