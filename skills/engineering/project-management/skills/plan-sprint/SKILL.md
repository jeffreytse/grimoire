---
name: plan-sprint
description: Use when starting a new sprint — to review the backlog, set the sprint goal, select stories, estimate capacity, and assign work before the team begins execution.
source: The Scrum Guide (Schwaber & Sutherland 2020), Atlassian Agile Coach documentation, Shape Up (Basecamp, 2019)
tags: [sprint-planning, scrum, agile, capacity-planning, backlog, estimation, team-coordination, delivery]
verified: true
---

# Plan Sprint

Run a structured sprint planning session that produces a clear sprint goal, a committed backlog slice, and capacity-matched assignments.

## Why This Is Best Practice

**Adopted by:** Spotify (squad sprint cadence), Atlassian (internal product teams), Amazon (two-pizza team sprint rituals), Linear (weekly cycle planning used internally and recommended to customers)
**Impact:** Teams that set explicit sprint goals complete 89% of committed work vs. 67% for teams without goals (VersionOne "State of Agile" 2022, n=1,000+ teams). Capacity planning reduces sprint spillover by 30–40% compared to ad-hoc commitment (McKinsey Agile survey, 2021). Proper sprint planning reduces mid-sprint scope changes by 50% (Scrum Alliance survey, 2019).
**Why best:** Sprint planning forces explicit prioritization and capacity math before work starts — the only moment where changing the plan is cheap. Ad-hoc ticket assignment produces invisible overcommitment and reactive scope cuts mid-sprint. Shape Up's appetite model (fixed time, variable scope) is an alternative for product teams, but Scrum planning with explicit goals achieves the same constraint clarity with more granular commitments.

Sources: The Scrum Guide 2020, VersionOne State of Agile 2022, Atlassian Agile Coach, Shape Up (Basecamp 2019)

## Steps

### 1. Prepare the backlog before the meeting (Product Owner)

Do this before the planning session — don't waste team time on ungroomed tickets.

- Top 1.5× sprint capacity worth of stories must be estimated and acceptance-criteria'd.
- Each story has: title, description, acceptance criteria, story point estimate, and no unresolved blockers.
- Stories are ordered by business priority (highest value first).

If the backlog top is not groomed, delay planning until it is. Starting planning with ungroomed stories wastes the entire team's time.

### 2. Open: state the sprint goal (10 min)

The Product Owner proposes a one-sentence sprint goal:

> "By end of sprint, users can [do X], which achieves [business outcome Y]."

The team challenges: is this achievable? Does it align with roadmap? Refine until everyone can repeat it from memory. The sprint goal is the north star — every story selected must contribute to it or be explicitly marked as maintenance/debt work.

### 3. Calculate team capacity (10 min)

```
Capacity (story points) = (team_members × sprint_days × focus_factor) / avg_story_duration
```

Simpler: use historical velocity (last 3 sprints average) as the capacity ceiling.

Adjust for:
- PTO / holidays: subtract proportionally (1 person out 2 days = -10% capacity for that person)
- Meetings / ceremonies: Scrum recommends max 10% of sprint time in ceremonies
- On-call rotation: subtract 20% for whoever is on-call during the sprint

Write the number down: `Sprint capacity: N story points`.

### 4. Select stories from the backlog (30–45 min)

Pull stories from the top of the backlog until you reach sprint capacity. For each story:

1. Read the acceptance criteria aloud.
2. Confirm the estimate is still accurate (re-estimate if scope changed since grooming).
3. Check for dependencies: does this story block or require another story in this sprint?
4. Ask: "Is anything unclear about how to implement this?" If yes, spike or defer.

Stop when total selected points = sprint capacity. Resist adding "one more small thing."

### 5. Break stories into tasks (optional, 15–20 min)

For complex stories (≥5 points), decompose into tasks:
- Each task: 1–2 days max, owned by one person, has a clear done state.
- Surface hidden work: migrations, feature flags, test coverage, documentation.

This step reveals hidden complexity — if decomposing reveals a 3-point story is actually 8 points of tasks, re-estimate and possibly defer.

### 6. Assign ownership

Each story gets one owner (DRI — Directly Responsible Individual). The DRI is accountable for delivery; others can help.

Rules:
- Don't assign one person >80% of team capacity — bus factor risk.
- Pair junior engineers with a senior mentor on complex stories.
- No story unowned at the end of planning.

### 7. Confirm and publish

Read back:
1. Sprint goal (one sentence)
2. Stories committed (list with points and owners)
3. Total points vs. capacity

Post to the team channel and update the project board. Any story not in the sprint goes back to the backlog — not "in progress."

## Rules

- Never commit more than 90% of calculated capacity — buffer for unplanned work.
- Sprint goal must be set before stories are selected, not after.
- Every selected story must have an owner before the meeting ends.
- Ungroomed stories (missing estimates or acceptance criteria) cannot be committed.
- Stories that don't contribute to the sprint goal must be explicitly labeled (maintenance, debt, incident follow-up).
- If planning takes >2 hours for a 2-week sprint, the backlog is not groomed — stop and groom first.

## Examples

**Well-formed sprint goal:**
> "Users can complete checkout with saved payment methods, reducing checkout abandonment."

**Poorly formed sprint goal:**
> "Work on payment stuff and fix some bugs."

**Capacity calculation:**
```
5 engineers × 10 days × 0.7 focus factor = 35 engineer-days
Team velocity: 40 points/sprint (last 3 sprint average)
Adjust: 1 engineer PTO 2 days → -4 points
Effective capacity: 36 points → commit up to 32 (90% ceiling)
```

## Common Mistakes

- **Committing to 100% capacity**: leaves no room for production incidents, review cycles, or complexity surprises. Use 80–90%.
- **Starting without a sprint goal**: teams drift to "finish whatever's at the top" with no coherent outcome.
- **Skipping acceptance criteria review**: "Done" becomes ambiguous mid-sprint; QA fails at the last moment.
- **Over-assigning one engineer**: single point of failure; their blocker kills the whole sprint.
- **Carrying over last sprint's unfinished stories without re-evaluating**: stale context means stale estimates.
- **Planning in a single marathon session**: planning >2 hours signals backlog is ungroomed — fix the input, not the meeting.
