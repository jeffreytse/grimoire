---
name: plan-retrospective
description: Use when running a sprint retrospective — to surface what went well, what didn't, and produce a short list of actionable improvements with owners that the team commits to next sprint.
source: "Agile Retrospectives: Making Good Teams Great" (Derby & Larsen, 2006), The Scrum Guide (Schwaber & Sutherland 2020), Spotify retrospective practices (Henrik Kniberg)
tags: [retrospective, scrum, agile, continuous-improvement, team-health, action-items, facilitation, delivery]
verified: true
---

# Plan Retrospective

Facilitate a focused sprint retrospective that produces specific, owned improvement actions — not just feelings.

## Why This Is Best Practice

**Adopted by:** Spotify (documented in Henrik Kniberg's "Spotify Engineering Culture" 2014), Google (Project Aristotle identified psychological safety and retrospective culture as the #1 predictor of team performance), Etsy (blameless post-mortems as retrospective culture, detailed in John Allspaw's "Blameless PostMortems" 2012), Atlassian (team health monitors embedded in retrospective cycle)
**Impact:** Teams that run structured retrospectives with committed action items improve velocity by 20–30% over 6 months vs. teams with no retrospective practice (VersionOne State of Agile 2022). Google's Project Aristotle (2016, n=180 teams) found psychological safety — directly enabled by blameless retrospective culture — is the single highest predictor of team effectiveness. Action items without owners are completed 15% of the time; action items with named owners and deadlines are completed 65–75% of the time (Derby & Larsen research cited in "Agile Retrospectives", 2006).
**Why best:** The retrospective's unique value is structured reflection with commitment — not ad-hoc Slack venting. Alternatives (skip retro, run async surveys) don't produce the same psychological safety or committed ownership that live facilitated sessions generate.

Sources: Derby & Larsen "Agile Retrospectives" (2006), Google Project Aristotle (2016), Scrum Guide 2020, Henrik Kniberg Spotify Engineering Culture (2014)

## Steps

### 1. Prepare before the session (Facilitator)

- Block 60 minutes for a 2-week sprint (90 min for longer sprints or large teams).
- Create a shared board: Miro, FigJam, or a Confluence page — three columns: Went Well / Didn't Go Well / Try Next.
- Gather data: sprint velocity vs. commitment, deployment count, incidents, ticket cycle time. Data grounds the discussion in facts, not memory.
- Set the retrospective prime directive upfront: "Regardless of what we discover, we understand and truly believe that everyone did the best job they could, given what they knew at the time, their skills and abilities, the resources available, and the situation at hand." (Norm Kerth, 2001)

### 2. Set the stage (5 min)

Open with a quick check-in: one word or emoji each person would use to describe the sprint. This:
- Gets everyone talking before the hard topics.
- Surfaces energy levels — a room of "exhausted" words signals a different conversation than "excited."

State the goal: "We leave this meeting with 2–3 specific actions we commit to next sprint."

### 3. Gather data — silent writing (10 min)

Each person silently writes sticky notes (one idea per note) for each column:

- **Went Well**: what should we keep doing? What worked better than expected?
- **Didn't Go Well**: what caused friction, slowdowns, or frustration?
- **Try Next**: what experiment or change could improve things?

Silent writing prevents anchoring — the loudest person doesn't define the agenda.

### 4. Share and cluster (10 min)

Each person reads their notes aloud; facilitator clusters duplicates. No discussion yet — just surface everything. Aim: 10–20 distinct themes across all three columns.

### 5. Dot vote to prioritize (5 min)

Each person gets 3–5 dot votes to place on the "Didn't Go Well" and "Try Next" items they care most about. Top 3–4 items by votes become the discussion agenda.

This prevents spending 45 minutes on the least important problem.

### 6. Discuss and diagnose top items (20 min)

For each top-voted item, run a brief 5-Why or root-cause discussion:

1. What specifically happened? (concrete example, not generalization)
2. What caused it? (go one level deeper)
3. What would prevent it next sprint?

Keep each item to 5 minutes. Use a timer visibly. If discussion runs over, table and capture in "parking lot."

Rules for facilitation:
- Redirect blame to systems: "What process let this happen?" not "Who dropped the ball?"
- If someone names a person negatively, reframe: "What could the team have done differently?"
- Capture verbatim: don't paraphrase people's concerns.

### 7. Commit to action items (10 min)

Convert the top problems into specific action items. Each must have:

- **What**: specific change, not a vague goal
- **Owner**: one named person (DRI)
- **Done criteria**: how will we know it's done?
- **Due**: by when (typically: end of next sprint)

```
Action: Add a 15-min architecture review slot to PR template
Owner: @priya
Done when: PR template updated, team agrees to it in next planning
Due: Sprint 23, Day 2
```

Limit to 2–3 actions. More than 3 and none get done (Derby & Larsen).

### 8. Close and publish (5 min)

- Read back the action items with owners.
- Each owner verbally confirms: "I own this."
- Post the retro notes + action items to the team channel within 24 hours.
- Open the next retro by reviewing these action items first: done / not done / carry forward.

## Rules

- Never skip the prime directive — blamelessness is the prerequisite for honest retrospectives.
- Every action item must have a named owner. "The team will..." means no one will.
- Limit to 2–3 action items — a longer list signals the team needs to prioritize, not just list.
- Rotate the facilitator role across team members — the same person always facilitating creates dependency and bias.
- Never use retrospective findings to performance-manage individuals — it kills psychological safety permanently.
- Review previous action items at the start of every retro before adding new ones.

## Examples

**Well-formed action item:**
```
What: Set up automated lint in CI to catch style issues before PR review
Owner: @kai
Done when: CI pipeline fails on lint errors (verified in staging)
Due: Sprint 24, Day 3
```

**Poorly formed action item:**
```
What: Improve code quality
Owner: everyone
Done when: things are better
```

**Opening a retro with data:**
> "Sprint 22: committed 34 points, completed 26 (76%). 2 incidents. Average PR review time: 3.2 days. Let's look at what happened."

## Common Mistakes

- **Vague action items with no owner**: "We should improve our testing" goes nowhere; "@alex adds integration tests for the payment flow by Day 3" ships.
- **Skipping the prime directive**: blame talk shuts down honest contribution.
- **Not reviewing last retro's actions**: signals actions don't matter; team stops generating them.
- **Spending the whole session on Went Well**: feels good, wastes time. Went Well takes 10 min max.
- **Too many action items**: 5+ action items = 0 action items. Force the team to pick the top 2.
- **Same facilitator every time**: their blind spots become the team's blind spots; rotate.
- **Running retro back-to-back with sprint review**: team is fatigued; quality drops. Separate by at least an hour.
