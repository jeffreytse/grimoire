---
name: design-synthesis-route
description: Use when planning a multi-step chemical synthesis from a target molecule back to available starting materials
source: 'Corey & Cheng "The Logic of Chemical Synthesis" (1989); Warren & Wyatt "Organic Synthesis: The Disconnection Approach" (2008); SciFinder/Reaxys synthetic planning methodology'
tags: [chemistry, synthesis, retrosynthesis, organic-chemistry]
verified: true
---

# Design Synthesis Route

Plan an efficient, feasible multi-step synthesis by working backward from target to available starting materials using retrosynthetic analysis.

## Why This Is Best Practice

**Adopted by:** Pharmaceutical process chemistry (FDA IND submissions), ACS Division of Organic Chemistry, Harvard/MIT/Caltech graduate synthesis courses, SciFinder and Reaxys route-planning algorithms.

**Impact:** Retrosynthetic analysis by Corey (Nobel Prize 1990) reduced average synthesis step count by 30–50% in complex natural product total synthesis; computer-assisted retrosynthesis (ASKCOS, ICSYNTH) achieves 80%+ accuracy for multi-step routes.

**Why best:** Working backward from target prevents forward-synthesis dead ends; disconnection approach systematically identifies bond-forming reactions, reducing reliance on intuition.

Sources: Corey & Cheng (1989) chapters 1–5; Warren & Wyatt 2nd ed. (2008); Clayden "Organic Chemistry" 2nd ed. (2012) ch. 30–31.

## Steps

1. **Draw the target molecule (TM)** — use skeletal structure; identify all functional groups, stereocenters, and ring systems.

2. **Assess complexity** — count carbons, stereocenters (n stereocenters → up to 2ⁿ stereoisomers), functional group count, and ring strain. Estimate required step count (~1–3 steps per complexity unit).

3. **Identify key disconnections** — apply retrosynthetic arrow (⟹) to break strategic bonds: bonds α to functional groups, bonds that simplify ring systems, bonds in the longest carbon chain.

4. **Apply transform library** — for each disconnection, identify the corresponding forward reaction: Diels-Alder (for 6-membered rings), aldol (for β-hydroxy carbonyls), Wittig/HWE (for alkenes), cross-coupling (for C–C/C–heteroatom bonds), etc.

5. **Generate synthons and reagent equivalents** — convert each synthon to a real reagent (e.g., acyl cation synthon → acyl chloride + Lewis acid; carbanion synthon → organolithium or Grignard).

6. **Evaluate each route** — score on: (a) step count (fewer = better), (b) overall yield estimate (product of step yields), (c) reagent cost and availability, (d) stereochemical control, (e) scalability and safety.

7. **Check literature precedent** — search SciFinder or Reaxys for each proposed step; confirm reaction conditions, yield range, and functional group compatibility. Prefer reactions with >70% reported yield.

8. **Write forward synthesis** — convert the retrosynthetic tree to a forward synthesis scheme with reagents, conditions, solvents, and expected yields per step.

9. **Identify protecting group strategy** — plan installation and removal of protecting groups for sensitive functional groups; use minimal number of protecting group operations.

10. **Flag safety and scale concerns** — identify hazardous reagents (organolithiums, azides, peroxides), exothermic steps, and steps unsuitable for scale-up.

## Rules

- Minimize protecting group steps — every PG operation adds 2 steps (install + remove) and reduces overall yield.
- Never plan a step without at least one literature precedent for the key bond-forming reaction.
- Check atom economy (Trost 1991) for each step — high atom economy reactions reduce waste and cost.
- Stereoselective steps must be planned with a specific stereochemical rationale, not assumed.

## Common Mistakes

- **Forward-only thinking** — starting from known materials and hoping to reach the target leads to circuitous routes with unnecessary steps.
- **Ignoring functional group compatibility** — planning a step that requires conditions incompatible with another functional group in the molecule.
- **Over-relying on protecting groups** — adding PG operations without checking if the reaction actually requires them.
- **Neglecting convergence** — linear synthesis amplifies step-count issues; convergent synthesis (assembling fragments) gives higher overall yield for complex targets.

## When NOT to Use

- For single-step reactions (direct synthesis from commercial materials) — just look up conditions
- For biocatalytic or biosynthetic routes where enzyme compatibility and cofactor availability drive planning
- For inorganic or organometallic synthesis where coordination chemistry rules, not disconnection approach
