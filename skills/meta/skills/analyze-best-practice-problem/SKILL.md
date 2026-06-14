---
name: analyze-best-practice-problem
description: Use when the user presents a problem that isn't well-defined — goal is unclear, scope is unstated, or what's described looks like a symptom rather than a root cause. Asks clarifying questions one at a time until the problem boundary is clear, then produces a problem statement, problem space map, and possible solution routes before handing off to solution skills.
source: Design Thinking problem framing (IDEO/Stanford d.school); McKinsey issue tree methodology (Rasiel, 1999); Meadows "Thinking in Systems" (2008) — system boundary definition; Toyota 5 Whys (Taiichi Ohno)
tags: [problem-framing, analysis, clarification, scoping, root-cause, problem-definition, vague-problem]
---

# Analyze Best Practice Problem

Clarify an ill-defined problem through structured questioning, then map the problem space and surface possible routes before handing off to solution skills.

## Why This Is Best Practice

**Adopted by:** Design Thinking's problem-definition phase (IDEO/Stanford d.school) is standard at Google, Apple, and IDEO before any solution is attempted — teams that skip it produce solutions targeting the wrong problem, requiring full redesign rather than iteration (IDEO HCD Field Guide, 2015). McKinsey's issue tree methodology (Rasiel, 1999) treats problem definition as a separate, mandatory step before any analysis begins. Donella Meadows (*Thinking in Systems*, 2008) identifies boundary definition — deciding what is inside vs. outside the system being analyzed — as the single most consequential decision in systems analysis, and the most commonly skipped.
**Impact:** The Institute of Medicine (*To Err is Human*, 1999) found that the majority of medical errors trace to acting on a symptom without diagnosing the root cause — the most expensive form of the "wrong problem" failure. In software engineering, the majority of rework traces to misunderstood requirements, not implementation errors (Standish Group CHAOS Report). In both fields: solving the right problem poorly produces less waste than solving the wrong problem well.
**Why best:** Practitioners who ask clarifying questions before solving report fewer "solution rejection" cycles — where a completed solution fails because it addressed the wrong aspect of the problem. The one-question-at-a-time constraint (vs. a batch of questions) is deliberate: batched questions overwhelm users and produce incomplete answers. Sequential questions surface progressively deeper information, with each answer informing the next question.

Sources: IDEO HCD Field Guide (2015); Rasiel, *The McKinsey Way* (1999); Meadows, *Thinking in Systems* (2008); Institute of Medicine (1999) *To Err is Human*; Standish Group CHAOS Report

## Steps

### Step 0: Match domain-specific analysis skill (silent)

Before running the default analysis steps, check if an installed skill is specifically designed for problem analysis in this domain.

Score all installed skills against the user's input using:
```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Filter to skills with at least one of these tags:
`problem-analysis`, `root-cause`, `root-cause-analysis`, `diagnostic`, `issue-tree`, `problem-framing`, `systems-thinking`, `failure-analysis`, `causal-analysis`, `decision-analysis`

**If one skill scores ≥ 0.5:**
Announce and invoke it:
```
This problem matches a domain-specific analysis practice: [skill-name] ([domain/subdomain]).
Applying it to define the problem before we look for solutions...
```
After it runs: if it produced a problem statement and/or structured map, use those as output. If the invoked skill also runs a clarification loop, let it — do not duplicate questions.

**If 2+ skills score ≥ 0.5:**
Present a compact ranked list and wait for selection:
```
Multiple problem analysis practices apply:
  ★ [top-skill] — [one sentence: what analysis it provides]  ← recommended
     [second-skill] — [one sentence]
