---
name: calculate-measurement-uncertainty
description: Use when reporting physical measurements that require quantified uncertainty estimates for scientific or metrological purposes
source: NIST/SEMATECH "Guide to the Expression of Uncertainty in Measurement" (GUM, ISO/IEC Guide 98-3:2008); Taylor "An Introduction to Error Analysis" (1997); BIPM Joint Committee for Guides in Metrology
tags: [physics, measurement, uncertainty, metrology]
verified: true
---

# Calculate Measurement Uncertainty

Evaluate and report measurement uncertainty using Type A (statistical) and Type B (non-statistical) components following GUM methodology.

## Why This Is Best Practice

**Adopted by:** BIPM (Bureau International des Poids et Mesures), NIST, ISO 17025 accredited laboratories worldwide, European Accreditation EA-4/02, FDA analytical method validation guidelines.

**Impact:** ISO/IEC 17025 requires GUM-compliant uncertainty reporting for accredited test results; without uncertainty quantification, measurement results cannot be compared across laboratories — a fundamental requirement for scientific reproducibility and trade metrological equivalence.

**Why best:** GUM provides a universal, internally consistent framework that covers both random and systematic contributions, replacing the obsolete "probable error" and "maximum error" methods.

Sources: JCGM 100:2008 (GUM); Taylor & Kuyatt NIST TN1297 (1994); Taylor "Introduction to Error Analysis" 2nd ed. (1997).

## Steps

1. **Identify all uncertainty sources** — list every factor that can affect the result: instrument resolution, calibration, repeatability, environmental conditions (temperature, humidity), operator, sampling.

2. **Evaluate Type A uncertainties (statistical)** — take n repeated measurements under the same conditions; calculate standard deviation s; Type A standard uncertainty: u_A = s / √n. Use n ≥ 10 for reliable Type A evaluation.

3. **Evaluate Type B uncertainties (non-statistical)** — for each remaining source, determine the standard uncertainty from: calibration certificate (divide expanded uncertainty by coverage factor k), instrument specification (divide half-range by √3 for rectangular distribution), or expert judgment.

4. **Convert all components to standard uncertainties** — ensure all u_i are expressed as standard uncertainties (k=1, approximately 68% confidence level).

5. **Combine using the law of propagation of uncertainty** — for y = f(x₁, x₂, ..., xₙ):  
   u_c(y) = √[Σᵢ (∂f/∂xᵢ)² · u(xᵢ)²]  
   For uncorrelated inputs with equal sensitivity (∂f/∂xᵢ = 1): u_c = √[Σ u_i²].

6. **Check for dominant contributions** — if one component u_i² > 0.9 × u_c², focus improvement effort on that source; further reducing small contributions yields negligible benefit.

7. **Calculate expanded uncertainty** — multiply combined standard uncertainty by coverage factor k: U = k × u_c. Use k=2 for ~95% confidence (normal distribution); use k from t-table if effective degrees of freedom ν_eff < 30 (Welch-Satterthwaite formula).

8. **Report the result** — state: "y = (result ± U) [unit], where U = k × u_c with k=2 (approximately 95% confidence level)." Example: "L = (23.47 ± 0.05) mm, k=2."

## Rules

- Never report more significant figures in the result than are supported by the uncertainty (round result to match uncertainty's least significant digit).
- Report both u_c (standard uncertainty) and U (expanded uncertainty) with the coverage factor k and confidence level.
- Correlations between input quantities must be accounted for — ignoring positive correlations underestimates combined uncertainty.
- Type B evaluations are not guesses — they must be based on documented sources (calibration certificates, manufacturer specs, published data).

## Common Mistakes

- **Reporting only random error** — omitting systematic contributions (calibration uncertainty, resolution) can understate total uncertainty by a factor of 2–10×.
- **Using maximum error as uncertainty** — "worst-case" addition of errors overestimates uncertainty and is not GUM-compliant.
- **Insufficient replicates for Type A** — n=3 gives a t-multiplier of 4.3 at 95% confidence vs. 2.0 for n→∞; small n dramatically inflates expanded uncertainty.
- **Wrong coverage factor** — using k=2 when effective degrees of freedom are low (e.g., ν_eff=5 needs k=2.57 for 95% coverage).

## When NOT to Use

- For qualitative or categorical measurements where numerical uncertainty is not meaningful
- For Monte Carlo uncertainty propagation (use JCGM 101:2008 Supplement 1 to GUM instead for non-linear models)
- For risk or safety margin calculations in engineering — use conservative design factors per relevant codes (ASME, Eurocode) rather than GUM uncertainty
