---
name: run-ab-test
description: Use when designing and executing an A/B test on a webpage, email, app, or ad to validate a conversion or engagement hypothesis
source: Optimizely "A/B Testing Guide"; VWO testing methodology; Kohavi et al. "Trustworthy Online Controlled Experiments" (2020, Cambridge)
tags: [ab-testing, experimentation, growth, conversion, statistics]
verified: true
---

# Run A/B Test

Design, execute, and interpret a statistically valid A/B test that produces actionable, trustworthy results for product and marketing decisions.

## Why This Is Best Practice

**Adopted by:** Google, Amazon, Netflix, Booking.com — all run thousands of A/B tests per year as the primary mechanism for product improvement
**Impact:** Netflix reports testing drove $1B+ in annual subscriber retention improvement; Booking.com runs 1,000+ concurrent tests; Kohavi (2020) shows organizations with mature experimentation programs grow 2–4× faster than peers
**Why best:** A/B testing is the gold standard for causal inference on live systems — without it, correlation is mistaken for causation, producing decisions that feel informed but reflect bias.

Sources: Kohavi, Tang & Xu "Trustworthy Online Controlled Experiments" (2020) Cambridge; Optimizely "Statistical Significance Calculator"; VWO "A/B Testing Guide" (2023)

## Steps

1. **Write the hypothesis** — use this format: "Because [observation/insight], we believe [change] will cause [result] for [audience], which we'll measure by [metric]." A bad hypothesis cannot produce learning regardless of the outcome.
2. **Define primary and secondary metrics** — primary metric: the one metric that decides winner/loser; secondary metrics: guardrail metrics that must not deteriorate (e.g., conversion rate is primary; revenue per visitor is guardrail).
3. **Calculate required sample size** — use a power calculator: input baseline conversion rate, minimum detectable effect (MDE, typically 10–20%), significance level (α = 0.05), and power (β = 0.80); this gives required visitors per variant.
4. **Estimate test duration** — duration = required sample size per variant ÷ daily traffic to the test page; minimum 7 days (weekly seasonality), maximum 4–6 weeks (novelty effect fades); if duration >6 weeks, increase MDE or accept lower power.
5. **Implement the test** — use a testing platform (Optimizely, VWO, AB Tasty, Statsig, Eppo, or built-in platform tools); randomize at user level (not session level) to avoid re-assignment; verify traffic split (50/50 for standard A/B).
6. **Run a pre-test check** — run the experiment for 24h with both variants identical (A/A test); if results show significant difference, fix the randomization or data collection bug before proceeding.
7. **Monitor for sample ratio mismatch (SRM)** — check daily that variant traffic volumes match the intended split ±1%; SRM (e.g., variant A gets 60% of traffic instead of 50%) indicates a tracking bug that invalidates results.
8. **Do not peek** — avoid checking results before the planned end date; peeking and stopping early inflates false positive rates from 5% to 30%+ (Kohavi 2020); use sequential testing if early stopping is required.
9. **Read results correctly** — declare a winner only if: p < 0.05 (or 95% CI excludes zero) AND primary metric moved in the expected direction AND guardrail metrics are not significantly harmed.
10. **Ship, learn, and document** — if significant: ship the winner; if not significant: document the null result as a learning (null results are equally informative); add to the experiment log; plan the next iteration.

## Rules

- One primary metric per test — multiple primary metrics require multiple hypothesis tests and Bonferroni correction; keep it simple.
- Never end a test early because the winner "looks obvious" — this is the most common source of invalid A/B test conclusions.
- Sample ratio mismatch must be investigated before reading results — SRM makes all statistics meaningless; it must be resolved or the test must be restarted.
- Minimum 7 days runtime for any web/app test — day-of-week effects can make early results reverse completely by day 7.
- Null results must be documented — a test that shows no effect tells you what not to build; it is equally valuable to the organization as a positive result.

## Common Mistakes

- **No pre-calculated sample size** — ending the test when "it feels like enough data" produces wildly wrong conclusions.
- **Changing the test mid-run** — modifying the variant, traffic percentage, or metric after the test starts invalidates all prior data.
- **Peeking and stopping early** — the #1 source of false positives in A/B testing; p-values are not valid at arbitrary stopping points.
- **Testing too many variants** — 4+ variants require much larger sample sizes and longer run times; 2-variant tests are faster and more reliable.
- **Ignoring guardrail metrics** — a variant that improves CTR while halving revenue per session is not a winner.

## When NOT to Use

- Traffic too low to reach statistical significance within a reasonable timeframe
- Changes that cannot be randomized at the user level (e.g., pricing changes visible to groups)
- Qualitative questions (use user interviews, surveys, or usability tests instead)
