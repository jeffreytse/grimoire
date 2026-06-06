# Skill Examples

Three annotated examples showing what PASS, NEEDS-REVISION, and REJECT look like in practice. Use these to calibrate before running `review-best-practice-skill` on your own file.

---

## Example 1: PASS

A complete skill that passes all review-best-practice-skill criteria. Annotations in `<!-- -->` comments explain what each section does right — remove them in your actual skill files.

```markdown
---
name: calculate-unit-economics          <!-- ✅ verb-first, approved verb (calculate-), specific subject -->
description: Use when evaluating whether a business model is viable before committing
  to go-to-market spend — covers LTV, CAC, payback period, and cohort assumptions.
                                        <!-- ✅ starts "Use when", triggering conditions only, no workflow summary -->
source: Y Combinator startup school curriculum, Andreessen Horowitz growth frameworks,
  SaaStr benchmarks for SaaS unit economics
                                        <!-- ✅ specific institutions, not "widely known" -->
tags: [unit-economics, ltv-cac, startup, business-viability, cost-optimization]
      <!-- ✅ problem: unit-economics, ltv-cac | method: ltv-cac | role: startup | outcome: business-viability, cost-optimization -->
      <!-- ✅ all 4 axes covered -->
---

# Calculate Unit Economics

Compute LTV, CAC, and payback period to verify a business model is viable before
scaling go-to-market spend.
                        <!-- ✅ one-sentence purpose statement immediately after title -->

## Why This Is Best Practice

**Adopted by:** Y Combinator, Andreessen Horowitz, and virtually all professional
investors require unit economics before funding. SaaStr documents that 90%+ of
successful SaaS companies track LTV:CAC ratio as a primary health metric.
                <!-- ✅ specific named institutions, claims majority adoption -->
**Impact:** Companies that validate unit economics before scaling customer acquisition
spend avoid the most common SaaS failure mode: growing a business where each customer
acquired costs more than the lifetime value they generate. a16z portfolio data shows
that companies with LTV:CAC > 3x before Series A have a 4× higher probability of
reaching Series B within 24 months.
                <!-- ✅ measurable outcome (4× probability), named source (a16z) -->
**Why best:** Unit economics calculation — vs. vanity metrics like total revenue or
user growth without margin context — is the only signal that distinguishes a viable
business from one that destroys value at scale. Top-line growth without unit economics
validation is the primary cause of post-Series A companies running out of cash despite
strong revenue growth.
                <!-- ✅ alternative named (vanity metrics), explains why this wins -->

Sources: Y Combinator startup school, a16z growth framework, SaaStr Annual benchmarks

## Steps

### 1. Calculate Customer Acquisition Cost (CAC)

```
CAC = Total Sales & Marketing Spend / New Customers Acquired
```

Use the same time period for both (monthly or quarterly). Include all spend:
salaries, ad spend, tools, agency fees. Exclude product and engineering costs.

### 2. Calculate Customer Lifetime Value (LTV)

```
LTV = Average Revenue Per Account (ARPA) × Gross Margin % × Average Customer Lifetime
```

Average Customer Lifetime = 1 / Monthly Churn Rate

If monthly churn is 2%, average lifetime = 50 months.

### 3. Compute LTV:CAC ratio

```
LTV:CAC = LTV / CAC
```

Interpretation:
- < 1: each customer destroys value — do not scale
- 1–3: marginal — investigate before scaling
- > 3: healthy — standard benchmark for investor-ready SaaS (SaaStr)
- > 5: strong — indicates pricing or retention advantage

### 4. Compute CAC Payback Period

```
CAC Payback Period = CAC / (ARPA × Gross Margin %)
```

Target: under 12 months for B2B SaaS, under 6 months for high-velocity B2C.
A payback period over 24 months is a cash flow risk regardless of LTV:CAC ratio.

### 5. Document cohort assumptions

State the assumptions explicitly:
- Which customer cohort is this based on? (all customers, or a specific acquisition channel?)
- What is the gross margin % used and where does it come from?
- What is the churn rate and how many months of data is it based on?

Undocumented assumptions make unit economics unfalsifiable and useless for decision-making.

## Rules

- Use gross margin, not revenue, for LTV — revenue-based LTV overstates value by the
  cost of serving each customer
- CAC and LTV must use the same customer segment — blended metrics hide channel-level problems
- Payback period under 12 months is required for fundraising at most institutional investors
  (YC, a16z, Sequoia benchmarks)

## Common Mistakes

**Using revenue instead of gross margin in LTV**: inflates LTV by 2–5× and produces
a false healthy ratio.

**Blending acquisition channels**: a Facebook channel with CAC of $200 and an organic
channel with CAC of $20 blend to $110 — a misleading average. Calculate per-channel.

**Ignoring expansion revenue**: for products with upsell/expansion, LTV should include
net revenue retention, not just initial contract value.
```

---

## Example 2: NEEDS-REVISION

A skill with two specific findings that block merge but don't require a full rewrite. Annotations show the exact finding and the fix.

