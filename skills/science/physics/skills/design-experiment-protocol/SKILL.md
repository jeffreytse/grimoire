---
name: design-experiment-protocol
description: Use when designing a physics experiment or laboratory investigation requiring control of variables and uncertainty quantification
source: NIST Guidelines for Evaluating and Expressing the Uncertainty of Measurement (1994), Taylor "An Introduction to Error Analysis" (1997), APS (American Physical Society) experiment design standards
tags: [science, physics, experiment-design, measurement, uncertainty, control-variables, protocol, scientific-method]
verified: true
---

# Design Experiment Protocol

Design a physics experiment with clearly defined variables, control conditions, measurement procedures, and uncertainty analysis.

## Why This Is Best Practice

**Adopted by:** NIST uncertainty guidelines are adopted by all major metrology institutes worldwide; Taylor's error analysis text is used in undergraduate physics labs at MIT, Caltech, and 1,000+ universities
**Impact:** NIST documents that poorly designed experiments with uncharacterized uncertainties account for the majority of irreproducible results in physics literature; rigorous protocol design eliminates this class of error

**Why best:** A physics experiment without a formal protocol produces measurements that cannot be interpreted, compared, or replicated. The protocol specifies what is measured, how it is measured, what is controlled, and how uncertainty is characterized. These decisions made before the experiment prevent the most common failure: collecting data that cannot answer the question the experimenter actually asked.

## Steps

1. **State the hypothesis and measurable prediction** — Write the hypothesis in the form: "If [mechanism] is correct, then [observable variable] will [change in specified direction] when [independent variable] is [changed in specified way]." The prediction must be quantitative and falsifiable.
2. **Identify variables** — Independent variable (what you change, one at a time). Dependent variable (what you measure). Controlled variables (everything else held constant). List all three explicitly.
3. **Specify the measurement method** — For each measured quantity: instrument, range, resolution, calibration procedure, and sampling rate. Never specify "measure the voltage" — specify "measure voltage across R1 using Fluke 87V multimeter in DC mode, ±0.1mV resolution, calibrated against known reference."
4. **Design the control conditions** — Identify and specify the method for controlling each variable that could confound results: temperature (thermostat or water bath), vibration (isolation table), ambient light (dark enclosure), humidity (desiccant or controlled environment).
5. **Determine sample size and repetitions** — For each measurement, determine the number of trials needed to characterize random error. Minimum: 5 trials; preferred: 10-30. Calculate the expected standard error = σ/√n; verify it is smaller than the required precision.
6. **Plan uncertainty analysis** — Identify error sources: Type A (statistical, from repeated measurement) and Type B (systematic, from instrument calibration, environmental conditions, model assumptions). Specify how each will be quantified.
7. **Establish the data recording protocol** — Specify the data table format, units, significant figures, and how anomalous readings are flagged and handled. Define the stopping criterion (number of trials or time window).
8. **Write the safety and equipment protocol** — List hazards (electrical, high-voltage, lasers, pressure, chemical), specify PPE, and define the shutdown procedure. No experiment proceeds without a completed safety review.

## Rules

- Change only one independent variable per experiment — changing two simultaneously makes the cause of any observed effect ambiguous.
- Record raw data, never only derived quantities — raw data allows reanalysis; derived data without raw data is irreproducible.
- Uncertainty must accompany every measurement result — a result without uncertainty (e.g., "g = 9.8 m/s²") is scientifically incomplete; write "g = 9.82 ± 0.05 m/s²."
- Never discard anomalous data without documented justification — unexplained outliers may be the most important data point.

## Examples

Experiment: Measure the effect of temperature on electrical resistance of a metal wire. Hypothesis: Resistance increases linearly with temperature for a pure conductor. Independent variable: temperature (20°C–100°C, 10°C increments). Dependent variable: resistance (R = V/I, measured via 4-wire Kelvin method, Keithley 2400, ±0.01Ω). Controlled: current through wire (1.0 mA constant), wire length (fixed), ambient humidity (<40%). Repetitions: 5 readings per temperature point. Uncertainty: Type A = standard deviation of 5 readings; Type B = 0.05Ω instrument calibration uncertainty. Combined uncertainty: propagated in quadrature.

## Common Mistakes

- Specifying the measurement instrument but not the calibration — systematic errors from uncalibrated instruments produce results that look precise but are inaccurate.
- No controlled variables list — experimenter assumes conditions are constant without verifying; temperature drift or vibration corrupts results invisibly.
- Too few repetitions — single-trial measurements cannot distinguish the true value from measurement noise; minimum 5 trials is required for meaningful Type A uncertainty.
- Confusing precision and accuracy — an instrument with 0.001mm resolution but ±0.5mm calibration accuracy produces precise-looking but inaccurate data; characterize both.

## When NOT to Use

- When conducting a qualitative or exploratory observation aimed at generating hypotheses rather than testing a specific quantitative prediction — rigorous variable control and uncertainty analysis impose premature constraints before the phenomenon is understood well enough to define measurable variables.
- When the experiment is a purely computational or simulation study with no physical measurement — Monte Carlo simulations and numerical models have their own validation methodology and do not require physical instrument calibration or Type A uncertainty from repeated trials.
- When the required precision exceeds what the available instruments and environment can provide — designing a protocol around measurement specifications that cannot be met in the available facility produces a formally correct but practically unexecutable protocol.