```

Then ask the user using the best available method for your platform:
- **Claude Code**: use `AskUserQuestion` — question: "Which analysis practice should I use first?", ★ recommended first with "(Recommended)" appended, `multiSelect: false`
- **OpenCode**: use `question` — same schema as `AskUserQuestion`
- **Gemini CLI**: use `ask_user` — question: "Which analysis practice should I use first?", `type: "select"`, ★ recommended first
- **All other platforms**: numbered list:
  ```
  1. [top-skill] ★ (recommended) — [what analysis it provides]
  2. [second-skill] — [what analysis it provides]
  Which should I use to analyze this problem first? (Enter number or name)
  ```

**If no skill scores ≥ 0.5:**
Proceed with default Steps 1–6 below. No announcement needed.

**If a domain skill ran and produced a problem statement:** skip Step 1 — trust the domain skill's output. Use its problem statement and map directly in Step 4.

**Domain skill lookup:** If a domain skill ran before this skill and produced a structured problem statement, use its output directly — skip Step 1 (re-elicitation). The domain skill's problem description is authoritative. This avoids asking the user to re-describe a problem already captured by a domain-specific skill.

### Step 1: Assess problem definition

From what the user said, check if all three criteria are satisfiable without asking:

- [ ] **Goal** — what outcome does solving this achieve?
- [ ] **Scope** — what's included in the problem? What's out?
- [ ] **Root cause vs symptom** — is what's described the actual problem, or a downstream effect?

If all three can be inferred → skip to Step 4 (problem space map). Do not ask unnecessary questions.
If any are missing → proceed to Step 2 (clarification loop).

### Step 2: Clarification loop

Ask ONE question per turn, in this priority order:

1. **Goal first** — "What outcome are you trying to achieve?"
2. **Scope** — "What's in scope for this? What are you explicitly not trying to solve?"
3. **Root cause check** — "Is [the thing described] the problem itself, or a symptom of something deeper?"
4. **Constraints** (only if still ambiguous after above) — "Any hard constraints I should know about — time, budget, team size, or tech?"

Stop the loop when all three criteria (goal, scope, root cause) are satisfiable from the conversation. Constraints are optional — include if clearly relevant to route selection.

Wait for the user's answer before asking the next question. Do not ask multiple questions in one turn.

### Step 3: Surface if symptom

If what's described is clearly a symptom of a deeper problem, surface it before mapping:

```
This looks like a symptom of [inferred root cause].

Real problem: [root cause]?
Or do you need to solve [stated symptom] specifically — for example, a quick fix while the root cause is addressed separately?
```

Wait for the user's answer before proceeding to the map.

### Step 4: Problem space map

Output the structured map:

```
Problem statement: [1-2 sentences — crisp goal + situation]

Scope
  In:  [what's included]
  Out: [what's explicitly excluded — or "not stated; assumed full scope"]

Constraints
  [type]: [value]
  (omit this section if none stated or inferable)

Stakeholders: [who is affected or must approve — omit if not relevant]

Root cause: [root cause statement — or "as stated" if already root-level]
```

### Step 5: Possible routes

Present 2–4 high-level solution directions. These are directions, not specific practices — keep abstract at this stage:

```
Possible routes:

A. [Route name] — [one line: what this direction does] | Trade-off: [gain / cost]
B. [Route name] — [one line: what this direction does] | Trade-off: [gain / cost]
C. [Route name] — [one line: what this direction does] | Trade-off: [gain / cost]
```

Common route types to adapt to domain context:
- Fix root cause (addresses the source; slower, permanent)
- Workaround symptom (fast relief; temporary, accumulates debt)
- Prevent recurrence (proactive; investment up front, eliminates future occurrences)
- Restructure to eliminate the problem class (architectural; highest effort, highest leverage)
- Delegate / outsource (shift ownership; works when it's not a core competency)

**Step 5→6 handoff:** After outputting the possible routes above, this skill's analytical work is complete. Do not proceed to score or suggest practices — that is `suggest-best-practice`'s job. Return the problem statement and routes as output. The calling skill (suggest-best-practice or the user) uses them as input to their next step.

### Step 6: Handoff

```
Which route fits best? I can then find the applicable best practices or build a full solution plan.
(Say the letter — or "practices" to see matching skills, or "plan" for a sequenced action plan)
```

- User says a route letter → apply the same handoff routing based on direction chosen
- "practices" → route to `suggest-best-practice` with the problem statement as input
- "plan" → route to `plan-best-practice-solution` with the problem space map as input

**If invoked automatically by another skill** (not by the user directly): skip Step 6. Return the problem statement and map as output to the invoking skill.

## Rules

- Never ask more than one question per turn
- Don't proceed to the map until all three criteria (goal, scope, root cause) are satisfiable from the conversation
- Routes are high-level directions, not skill names — keep abstract until handoff
- If the problem is already well-defined (Step 1 YES), go directly to Step 4 — do not run the clarification loop
- If invoked automatically by another skill, skip the handoff offer and return structured output instead
- Constraints are not required to proceed — only ask if they're clearly relevant to route selection
- The symptom check (Step 3) is optional — only surface it when there's clear evidence the described thing is downstream of a deeper cause

## Common Mistakes

**Asking multiple questions at once**: "What's your goal, scope, and timeline?" overwhelms the user and produces incomplete answers. One question per turn — the answer informs the next.

**Running the clarification loop when the problem is already clear**: if goal, scope, and root cause are all inferable, go straight to the map. Unnecessary questions feel patronizing.

**Abstract routes that don't help decide**: "Option A: technical approach / Option B: non-technical approach" isn't useful. Name the actual direction and its trade-off.

**Skipping the symptom check**: solving a symptom without flagging the root cause creates a false sense of resolution. If there's evidence of a deeper cause, surface it in Step 3.

**Forgetting the invocation context**: when called automatically by another skill, don't present the handoff offer — return structured output so the calling skill can continue.
