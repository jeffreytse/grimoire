---
name: calculate-reaction-yield
description: Use when determining the theoretical, actual, or percent yield of a chemical reaction from experimental data
source: IUPAC nomenclature and yield calculation standards; March "Advanced Organic Chemistry" (2007); ACS Style Guide laboratory reporting
tags: [chemistry, stoichiometry, yield, laboratory]
verified: true
---

# Calculate Reaction Yield

Compute theoretical yield, actual yield, and percent yield for a chemical reaction using stoichiometric principles.

## Why This Is Best Practice

**Adopted by:** IUPAC standardized reporting, ACS journal submission requirements, FDA manufacturing validation (process yield ≥98% threshold), EPA Green Chemistry metrics.

**Impact:** Accurate yield calculation enables reaction optimization — a 10% yield improvement in a 10-step synthesis compounds to >2× overall yield gain. Green Chemistry metrics (E-factor, PMI) depend on accurate yield data.

**Why best:** Stoichiometric yield calculation from balanced equations is the only way to objectively compare reactions across different scales, solvents, and conditions.

Sources: IUPAC Green Book (2007) §2.10; March & Smith "March's Advanced Organic Chemistry" 7th ed. (2013); ACS Style Guide 3rd ed.

## Steps

1. **Write and balance the equation** — ensure atom balance and charge balance; identify the limiting reagent and all stoichiometric coefficients.

2. **Identify the limiting reagent** — calculate moles for each reactant: n = mass(g) / MW(g/mol). The reagent with the smallest n × (1/stoichiometric coefficient) is limiting.

3. **Calculate theoretical yield** — multiply moles of limiting reagent by the stoichiometric ratio to the product, then convert to grams: theoretical yield (g) = n_limiting × (MW_product / stoichiometric ratio).

4. **Measure actual yield** — weigh the isolated, purified product after workup and drying to constant mass. Record on an analytical balance (±0.1 mg precision).

5. **Calculate percent yield** — % yield = (actual yield / theoretical yield) × 100. Values >100% indicate impurity, incomplete drying, or weighing error.

6. **Check for side products** — if percent yield is low (<50%), identify likely side reactions and byproducts; use TLC, NMR, or HPLC to assess purity of the product.

7. **Calculate atom economy (optional)** — atom economy = (MW_desired product / sum of MW all products) × 100; report alongside % yield for green chemistry assessment.

8. **Report with units and purity** — always state: actual yield (g or mg), % yield, and purity (% by HPLC, NMR, or mp range). A 99% yield of 50% pure product is a 50% effective yield.

## Rules

- Always report purity alongside percent yield — isolated yield without purity data is misleading.
- A percent yield >100% always indicates an error: residual solvent, incomplete drying, or a weighing mistake.
- Use the limiting reagent for all yield calculations; using the excess reagent overstates theoretical yield.
- For multi-step syntheses, calculate overall yield as the product of individual step yields.

## Common Mistakes

- **Wrong limiting reagent** — failing to account for stoichiometric coefficients leads to incorrect theoretical yield.
- **Impure product weighed** — weighing before full drying or workup inflates actual yield; always dry to constant mass.
- **Using theoretical yield from excess reagent** — inflates denominator, artificially lowers reported % yield.
- **Ignoring catalyst stoichiometry** — for catalytic reactions, the catalyst is never the limiting reagent; limit by substrate.

## When NOT to Use

- For reactions where yield is defined differently (e.g., equilibrium conversions — use conversion and selectivity metrics instead)
- For biological enzyme assays (use activity units or specific activity instead)
- For radiochemical synthesis (use radiochemical yield and molar activity; mass may be too small to weigh)
