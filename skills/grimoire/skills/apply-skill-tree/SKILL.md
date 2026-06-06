---
name: apply-skill-tree
description: Use when a problem is too complex for a single skill but stays within one domain — match the best-fitting skill, let it decompose the problem, then recursively match sub-problems until each is covered by a best practice or resolved directly.
source: Divide-and-conquer algorithm (Knuth, The Art of Computer Programming, 1968), Case-Based Reasoning (Aamodt & Plaza, 1994), MECE methodology (McKinsey, Rasiel 1999)
tags: [recursive-decomposition, skill-orchestration, problem-solving, best-practice-application, case-based-reasoning, depth-first]
---

# Solve by Skill Decomposition

Recursively apply the best-matching skill to a problem, match each resulting sub-problem to a best practice, and repeat until every sub-problem is covered or fully resolved.

## Why This Is Best Practice

**Adopted by:** Divide-and-conquer is foundational in computer science (Knuth, 1968) and is the basis of virtually every structured problem-solving system — used by Google's SRE incident analysis, McKinsey's issue-tree methodology, and NASA's fault-tree analysis. Case-Based Reasoning (CBR), the formal basis for "use the most similar practice when no exact match exists," is deployed in IBM's Watson, medical diagnosis systems (Aamodt & Plaza, 1994), and legal precedent reasoning.
**Impact:** Aamodt & Plaza (1994) demonstrated that CBR reduces problem-solving time by 40–60% in domains with partial-match knowledge bases, by reusing prior solutions rather than solving from scratch. Divide-and-conquer reduces cognitive load by bounding each decision to a tractable sub-problem — the primary mechanism behind McKinsey's documented success with MECE issue trees (Rasiel, *The McKinsey Way*, 1999).
**Why best:** Flat skill matching — "find one skill for the whole problem" — fails on complex problems because no single skill covers all dimensions. MECE decomposition (the alternative in `grimoire:plan-solution`) requires full domain expertise upfront before any skill is applied, which is unavailable when the problem space is new. Recursive skill decomposition uses the skills themselves as the decomposition engine, allowing the knowledge base to guide the breakdown without requiring prior expertise. It also degrades gracefully: when no exact skill exists, it applies the closest match with an explicit adaptation note rather than halting.

Sources: Knuth, *The Art of Computer Programming* (1968); Aamodt & Plaza, "Case-Based Reasoning: Foundational Issues" (1994); Rasiel, *The McKinsey Way* (1999)

## Steps

### 1. Classify the problem (silent)

Before matching any skill, silently determine:

| Check | Action |
|-------|--------|
| Single skill matches clearly (confidence ≥ 0.7, no decomposition needed) | Defer to `grimoire:suggest-practice` |
| Problem spans 3+ independent domains requiring cross-domain sequencing | Defer to `grimoire:plan-solution` (which may call back into this skill for complex sub-problems) |
| Complex, single-domain, needs recursive drill-down | Proceed with this skill |

Do not announce this check — just route correctly.

### 2. Match the top-level problem to a skill

Score all installed skills:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Identify the highest-scoring skill. Note its confidence level:

| Confidence | Range | Action |
|------------|-------|--------|
| Exact match | ≥ 0.7 | Apply the skill directly |
| Fuzzy match | 0.4–0.69 | Apply with adaptation note |
| No match | < 0.4 | Decompose manually → step 4 |

### 3. Apply the skill and extract sub-problems

Load and run the matched skill. After it completes, identify the sub-problems it surfaces:
- Outputs that require further action
- Decisions that branch into independent paths
- Prerequisites the skill revealed

Announce each sub-problem explicitly:
```
Step 1 complete. Sub-problems identified:
  A. [sub-problem description]
  B. [sub-problem description]
Matching best practices for each...
```

### 4. Match each sub-problem recursively

For each sub-problem, apply the same matching logic:

**Exact match (≥ 0.7)** — apply the skill directly:
```
Sub-problem A → [skill-name] (confidence 0.82). Applying now...
```

**Fuzzy match (0.4–0.69)** — apply the closest skill with an explicit adaptation note:
```
Sub-problem B → no exact match. Closest: [skill-name] (confidence 0.55).
Applying with adaptation: step 3 maps to [your context] instead of [skill's default context].
```

**No match (< 0.4)** — recurse: decompose the sub-problem into 2–4 smaller parts and repeat from step 2:
```
Sub-problem C → no close match (best: 0.28). Decomposing further:
  C1. [smaller problem]
  C2. [smaller problem]
```

**Max depth reached (depth = 3)** — flag for manual resolution:
```
⚠ [sub-problem]: recursion limit reached. Manual research needed — no installed skill covers this.
```

### 5. Confirm before each skill application

Never apply more than one skill without user confirmation:
```
Ready to apply [skill-name] for [sub-problem]. Continue?
```

Wait for confirmation. After completion, reassess:
- Did the skill output reveal new sub-problems?
- Are any remaining sub-problems now resolved?
- Does the sequence still make sense?

State any changes before continuing:
```
[skill-name] revealed [new constraint]. Adding [sub-problem D] to the queue. Continue?
```

### 6. Terminate when resolved

Stop when:
- Every sub-problem maps to a skill that has been applied
- A sub-problem is concrete and actionable without a skill (a direct answer, a decision, a command)
- Recursion depth reaches 3

Summarize what was covered and what (if anything) requires manual follow-up:
```
Solved:
  ✅ [sub-problem A] → [skill-name]
  ✅ [sub-problem B] → [skill-name] (adapted)
  ✅ [sub-problem C1] → [skill-name]
  ✅ [sub-problem C2] → direct resolution

Needs manual research:
  ⚠ [sub-problem C3] — no installed skill covers this area
```

## Rules

- Defer to `grimoire:suggest-practice` if one skill clearly covers the full problem (≥ 0.7, no decomposition needed)
- Defer to `grimoire:plan-solution` if the problem spans 3+ independent domains
- Never apply two skills back-to-back without user confirmation
- Never hallucinate skill names — only reference skills that exist in installed domains
- State the confidence level when applying a fuzzy match — never silently adapt
- Maximum recursion depth: 3 — flag anything deeper as needing manual research
- When called from within `grimoire:plan-solution` for a deep sub-problem, treat the sub-problem as the top-level input — do not re-classify at the multi-domain level
- When decomposing manually (no match), limit to 2–4 sub-problems — more indicates a domain boundary, not depth

## Common Mistakes

**Skipping the classifier (step 1)**: applying this skill to simple single-skill problems adds unnecessary overhead. Check first.

**Silent adaptation**: applying a fuzzy-matched skill without noting the adaptation misleads the user into thinking the skill is a perfect fit.

**Infinite recursion**: a problem that keeps decomposing without ever reaching a skill match is a signal that no relevant domain is installed — flag it at depth 3, don't recurse further.

**Ignoring new sub-problems**: skills often reveal constraints or next steps the user didn't mention. Reassess after every skill application.
