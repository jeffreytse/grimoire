---
name: suggest-practice
description: Use when the user describes any situation, problem, goal, complaint, or question — including when they want to browse available best practices for a topic, don't know which domain applies, or don't know a best practice exists for their situation.
source: Information retrieval best practices (van Rijsbergen, 1979), Nielsen Norman Group search UX guidelines
tags: [skill-discovery, auto-classify, problem-routing, problem-solver, situation-analysis, practice-recommendation]
---

# Suggest Practice

Accept any natural language input, auto-classify it to the relevant domain(s), and
either apply the best-matching skill directly or present a ranked list for the user
to choose from.

## Why This Is Best Practice

**Adopted by:** Zero-knowledge entry points are standard in expert systems and
recommendation engines — Google Search auto-classifies queries without requiring users
to specify intent type, Spotify Discover Weekly infers genres from listening behavior
without requiring genre tags, and Wolfram Alpha auto-routes computations to the correct
solver domain without user specification.
**Impact:** Requiring users to know the domain before getting help is the primary
barrier to knowledge system adoption. Studies on expert system usability (Prerau, 1990)
found that systems requiring users to pre-classify their problems had 3–5× lower
adoption than systems that accepted natural language and routed automatically.
**Why best:** A single entry point — vs. domain-specific skills that require users to
already know which domain and practice applies to their situation — eliminates the
"which skill do I use?" confusion. Users don't need to pre-classify their problem;
`suggest-practice` infers domain and intent automatically.

Sources: Prerau (1990) expert system usability research, Nielsen Norman Group search
UX guidelines, Google Search intent classification research

## Steps

### 1. Extract intent signals (no clarifying questions yet)

From the user's input, silently identify:

| Signal | Extract |
|--------|---------|
| **Goal** | What outcome does the user want? |
| **Symptoms** | What problems are they experiencing? |
| **Domain cues** | Industry, role, tool names, context words |
| **Constraints** | Time pressure, team size, resources, urgency |
| **Problem type** | Prevention / diagnosis / optimization / compliance / learning |

Do not ask the user for any of this — infer it from what they wrote.

### 2. Score candidate skills

For each skill in the installed grimoire domains:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

- **tag_overlap**: count of skill tags matching extracted keywords (normalized 0–1)
- **description_match**: does `Use when...` describe this situation? (0 = no, 1 = yes)
- **domain_plausibility**: is this domain plausible given context cues? (0 = no, 0.5 = possible, 1 = likely)

Normalize final scores to 0–1.

### 3. Route by confidence

**Existing solution** — user describes something they already built, planned, or decided and wants evaluation ("is this good?", "what am I missing?", "does this follow best practices?"):

Delegate to `review-practice-fit`. Announce:
```
You have an existing solution. Applying review-practice-fit to evaluate it against best practices...
```

**Multi-domain** — user has a new problem that spans 3+ domains and requires applying ALL of them (not just one):

Delegate to `plan-solution`. Announce:
```
Your situation spans multiple domains and requires coordinating several best practices.
Applying plan-solution to build a sequenced action plan...
```

**High confidence** — one skill scores ≥ 0.7 and clearly leads the rest:

Load and apply the skill directly. Announce before applying:
```
Situation matches: [skill-name] ([domain/subdomain])
Applying now...
```

**Multiple plausible matches** — 2–5 skills with similar scores (top score < 0.7,
or top 2 scores within 0.15 of each other), but from the SAME domain or the user
only needs ONE of them:

Present a ranked list and wait for user selection:
```
Your situation matches several best practices. Which fits best?

1. [skill-name] — [one sentence: what problem it solves]
   Domain: [domain/subdomain]  |  Install: /plugins add github:jeffreytse/grimoire/[path]

2. [skill-name] — [one sentence]
   Domain: [domain/subdomain]  |  Install: /plugins add github:jeffreytse/grimoire/[path]
```
After user selects, load and apply the chosen skill.

**Browse mode** — user explicitly says "show me options", "what practices exist for X",
or "what should I know about Y" without wanting to act yet:

Present the ranked list only, do not apply any skill:
```
Best practices for: [topic]

1. [skill-name] — [one sentence: what it solves]
   Domain: [domain/subdomain]  |  Install: /plugins add github:jeffreytse/grimoire/[path]

2. ...
```
Announce at end: "Say the number or skill name to apply one."

**No match** — all skills score < 0.3:

State clearly that no skills currently cover this area, then ask ONE targeted
clarifying question to narrow the domain:
```
No installed skills match this situation closely.
[One clarifying question to narrow the domain — e.g., "Is this about your code, your health, your finances, or something else?"]
```

### 4. Apply the matched skill

Load the skill using the Skill tool and follow its steps exactly.

### 5. Check for cross-domain coverage

After applying the primary skill, check: does the user's situation span additional
domains with independent high-confidence matches?

- **1 additional domain**: ask once: "This situation also touches [domain] — [skill-name] applies to that aspect. Want me to apply it?"
- **2+ additional domains**: delegate to `plan-solution` — "This situation spans multiple domains. Want me to build a full solution plan?"

Do not chain more than 2 skills without user confirmation. If 3+ skills are needed, use `plan-solution`.

## Rules

- Never ask the user which domain their problem belongs to — that's the skill's job to figure out
- Auto-apply only when 1 skill clearly dominates (score ≥ 0.7, clear gap to second)
- Never auto-apply when uncertain — present choices instead
- Max 5 options in a ranked list; drop lower-scoring matches beyond 5
- Never hallucinate skill names — only reference skills that exist in installed domains
- Announce the matched skill before applying — don't silently load skills
- If a skill is not installed, include the install command in the suggestion

## Examples

**Example 1 — High confidence, single domain**
> "My pull requests keep getting rejected in code review"

Extract: goal=pass-review, symptoms=PR-rejection, domain-cues=code/pull-request/review
Top match: `code-review` (engineering/development) — score 0.82
→ Apply directly: "Situation matches: code-review (engineering/development). Applying now..."

---

**Example 2 — Multiple matches**
> "I always feel exhausted after training"

Extract: goal=recover-better, symptoms=fatigue/exhaustion, domain-cues=training/workout
Top matches (similar scores):
1. `optimize-recovery` (health/fitness) — 0.61
2. `calculate-macros` (health/nutrition) — 0.54
3. `design-training-program` (health/fitness) — 0.49
→ Present ranked list, wait for user selection

---

**Example 3 — No match, clarifying question**
> "I don't know what to do with my life"

Extract: goal=unclear, symptoms=directionlessness, domain-cues=none
All scores < 0.3
→ "No installed skills closely match this. Is this about your career, your health, your finances, or something else?"

## Common Mistakes

**Asking the user to classify their own problem**: "Is this an engineering problem or
a business problem?" — never do this. Route silently, present choices only when
genuinely ambiguous.

**Over-confident routing**: applying a skill when the top score is 0.5 with a close
second-place match. Present choices when it's close.

**Applying too many skills**: chaining 3+ skills without confirmation overwhelms the
user. One at a time, ask before adding the second domain.

**Hallucinating skills**: if no skill exists for the situation, say so. Don't invent
skill names.
