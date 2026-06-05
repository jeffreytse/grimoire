---
name: calculate-statistical-power
description: Use when determining required sample size for a study, or evaluating whether a completed study had adequate power to detect its effect
source: Cohen "Statistical Power Analysis for the Behavioral Sciences" (1988); Button et al. "Power Failure" Nature Reviews Neuroscience (2013); G*Power methodology (Faul et al. 2007)
tags: [statistics, sample-size, power-analysis, experimental-design]
verified: true
---

# Calculate Statistical Power

Determine the minimum sample size needed to detect a biologically meaningful effect with specified confidence, preventing both underpowered and wasteful studies.

## Why This Is Best Practice

**Adopted by:** NIH grant review criteria (power statement required since 2014), FDA guidelines for clinical trials, Nature/Science submission checklists, ARRIVE 2.0 animal research guidelines.

**Impact:** Button et al. (2013) found median power of 21% in neuroscience studies — meaning 79% of negative results were false negatives. Studies with ≥80% power reduce false-negative rate to <20% and improve replication rates by ~2×.

**Why best:** Power analysis forces researchers to specify effect size before data collection, preventing HARKing; it quantifies the tradeoff between Type I error (α), Type II error (β), effect size, and n.

Sources: Cohen (1988) chapters 2–8; Button et al. Nature Rev Neurosci 14:365–376 (2013); Faul et al. Behav Res Methods 39:175–191 (2007).

## Steps

1. **Set α (significance level)** — use α=0.05 as default; use α=0.01 for exploratory studies with many comparisons; adjust for multiple comparisons (Bonferroni: α/k).

2. **Set desired power (1−β)** — use 0.80 as minimum; use 0.90 for high-stakes experiments (clinical, irreversible interventions).

3. **Specify the statistical test** — identify the test you will use: t-test, ANOVA, chi-square, correlation, regression, survival analysis. Power formulas differ by test.

4. **Estimate the effect size** — use Cohen's d for t-tests, f for ANOVA, r for correlations, w for chi-square. Sources in order of preference: (a) pilot data, (b) meta-analytic estimate, (c) smallest effect of biological importance, (d) Cohen's conventional values (small d=0.2, medium d=0.5, large d=0.8).

5. **Calculate n using G*Power or formula** — for two-sample t-test: n = 2(z_α/2 + z_β)² / d²; use G*Power 3.1 software for complex designs or run in R: `pwr::pwr.t.test(d=0.5, sig.level=0.05, power=0.80)`.

6. **Account for attrition** — inflate n by expected dropout rate: n_adjusted = n / (1 − dropout_rate). Use 10–20% for animal studies, 20–30% for human clinical trials.

7. **Report power justification** — write: "We require n=X per group to detect d=Y with 80% power at α=0.05 (two-tailed), based on [source of effect size estimate]."

8. **For completed studies** — calculate observed power only to contextualize a non-significant result; do not use observed power to retroactively justify design.

## Rules

- Never use post-hoc observed power to interpret a non-significant p-value — it is mathematically circular (Hoenig & Heisey 2001).
- Use pilot study effect sizes only as a starting point; pilot n is too small for reliable effect size estimation — prefer literature meta-analyses.
- Report all four parameters (α, power, effect size, n) in the methods section — partial reporting is insufficient for replication.
- When effect size is unknown, perform sensitivity analysis across a plausible range of effect sizes.

## Common Mistakes

- **Using "medium" effect size without justification** — defaulting to d=0.5 without biological rationale inflates or deflates n arbitrarily.
- **Ignoring multiple comparisons** — running 10 tests at α=0.05 without correction yields ~40% familywise error rate.
- **No attrition adjustment** — reaching target enrollment without buffer often results in underpowered final analyses.
- **Confusing one-tailed and two-tailed** — one-tailed tests require directional hypothesis stated a priori; using them post-hoc is p-hacking.

## When NOT to Use

- For purely exploratory/descriptive studies where no hypothesis test is planned (report effect sizes and CIs instead)
- For N-of-1 or case study designs where statistical inference is not the goal
- When performing Bayesian analysis (use Bayes Factor or decision-theoretic sample size planning instead)
