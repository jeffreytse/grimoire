---
name: design-experiment
description: Use when designing an A/B test or online controlled experiment to evaluate a product change
source: Kohavi, Tang & Xu "Trustworthy Online Controlled Experiments" (Cambridge University Press 2020); Evan Miller "How Not To Run an A/B Test" (evanmiller.org)
tags: [experimentation, ab-testing, statistics, product, data-science, hypothesis-testing]
verified: true
---

# Design Experiment

Design statistically sound A/B experiments that produce trustworthy causal evidence for product decisions.

## Why This Is Best Practice

**Adopted by:** Microsoft (ExP platform, Kohavi's team), Google, Netflix, Airbnb — all maintain internal experimentation platforms
**Impact:** Kohavi et al. report that only 1/3 of A/B experiments at Microsoft produce a positive result; without rigorous experiment design, teams ship features that feel successful but have no causal impact or actively harm metrics.

**Why best:** A/B testing is the only method that provides causal evidence in product development. Without it, correlation-based decisions (feature launched, DAU went up — success!) fail to account for confounders. Trustworthy experiments require pre-registration of hypotheses, power analysis, and fixed analysis windows.

## Steps

1. **Define the hypothesis** — Write a falsifiable hypothesis: "Changing X will increase metric Y by Z% because [mechanism]." Define primary metric (OEC — Overall Evaluation Criterion) before running.
2. **Select guardrail metrics** — Identify metrics that must not regress (e.g., latency, error rate, revenue) even if the primary metric improves.
3. **Run a power analysis** — Calculate required sample size: specify minimum detectable effect (MDE), significance level α=0.05, power 1-β=0.80. Use a power calculator (e.g., `pwr` in R, `statsmodels` in Python). Do not start the test without sufficient sample size.
4. **Randomize correctly** — Randomize at the user level (not session or request) to avoid carryover effects. Use a consistent hash of user ID + experiment ID for deterministic assignment.
5. **Run the experiment for the pre-determined duration** — Never stop early based on results (peeking). Run for at least one full business cycle (typically 1-2 weeks) to account for day-of-week effects.
6. **Analyze with the pre-registered method** — Apply t-test or Mann-Whitney U; report p-value, confidence interval, and effect size. Check for SRM (Sample Ratio Mismatch) first — if control/treatment ratio deviates >5% from target, the experiment is invalid.
7. **Document and share results** — Record hypothesis, methodology, results (including null results), and decision. Null results are valuable — they prevent re-running the same failed experiment.

## Rules

- Never change the experiment while it is running — adds confounders that invalidate results.
- Ship only if the primary metric improves AND no guardrail metrics regress significantly.
- Publish null results — a well-run experiment that shows no effect is a valid and valuable outcome.
- Avoid running multiple overlapping experiments on the same users without interaction analysis.

## Examples

Hypothesis: Adding social proof ("1,200 users bought this today") to the product page will increase add-to-cart rate by 5%.
OEC: Add-to-cart rate. Guardrail: Page load time P95, return rate.
MDE: 5%, α=0.05, power=0.80 → required n=8,400 per variant → 2 weeks at current traffic.
Result: +3.2% (95% CI: [1.1%, 5.3%], p=0.003). Ship.

## Common Mistakes

- **Peeking and stopping early** — inflates false positive rate; at α=0.05, peeking 5× gives a true false-positive rate of ~19%.
- **No power analysis** — underpowered experiments produce inconclusive results and waste traffic.
- **Using page views as the unit of randomization** — same user sees both variants; violates stable unit treatment value assumption (SUTVA).

## When NOT to Use

- When daily active users are too low to reach the required sample size within a reasonable time window, running the experiment will produce underpowered, inconclusive results regardless of how long it runs.
- When the change being evaluated is a critical security patch, compliance fix, or regression repair, ethical and legal obligations require shipping immediately rather than exposing a control group to a known defect.
- When the product surface has strong network effects between users (e.g., social feeds, multiplayer features), user-level randomization creates interference between variants that invalidates the causal assumption underpinning the experiment.
