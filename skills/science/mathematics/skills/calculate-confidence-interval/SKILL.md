---
name: calculate-confidence-interval
description: Use when reporting statistical estimates that require quantified uncertainty bounds for scientific, clinical, or policy decision-making
source: Neyman "Outline of a Theory of Statistical Estimation" (1937); Fisher "Statistical Methods for Research Workers" (1925); APA statistical reporting guidelines
tags: [statistics, confidence-interval, inference, reporting]
verified: true
---

# Calculate Confidence Interval

Compute and correctly interpret confidence intervals to quantify estimation uncertainty for means, proportions, differences, and regression parameters.

## Why This Is Best Practice

**Adopted by:** APA Publication Manual 7th ed. (requires CIs for all primary statistics), CONSORT clinical trial reporting standards, New England Journal of Medicine statistical reporting guidelines, Cochrane Collaboration meta-analysis methods.

**Impact:** Gardner & Altman (BMJ 1986) demonstrated that CIs convey more information than p-values alone; NEJM mandated CIs in 1988 and subsequent studies showed 40% reduction in dichotomous "significant/not significant" misinterpretation. APA made CIs required in 2010.

**Why best:** A CI provides the precision of an estimate (width) and its plausible range simultaneously, whereas a p-value alone conflates effect size with sample size. CIs are directly usable for clinical or practical significance assessment.

Sources: Neyman Philos Trans R Soc (1937); Gardner & Altman BMJ 292:746 (1986); APA 7th ed. §7.5.

## Steps

1. **Identify the estimator and its distribution** — determine what parameter you are estimating (mean, proportion, difference, correlation, regression coefficient) and the sampling distribution of the estimator (normal, t, chi-square, F, binomial).

2. **Choose the confidence level** — use 95% as default (α=0.05); use 99% for regulatory submissions or when Type I error must be minimized; use 90% only when justified by domain convention. State the level explicitly.

3. **Calculate for a population mean (σ unknown)** — CI = x̄ ± t_(α/2, n−1) × (s/√n), where t_(α/2, n−1) is the critical value from the t-distribution with n−1 degrees of freedom; use z=1.96 only when n ≥ 30 and σ is known.

4. **Calculate for a proportion** — for n×p ≥ 5 and n×(1−p) ≥ 5, use Wilson interval (preferred over Wald): p̂ ± z_(α/2) × √[p̂(1−p̂)/n]. For small samples, use exact Clopper-Pearson method.

5. **Calculate for difference of two means** — CI = (x̄₁ − x̄₂) ± t × SE_diff where SE_diff = √(s₁²/n₁ + s₂²/n₂); use Welch-Satterthwaite degrees of freedom if variances are unequal.

6. **Calculate using software** — in R: `t.test(x, conf.level=0.95)$conf.int`; in Python: `scipy.stats.t.interval(0.95, df=n-1, loc=mean, scale=sem)`; in SPSS: enable "Confidence intervals" in Descriptives/Compare Means.

7. **Bootstrap when assumptions fail** — for non-normal distributions or complex estimators, use percentile bootstrap: resample with replacement B=10,000 times, compute statistic each time, take 2.5th and 97.5th percentiles as CI bounds.

8. **Report completely** — state: estimate [lower CI, upper CI], 95% CI. Example: "Mean difference = 3.2 kg [1.8, 4.6], 95% CI." Always include the confidence level; "95% CI" is not implicit.

## Rules

- Never interpret a 95% CI as "there is a 95% probability the true value lies in this interval" — frequentist CIs are properties of the procedure, not Bayesian credible intervals. The correct interpretation: "95% of intervals constructed this way will contain the true parameter."
- A CI that includes zero (for differences) or 1.0 (for ratios) does not "prove no effect" — it means the data are consistent with zero, not that zero is the true value.
- Report the CI alongside the point estimate; a CI without a point estimate is incomplete.
- Do not adjust CI width post-hoc based on the result — CI width is fixed by the pre-specified α and n.

## Common Mistakes

- **Using z instead of t for small n** — z=1.96 underestimates the critical value for n<30, yielding CIs that are too narrow (undercoverage).
- **Wald interval for proportions near 0 or 1** — the Wald interval performs poorly when p̂ < 0.1 or > 0.9; use Wilson or Clopper-Pearson.
- **Conflating CI with prediction interval** — a CI bounds the population mean; a prediction interval (wider) bounds a single future observation.
- **Multiple comparison inflation** — reporting 10 CIs at 95% confidence gives ~40% familywise error; use Bonferroni correction (95%/k) or simultaneous confidence bands.

## When NOT to Use

- When a Bayesian credible interval better matches the inferential question (use posterior distribution with specified prior)
- When the goal is prediction of individual outcomes (use prediction intervals instead)
- For non-parametric rank-based statistics where standard CI formulas do not apply (use exact or permutation-based CIs)
