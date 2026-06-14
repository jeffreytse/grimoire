---
name: plan-best-practice-solution
description: Use when the user's problem spans multiple domains, requires coordinating several best practices, or is too complex for a single skill — e.g. "launch a startup", "handle a workplace medical emergency", "going through a divorce while buying a house."
source: McKinsey Problem Solving (MECE methodology, Rasiel 1999), Kepner-Tregoe problem analysis, Design Thinking (IDEO/Stanford d.school)
tags: [complex-problem, multi-domain, problem-decomposition, mece, skill-orchestration, problem-solver, structured-solution, cross-domain]
---

# Plan Solution

Decompose a complex multi-domain problem into sub-problems, match each to the highest-confidence best-practice skill, sequence them by dependency, and execute one skill at a time.

## Why This Is Best Practice

**Adopted by:** McKinsey, BCG, and Bain use MECE (Mutually Exclusive, Collectively Exhaustive) issue trees as their primary problem-solving methodology for complex client situations. NASA and Boeing use Kepner-Tregoe structured problem analysis for high-stakes multi-system failures. Design Thinking's problem-decomposition phase is standard at Google, IDEO, and Apple before any solution is attempted.
**Impact:** Rasiel (*The McKinsey Way*, 1999) documents that MECE decomposition is the primary tool consultants use to avoid missing problem dimensions — the leading cause of strategy failures. IDEO's Human-Centered Design methodology mandates an explicit problem-definition phase before ideation; teams that skip it produce solutions targeting the wrong problem, requiring full redesign rather than iteration (IDEO HCD Field Guide, 2015). Kepner-Tregoe structured problem analysis is used by Boeing and NASA for high-stakes multi-system failures precisely because unstructured analysis under pressure collapses to the most visible symptom, not the root cause.
**Why best:** Single-skill application fails for complex problems because skills are scoped to one concept. Without decomposition: (1) important problem dimensions are missed; (2) skills are applied in the wrong order — applying a hiring plan before validating a business model wastes work; (3) conflicting guidance across domains goes unresolved. Structured decomposition first, then sequenced skill application, eliminates all three failure modes. Ad-hoc "try one skill, see what sticks" is the alternative — it produces incomplete coverage and rework.

Sources: Rasiel, *The McKinsey Way* (1999); Kepner-Tregoe Problem Analysis; IDEO Human-Centered Design Field Guide (2015)

## Steps

### 1. Extract problem dimensions (silent)

From the user's input, silently identify:

| Dimension | Extract |
|-----------|---------|
| **Goal** | What outcome does the user want? |
| **Domains** | Which fields are involved? (engineering, law, finance, health, etc.) |
| **Constraints** | Time, resources, legal, technical limits |
| **Dependencies** | What must be true before other things can happen? |
| **Urgency** | Is any sub-problem time-sensitive or blocking others? |

Do not ask the user for any of this — infer from what they wrote.

**Problem clarity check:** After extracting dimensions, apply skill judgment: can the goal and at least 2 domains be identified from what the user said? If the goal is completely uninferable, or what's described is clearly a symptom with no root cause context → invoke `analyze-best-practice-problem` first. Use the problem space map from its output to populate the dimensions above, then continue to Step 2.

**Complexity check:** If only one domain is involved and the problem maps cleanly to a single skill, delegate to `suggest-best-practice` instead. `plan-best-practice-solution` is for genuinely multi-domain or multi-step problems.

**Insufficient domains:** If the goal is inferable but fewer than 2 distinct domains can be identified, AND the input signals a multi-domain problem (mentions multiple stakeholders, multiple life areas, or cross-functional concerns), ask ONE question before proceeding to Step 2:
"What areas does this touch? (e.g., legal, technical, financial, operational, health, personal)"
Use the answer to populate the Domains row before decomposing.

### 2. Decompose into MECE sub-problems

Apply MECE decomposition:
- Each sub-problem addresses one distinct dimension of the overall problem
- Sub-problems don't overlap — a skill that solves sub-problem A doesn't also solve sub-problem B
- Together they cover the full problem — no important dimension omitted

Maximum 7 sub-problems. If more emerge, group related ones under a shared theme.

### 3. Match skills to sub-problems

For each sub-problem, score all candidate skills:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Classify the result per sub-problem:

| Result | Condition | Action |
|--------|-----------|--------|
| **Clear match** | 1 skill ≥ 0.7, second < 0.4 | Assign directly — no user choice needed |
| **Multiple candidates** | 2+ skills ≥ 0.4 | Mark for user choice — record all candidates, flag ★ recommendation (highest score) |
| **No match** | All skills < 0.4 | Flag with ⚠ — manual research needed |

If no installed skill covers a sub-problem, flag it explicitly:
`⚠ [sub-problem description]: no installed skill — manual research needed`

### 4. Sequence by dependency

Order skills by logical dependency — not arbitrary order:

| Dependency rule | Example |
|-----------------|---------|
| Validate before build | Business model before go-to-market |
| Legal before commitment | Review contract before signing or hiring |
| Diagnose before fix | Root cause before solution design |
| Calculate before plan | Unit economics before funding strategy |
| Foundation before structure | Architecture before implementation |

Skills with no prerequisites go first. Skills whose output feeds another skill go next.

### 5. Present the solution plan

Present the full sequenced plan. For sub-problems with a clear match, show directly. For sub-problems with multiple candidates, show all options inline with ★ recommendation and collect the user's choice before execution begins.

