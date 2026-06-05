---
name: design-kpi-dashboard
description: Use when building or redesigning a performance dashboard to track organizational, team, or product KPIs and support data-driven decisions
source: Kaplan & Norton "The Balanced Scorecard" HBR (1992); Eckerson "Performance Dashboards" (2010); Doerr "Measure What Matters" (2018)
tags: [operations, kpi, dashboard, metrics, data-driven, balanced-scorecard]
verified: true
---

# Design KPI Dashboard

Design a performance dashboard that surfaces the right metrics to the right audience for faster, better-informed decisions.

## Why This Is Best Practice

**Adopted by:** Fortune 500 operations teams, Tableau customers, Google, and organizations using OKR or Balanced Scorecard frameworks
**Impact:** Kaplan & Norton demonstrated that organizations using Balanced Scorecard dashboards outperformed peers by 35% on shareholder return. Eckerson research found companies with effective performance dashboards make decisions 5x faster than those relying on static reports. Doerr documented that OKR-aligned dashboards at Google increased goal attainment by 30%.
**Why best:** Dashboards fail not from poor visualization but from poor metric selection. Most organizations track what is easy to measure (outputs) rather than what drives outcomes (leading indicators). A structured design process ensures the dashboard measures what matters.

Sources: Kaplan & Norton "The Balanced Scorecard" HBR (1992); Eckerson "Performance Dashboards: Measuring, Monitoring, and Managing Your Business" (2010); Doerr "Measure What Matters" (2018)

## Steps

1. **Define the audience and decision context** — different audiences need different dashboards. An executive needs 5–8 high-level strategic metrics; an operations manager needs 15–20 operational indicators. Define: who sees this, what decisions they make with it, and at what frequency.

2. **Align metrics to strategy** — use the Balanced Scorecard four perspectives as a framework: (1) Financial (revenue, margin, cash), (2) Customer (NPS, churn, satisfaction), (3) Internal Process (cycle time, quality, throughput), (4) Learning & Growth (employee engagement, capability building). Ensure each perspective is represented.

3. **Separate leading from lagging indicators** — lagging metrics (revenue, quarterly NPS) show what happened. Leading metrics (pipeline coverage, product usage by new users) predict what will happen. A good dashboard includes both, with more leading indicators for operational audiences.

4. **Select 5–10 metrics maximum per dashboard** — every metric added reduces attention paid to every other metric. Apply the "fewer, better" principle. If a metric doesn't drive a decision or action, remove it.

5. **Define each metric precisely** — for every KPI, document: name, definition, calculation formula, data source, refresh frequency, owner, and target. Ambiguous definitions produce debates about the number, not the decision.

6. **Set targets and thresholds** — for each metric, define: green (on track), yellow (watch), and red (requires action) thresholds. Targets without thresholds require judgment on every review; thresholds automate the interpretation.

7. **Design the visual hierarchy** — place strategic summary metrics at the top (highest visibility), operational drill-downs below. Use color coding consistently. Red/yellow/green is universal; don't invent custom color conventions.

8. **Connect metrics to owners** — every metric must have a named owner accountable for the number and for driving corrective action when the metric is red. Ownerless metrics are watched by no one.

9. **Automate data pipelines** — manually updated dashboards become stale within weeks. Connect to live data sources (database, API, BI tool). Refresh frequency should match decision frequency: daily for operations, weekly for management, monthly for executives.

10. **Review and prune quarterly** — remove metrics that generate no decisions, add metrics aligned to current strategic priorities. A dashboard that never changes stops being read.

## Rules

- Never include a metric you cannot act on — if no action follows from the red state, remove the metric.
- All metrics must have a named owner — shared ownership is no ownership.
- Targets must be set before the period begins, not after results are known.
- Dashboard must refresh on a cadence that matches the decision cycle it supports.
- Limit executive dashboards to 8 metrics maximum — beyond 8, executives default to the 2–3 they already care about.

## Common Mistakes

- **Measuring what's easy, not what matters** — page views, hours worked, and lines of code are output metrics that rarely drive strategic outcomes.
- **Too many metrics** — dashboards with 30+ KPIs are ignored; leaders focus on 2–3 familiar numbers and ignore the rest.
- **No baseline or target** — a metric without context ("NPS: 42") is meaningless; show trend, benchmark, and target alongside.
- **Static dashboards** — dashboards updated manually monthly are already stale by the time they're shared; automate or don't build.

## When NOT to Use

- One-time analysis or research questions (use ad-hoc reporting or a data query instead)
- Exploratory data analysis before metrics are defined (use EDA first, then build the dashboard)
- Teams without data infrastructure to support automated refresh (build the data foundation first)
