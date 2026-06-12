---
name: run-growth-experiment
description: Use when designing and executing a structured growth experiment to improve a key metric
source: Sean Ellis & Morgan Brown "Hacking Growth" (2017), AARRR Pirate Metrics (Dave McClure), ICE Scoring Model
tags: [growth, experimentation, aarrr, ice-scoring, ab-testing, metrics, hypothesis]
verified: true
---

# Run Growth Experiment

Design, run, and learn from a growth experiment using a structured hypothesis-driven process.

## Why This Is Best Practice

**Adopted by:** Facebook Growth team, Airbnb, Uber — all pioneered rapid experimentation culture documented in "Hacking Growth"
**Impact:** Sean Ellis documented that companies running 10+ experiments per week grow 2x faster than those running fewer than 5; Facebook ran 1,000+ experiments per day at peak

**Why best:** Ad hoc "let's try this" changes produce noise, not learning. A structured experiment framework separates signal from noise, quantifies impact, and builds institutional knowledge about what works for a specific product and audience. ICE scoring prevents teams from chasing high-effort, low-impact ideas.

## Steps

1. **Identify the AARRR stage** — Determine which funnel stage has the biggest gap: Acquisition, Activation, Retention, Referral, or Revenue. Focus experiments where the gap is largest.
2. **Generate hypotheses** — Format: "We believe [change] will improve [metric] by [magnitude] because [rationale] for [audience]."
3. **Score with ICE** — Rate each hypothesis on Impact (1-10), Confidence (1-10), Ease (1-10); average the three scores. Run highest-ICE experiments first.
4. **Define success metric and guardrail metric** — Primary metric the experiment moves; guardrail metric that must not regress (e.g., revenue can't drop while testing onboarding change).
5. **Calculate sample size** — Use a power calculator (e.g., Evan Miller's); target 80% statistical power, 95% confidence. Never stop a test early based on early results.
6. **Build and launch** — Implement the minimum viable version of the change. Use feature flags for clean rollback. Document the exact variant differences.
7. **Analyze results** — After reaching sample size, compare primary metric and guardrail metric by variant. Report p-value, confidence interval, and effect size.
8. **Document and socialize** — Write a 1-page experiment brief: hypothesis, method, result, decision (ship / kill / iterate), and learnings. File in a shared experiment log.

## Rules

- Never run experiments without a pre-defined success criterion — post-hoc metric selection is p-hacking.
- Never call a test early — peeking inflates false positive rate; wait for the calculated sample size.
- Run one experiment per funnel stage at a time to avoid interaction effects.
- Document failed experiments with equal rigor — "what didn't work" is as valuable as "what did."
- ICE score is a prioritization tool, not a guarantee; validate assumptions before full commitment.

## Examples

Hypothesis: "Adding a progress bar to the onboarding flow will increase activation rate (first project created) by 15% because users will understand how close they are to value." ICE: Impact 8, Confidence 6, Ease 9 = score 7.7. Run for 2 weeks (calculated 1,200 users per variant). Result: +22% activation, p=0.02. Decision: ship. Learning: visual progress signals significantly reduce onboarding abandonment for this audience.

## Common Mistakes

- Testing without enough traffic — underpowered tests produce inconclusive or false results.
- Running 5 experiments simultaneously on the same funnel — interaction effects make attribution impossible.
- Treating ICE scores as objective truth — they reflect team assumptions; validate high-ICE ideas with qualitative research first.
- Shipping winners without monitoring retention impact — a change that boosts activation but hurts 30-day retention is a loss.

## When NOT to Use

- Do not run an A/B experiment when daily traffic is below the minimum sample size threshold for the expected effect size — underpowered tests will reach statistical significance by chance and produce misleading ship decisions.
- Do not use this structured experiment framework for one-time irreversible changes such as a pricing model overhaul or a full product rebrand, where the cost of running a controlled test exceeds the cost of a phased rollout with monitoring.
- Do not apply ICE-scored experimentation to compliance-driven changes or legal requirements — these must ship regardless of expected impact and do not benefit from prioritization scoring.