```
Your situation spans [N] domains. Here is the solution plan:

1. [skill-name] — [what sub-problem it solves]
   Domain: [domain/subdomain]

2. Multiple practices apply — choose one:
   ★ [top-skill] — [one sentence: what it solves]  ← recommended
      [second-skill] — [one sentence: what it solves]
      [third-skill] — [one sentence: what it solves]

3. [skill-name] — [what sub-problem it solves]
   Domain: [domain/subdomain]

⚠ [sub-problem]: no installed skill — manual research needed.
```

After presenting, if any steps have multiple candidates, collect user choices using
the best available method for your platform:

- **Claude Code**: use `AskUserQuestion` — one question per ambiguous step,
  ★ recommended option first with "(Recommended)" appended, `multiSelect: false`
- **Gemini CLI**: use `ask_user` — same structure, `type: "select"`, options list,
  ★ recommended first
- **OpenCode**: use `question` — same schema as `AskUserQuestion`
- **All other platforms** (Codex, Cursor, Copilot, etc.): present numbered list
  and wait for user to type a number or skill name:
  ```
  Step N has multiple options:
  1. [top-skill] ★ (recommended) — [what it solves]
  2. [second-skill] — [what it solves]
  Which would you like? (Enter number or skill name)
  ```

Collect all choices before starting execution. Only proceed once every step has a decided skill.

### 6. Execute one skill at a time

For each skill in the sequence (using user-decided skills from Step 5):
1. Announce: `Applying step N: [skill-name]`
2. Load and run the skill fully
3. After completion, silently proceed to the next step

**Failure-handling rule:** If a skill in the plan sequence fails to complete (errors, produces no output, or the user abandons it mid-way), do not silently skip to the next skill. Instead:
1. Show what failed: `[skill-name] did not complete — [reason if known]`
2. Offer: `Retry [skill-name], skip it, or stop the plan? [retry / skip / stop]`
3. If skipped: note in the final plan summary that [skill-name] was skipped and its dependencies may be incomplete

A plan with silent failures is worse than a failed plan — the user thinks they applied all practices but they didn't.

Only pause if the plan changes — new constraints revealed, a step becomes unnecessary, or a dependency changes:
```
Step N revealed [new constraint]. Revised plan: removing step M, adding [skill-name].
Continue with revised plan?
```

### 7. Adapt after each step

After each skill completes, reassess:
- Did the output reveal new constraints or sub-problems?
- Does the remaining sequence still make sense?
- Are any later skills now unnecessary?

State any changes explicitly before continuing:
```
Step 2 revealed [new constraint]. Adjusting: removing step 4, adding [skill-name] before step 3.
Continue?
```

### 8. Execute the plan

Invoke each skill in the plan sequence in the order shown in Step 5. Proceed to the next skill as soon as the current one completes — no confirmation needed between skills unless the plan changes.

Follow the failure-handling rule: if any skill fails to complete, show what failed, offer retry/skip/stop, and note any skipped skill in the final summary.

## Rules

- If the problem is single-domain and maps to one skill, defer to `suggest-best-practice` — don't over-engineer
- Only pause between skills if the plan changes — not between every skill
- Never hallucinate skill names — only reference skills that exist in installed grimoire domains
- Flag sub-problems with no matching skill explicitly — don't skip them silently
- State the reason for sequencing decisions — don't just present an order without explaining why
- Maximum 7 sub-problems — group if more emerge
- If a sub-problem is itself complex and single-domain, delegate to `apply-best-practice-tree` for recursive drill-down rather than forcing it into the flat sequence
- **Boundary with apply-best-practice-tree:** Use `plan-best-practice-solution` when the problem spans multiple independent domains (law + engineering, or finance + marketing). Use `apply-best-practice-tree` when the problem is complex but stays within one domain. The distinction: cross-domain problems need coordination between separate skill chains; single-domain complex problems need recursive decomposition within one chain.

## Examples

> Skill names in examples are illustrative — actual skills depend on what domains are installed. If a skill is not installed, `plan-best-practice-solution` flags it with ⚠ and notes manual research is needed.

**Example 1 — Startup launch**
> "I want to launch a SaaS startup"

Sub-problems: business model, unit economics, legal structure, go-to-market, hiring
Sequence: `design-business-model` → `calculate-unit-economics` → `review-saas-contract` → `design-go-to-market` → `plan-hiring`
Reason: validate model and economics before legal commitments; legal before hiring

---

**Example 2 — Career transition**
> "I'm a senior engineer who wants to move into engineering management"

Sub-problems: skills gap assessment, compensation negotiation, leadership approach, personal brand
Sequence: `audit-technical-debt` → `negotiate-compensation` → `design-onboarding-program` → `write-leadership-principles`
Reason: understand current position before negotiating; negotiate role before starting; onboarding approach before managing

---

**Example 3 — Single domain → delegate**
> "My pull requests keep getting rejected"

One domain, one clear skill match — delegate: "Single-domain problem. Applying `suggest-best-practice`..."

## Common Mistakes

**Over-applying to simple problems**: one domain, one clear skill → use `suggest-best-practice`. Reserve this skill for problems that genuinely span multiple fields.

**Ignoring dependencies**: a flat unsequenced list creates rework. Always explain the order.

**Hallucinating skills**: if no skill covers a sub-problem, say so. Don't invent names.

**Chaining without confirmation**: running multiple skills back-to-back overwhelms. One at a time.
