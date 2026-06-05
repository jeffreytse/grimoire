---
name: run-usability-testing
description: Use when planning or conducting usability tests to identify UI friction, task completion failures, or navigation problems with real users
source: Nielsen "Usability Engineering" (1993); Nielsen Norman Group usability testing guidelines; ISO 9241-11:2018 usability definition and measurement
tags: [usability, research, ux, testing, user-research]
verified: true
---

# Run Usability Testing

Systematically observe real users attempting real tasks to surface usability problems before or after launch.

## Why This Is Best Practice

**Adopted by:** Google, Apple, Microsoft, Amazon — all run weekly usability sessions; mandated by US federal Section 508 compliance process
**Impact:** Nielsen (1993) showed 5 users find 85% of usability problems; fixing usability issues costs 10–100× less before development than after launch (Pressman software engineering research)
**Why best:** ISO 9241-11 defines usability as effectiveness + efficiency + satisfaction — testing operationalizes all three with actual task data rather than assumptions.

Sources: Nielsen "Usability Engineering" Ch. 6; Nielsen Norman Group "How Many Test Users"; ISO 9241-11:2018

## Steps

1. **Define research questions** — write 3–5 specific questions you need answered (e.g., "Can users find the export function without assistance?").
2. **Select participant profile** — define 2–3 user segments; recruit 5 participants per segment matching real user demographics and technical proficiency.
3. **Write task scenarios** — create realistic, goal-based tasks using natural language; avoid UI vocabulary (say "pay your bill" not "click the payment button").
4. **Choose test format** — moderated in-person: richest data; moderated remote (Zoom + screen share): broad reach; unmoderated (Maze, UserTesting): scale with less depth.
5. **Prepare test materials** — build or stage prototype/product at correct fidelity; create screener, consent form, NDA, and observer guide.
6. **Conduct pilot test** — run one session with a colleague to validate task clarity, timing (~45–60 min total), and recording setup.
7. **Run sessions** — facilitate without leading: use "think aloud" protocol; ask "what are you thinking?" not "is this confusing?"; record screen, audio, and (optionally) face.
8. **Capture observations** — log quotes, behaviors, hesitations, and errors per task in a structured observation sheet (participant × task grid).
9. **Analyze and prioritize** — identify patterns across participants; rate severity using Nielsen's 0–4 scale (0 = not a problem, 4 = usability catastrophe); prioritize by frequency × severity.
10. **Report and action** — present findings as problem statements with evidence clips; pair each finding with a recommended design direction; schedule follow-up test to validate fixes.

## Rules

- Never correct users during tasks — intervene only if they are completely stuck for >3 minutes or distressed.
- Tasks must be completable on the prototype; do not test features that do not exist.
- Record observer impressions silently during sessions; debrief after, not during.
- Report findings as behavioral evidence, not opinions: "3 of 5 participants could not find X" not "the navigation is confusing."
- Fix the top 3 severity-4 issues before running another round; iteration is the mechanism of improvement.

## Common Mistakes

- **Testing with colleagues or designers** — familiarity bias invalidates results; recruit actual target users.
- **Leading questions** — "Was that confusing?" biases toward yes; ask "What were you trying to do there?"
- **Recruiting only 1–2 participants** — too few to identify patterns; minimum 5 per segment.
- **Testing too late** — usability testing a shipped product with no change budget wastes effort; test prototypes.
- **Reporting all findings equally** — 20+ minor issues overwhelm stakeholders; severity ratings focus the fix list.

## When NOT to Use

- Validating visual aesthetics alone — use preference testing or surveys instead
- Measuring adoption at scale — use analytics and cohort analysis
- Comparing two designs quantitatively — use A/B testing with sufficient statistical power