```markdown
---
name: design-sprint-retrospective
description: Use when facilitating an agile sprint retrospective to surface what went
  well, what didn't, and what to improve for the next sprint.
source: Scrum Guide (Schwaber & Sutherland), Atlassian agile playbook
tags: [retrospective, agile, facilitator, team-improvement]
      <!-- ⚠️ NEEDS-REVISION: tags missing outcome axis -->
      <!-- findings: problem=retrospective, method=agile, role=facilitator, outcome=MISSING -->
      <!-- fix: add outcome tag — e.g., team-improvement is borderline; add process-improvement or velocity-improvement -->
---

# Design Sprint Retrospective

Structure a sprint retrospective that produces actionable improvements, not just discussion.

## Why This Is Best Practice

**Adopted by:** Retrospectives are required ceremonies in the Scrum framework, used by
virtually all agile software teams. The Scrum Guide mandates a retrospective after every
sprint. Atlassian's State of Agile survey (2023) found retrospectives used by 87% of
teams practicing agile.
**Impact:** Teams that hold regular retrospectives improve their process over time.
Without retrospectives, recurring problems persist indefinitely.
           <!-- ⚠️ NEEDS-REVISION: Impact is vague — "improve their process" has no number or study -->
           <!-- finding: no measurable outcome, no cited study -->
           <!-- fix: "Atlassian's State of Agile 2023 found teams holding regular retrospectives -->
           <!--       were 2× more likely to report high team satisfaction scores than those that didn't" -->
           <!-- OR: remove the claim and keep only the causal argument: "Without retrospectives, -->
           <!--     recurring impediments have no structured venue for resolution — they persist sprint-to-sprint." -->
**Why best:** A structured retrospective with a defined format — vs. ad-hoc end-of-sprint
conversations — ensures every team member can raise concerns safely and produces a
committed action item, not just discussion.

Sources: Scrum Guide 2020, Atlassian State of Agile 2023
```

**review-best-practice-skill verdict:**

```
## Skill Review: design-sprint-retrospective

### Overall Verdict: NEEDS-REVISION

| Criterion | Result | Finding |
|-----------|--------|---------|
| tags | ⚠️ | Missing outcome axis — add process-improvement or similar |
| Impact | ⚠️ | "improve their process" has no number or study — cite Atlassian data or rewrite as causal argument |
| (all other criteria) | ✅ | — |

### Required fixes before merge
1. Add outcome axis tag (e.g., process-improvement, velocity-improvement, team-satisfaction)
2. Impact line: replace "improve their process over time" with cited Atlassian number
   OR rewrite as: "Without retrospectives, recurring impediments have no structured
   venue for resolution — they persist sprint-to-sprint regardless of severity."
```

---

## Example 3: REJECT

A skill with a blocking REJECT finding in frontmatter. The entire PR is blocked until the required field is fixed.

```markdown
---
name: write-status-update
description: Provides guidance on writing effective status updates for engineering teams.
             <!-- ❌ REJECT: description does not start with "Use when" -->
             <!-- description summarizes content ("provides guidance") instead of triggering conditions -->
             <!-- fix: "Use when writing a weekly or sprint status update for engineering stakeholders — -->
             <!--        covers audience framing, key metrics to include, and escalation signals." -->
source: Google engineering communication practices, Stripe internal writing guide
tags: [communication, status-update, engineering, clarity]
      <!-- ⚠️ additionally: tags missing role and outcome axes (would be NEEDS-REVISION if description were fixed) -->
---
```

**review-best-practice-skill verdict:**

```
## Skill Review: write-status-update

### Overall Verdict: REJECT

| Criterion | Result | Finding |
|-----------|--------|---------|
| description | ❌ | Does not start with "Use when" — describes content instead of triggering conditions |
| tags | ⚠️ | Missing role axis (engineer? manager?) and outcome axis (clarity? alignment?) |
| (all other criteria not evaluated) | — | REJECT on frontmatter blocks further evaluation |

### Required fixes before merge
1. Rewrite description: must start with "Use when" and describe triggering conditions only.
   Example: "Use when writing a weekly or sprint status update for engineering stakeholders —
   covers what to include, how to frame for executives vs. peers, and when to escalate."
2. After fixing description, re-run review-best-practice-skill — tags will also need outcome and role axes.
```

**Why REJECT blocks everything else:** A REJECT finding on any required frontmatter field means the skill cannot be routed correctly by `suggest-best-practice` and fails the basic contract of the skill format. Fix the REJECT finding first, then re-run `review-best-practice-skill` to see if any NEEDS-REVISION findings remain.

---

## Summary

| Verdict | Meaning | Action |
|---------|---------|--------|
| PASS | All criteria ✅ | Open PR |
| NEEDS-REVISION | Specific criteria ⚠️ — fixable without rewrite | Use `revise-best-practice-skill`, fix flagged sections only |
| REJECT | Required field missing or core criterion failed ❌ | Fix blocking issue, re-run review-best-practice-skill |
