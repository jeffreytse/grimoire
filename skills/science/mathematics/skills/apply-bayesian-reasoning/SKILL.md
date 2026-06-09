---
name: apply-bayesian-reasoning
description: Use when updating beliefs from evidence, estimating probabilities under uncertainty, or making decisions where prior knowledge and new data must be combined
source: 'Bayes "An Essay towards Solving a Problem in the Doctrine of Chances" (1763); Jaynes "Probability Theory: The Logic of Science" (2003); Gelman et al. "Bayesian Data Analysis" (2013)'
tags: [statistics, bayesian, probability, decision-making]
verified: true
---

# Apply Bayesian Reasoning

Update prior beliefs with new evidence using Bayes' theorem to produce calibrated posterior probabilities for inference and decision-making.

## Why This Is Best Practice

**Adopted by:** Clinical diagnostic reasoning (pre-test/post-test probability), spam filtering (naive Bayes), Google search ranking, NASA mission reliability analysis, FDA Bayesian adaptive clinical trial designs (2019 guidance).

**Impact:** Bayesian adaptive clinical trial designs reduce average sample size by 20–30% vs. fixed designs (FDA 2019); Bayesian diagnostic reasoning reduces unnecessary testing by 40% in low-prevalence conditions where base-rate neglect is common (Eddy 1982).

**Why best:** Bayesian reasoning is the mathematically optimal method for incorporating prior information with likelihood evidence; it produces probabilities directly interpretable as degrees of belief, unlike frequentist p-values which do not.

Sources: Bayes (1763) Philos Trans; Jaynes (2003) Cambridge UP; Gelman et al. "BDA3" (2013); Eddy "Probabilistic Reasoning in Clinical Medicine" (1982).

## Steps

1. **Define the hypothesis and its complement** — state H (hypothesis of interest) and H̄ (alternative or complement). Be precise: "The patient has disease D" not "something is wrong."

2. **Establish the prior probability P(H)** — use base rate data: population prevalence for diagnostic questions, historical frequency for reliability questions, expert consensus for novel domains. Document the source. If truly uninformative, use a flat prior P(H) = 0.5 but justify it.

3. **Determine the likelihood ratio (LR)** — for each piece of evidence E, find: LR = P(E|H) / P(E|H̄). For diagnostic tests: LR+ = sensitivity / (1 − specificity); LR− = (1 − sensitivity) / specificity.

4. **Apply Bayes' theorem** — P(H|E) = [P(E|H) × P(H)] / P(E), where P(E) = P(E|H)×P(H) + P(E|H̄)×P(H̄). Equivalently using odds form: posterior odds = LR × prior odds, where odds = p / (1−p).

5. **Use log-odds for sequential updating** — when updating with multiple independent pieces of evidence: log-odds_posterior = log-odds_prior + Σ log(LRᵢ). Convert back: p = 1 / (1 + e^(−log-odds)).

6. **Specify the likelihood function for parameter estimation** — for estimating a parameter θ from data x, write the likelihood L(θ|x) = P(x|θ) and choose a prior π(θ); compute posterior π(θ|x) ∝ L(θ|x) × π(θ).

7. **Compute or approximate the posterior** — use conjugate priors for closed-form solutions (Beta-Binomial, Normal-Normal, Gamma-Poisson); use Markov Chain Monte Carlo (MCMC via Stan or PyMC) for complex models.

8. **Summarize the posterior** — report: posterior mean or median (point estimate), credible interval (e.g., 95% HDI — highest density interval), and posterior probability of hypotheses (e.g., P(θ > 0 | data)).

9. **Perform sensitivity analysis** — vary the prior across a plausible range; if the posterior is robust to prior choice, the data dominate; if sensitive, acknowledge the prior's influence explicitly.

10. **Make decisions using expected utility** — multiply posterior probabilities by outcome utilities and sum; choose the action that maximizes expected utility. Separate inference (posterior) from decision (action).

## Rules

- Prior probabilities must be based on documented evidence or explicitly stated assumptions — "uninformative" priors are not neutral; they encode assumptions about scale and range.
- Do not use Bayes' theorem with conditional independence assumed — verify (or acknowledge) that evidence items are not correlated given H.
- Credible intervals (Bayesian) are not confidence intervals (frequentist) — a 95% credible interval means P(θ ∈ CI | data) = 0.95; this is the statement most people incorrectly attribute to confidence intervals.
- When the prior is diffuse and n is large, Bayesian and frequentist results converge — differences appear most when n is small or the prior is informative.

## Common Mistakes

- **Base rate neglect** — ignoring P(H) and focusing only on the likelihood; even a test with 99% sensitivity gives mostly false positives in a 1% prevalence population.
- **Treating posterior as frequentist p-value** — P(H|data) and 1−p are not the same; "posterior probability of H is 0.95" is not equivalent to "p=0.05."
- **Updating on the same data twice** — using the same data to set the prior and compute the posterior double-counts evidence and overconfidently updates beliefs.
- **Ignoring the prior's influence at small n** — with n=5 observations, the prior dominates; reporting only the posterior without the prior hides this dependence.

## When NOT to Use

- When a frequentist hypothesis test is required by the pre-registered analysis plan or regulatory requirement (run the pre-specified test and report separately)
- For purely exploratory pattern-finding where any model is speculative (use descriptive statistics and visualizations first)
- When the domain has no meaningful prior information and a flat prior would produce results identical to a frequentist MLE (use MLE with uncertainty quantification instead)
