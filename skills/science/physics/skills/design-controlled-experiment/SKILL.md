---
name: design-controlled-experiment
description: Use when planning a physics experiment that requires controlling variables to isolate causal relationships between physical quantities
source: Fisher "The Design of Experiments" (1935); Box, Hunter & Hunter "Statistics for Experimenters" (2005); NIST/SEMATECH e-Handbook of Statistical Methods
tags: [physics, experiment-design, controls, scientific-method]
verified: true
---

# Design Controlled Experiment

Plan a physics experiment that isolates a single causal variable while controlling all others, enabling valid causal inference.

## Why This Is Best Practice

**Adopted by:** NIST measurement science protocols, CERN experimental design review process, APS (American Physical Society) experimental reporting standards, NSF research grant requirements.

**Impact:** Factorial design (Box & Hunter) reduces the number of required experimental runs by 50–90% compared to one-factor-at-a-time approaches while simultaneously revealing interactions; NIST calibration studies using designed experiments achieve measurement uncertainty reductions of 30–60%.

**Why best:** Proper control structure is the only mechanism by which an experiment can support causal claims rather than correlational observations; without it, confounders are indistinguishable from true effects.

Sources: Fisher "Design of Experiments" (1935); Box, Hunter & Hunter 2nd ed. (2005) ch. 3–5; NIST/SEMATECH e-Handbook §5.

## Steps

1. **Define the causal question** — state explicitly: "Does [independent variable X] cause a change in [dependent variable Y], holding [list of control variables Z₁, Z₂, ...] constant?"

2. **Map all variables** — categorize every variable as: independent (manipulated), dependent (measured), controlled (held constant), or nuisance (known to vary but not of interest → randomize or block).

3. **Choose the experimental design** — select from: completely randomized design (CRD, for homogeneous material), randomized block design (RBD, for known nuisance variables), factorial design (multiple independent variables and their interactions), or response surface design (optimization).

4. **Define measurement protocol** — specify: instrument, calibration procedure, range, resolution, sampling rate, and number of repeated measurements per condition. Reference `calculate-measurement-uncertainty`.

5. **Randomize run order** — use a random number generator to determine the order of experimental runs; this distributes unknown time-dependent effects (instrument drift, temperature fluctuations) across all conditions.

6. **Implement blocking** — if a nuisance variable cannot be controlled (e.g., different batches of material, different days), group runs into blocks where the nuisance is constant and include block as a factor in analysis.

7. **Replicate** — perform ≥3 independent replications (not repeated measurements of the same run) to estimate run-to-run variability and support statistical inference.

8. **Record all conditions** — log every environmental parameter (temperature, pressure, humidity, operator, instrument serial number) at the time of each run, even those believed to be controlled.

9. **Analyze with ANOVA or regression** — test the effect of the independent variable while accounting for block effects; report effect size (η² or ω²), F-statistic, and p-value; plot residuals to verify model assumptions.

10. **Interpret within scope** — conclusions apply only to the range of conditions tested (no extrapolation beyond the experimental domain without physical model justification).

## Rules

- Control variables must be actively monitored and documented — assuming they are constant without measurement introduces undetectable confounders.
- Never change more than one variable between conditions in a controlled experiment — the entire purpose is isolation.
- Randomization is mandatory; alternating treatment-control-treatment is not randomization.
- Report all failed or excluded runs with reason for exclusion — selective omission is data fabrication.

## Common Mistakes

- **Confounding with time** — running all treatment conditions first, then all controls, means any time-dependent drift (instrument warm-up, sample aging) is confounded with the treatment effect.
- **Technical replicates only** — measuring the same sample three times is not replication; it estimates measurement noise, not run-to-run experimental variability.
- **No baseline or reference condition** — without a reference point (zero-field, room temperature, no treatment), the magnitude of the effect cannot be quantified.
- **Insufficient range of independent variable** — testing X over too narrow a range may miss the functional relationship or underestimate the effect.

## When NOT to Use

- For observational studies where manipulation is impossible (use matched cohort or natural experiment methods)
- For theoretical/computational studies (use sensitivity analysis and verification/validation methodology instead)
- For engineering design optimization over a large parameter space (use Design of Experiments software: JMP, Minitab, or Python pyDOE)
