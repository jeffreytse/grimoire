---
name: apply-dimensional-analysis
description: Use when deriving relationships between physical quantities, checking equations for consistency, or scaling physical models to different sizes or conditions
source: Buckingham Pi theorem (Buckingham 1914); Bridgman "Dimensional Analysis" (1922); Barenblatt "Scaling, Self-Similarity, and Intermediate Asymptotics" (1996)
tags: [physics, dimensional-analysis, scaling, fluid-mechanics]
verified: true
---

# Apply Dimensional Analysis

Use the Buckingham Pi theorem to derive dimensionless groups, check dimensional consistency, and scale physical relationships without solving governing equations.

## Why This Is Best Practice

**Adopted by:** Fluid mechanics (Reynolds number, Mach number, Froude number), heat transfer (Nusselt, Prandtl, Rayleigh numbers), structural engineering (geometric similarity), NASA aerodynamic model testing, chemical engineering (Damköhler number).

**Impact:** Dimensional analysis reduced experimental effort in the Wright Brothers' wind tunnel tests by ~70% by identifying similarity parameters; Kolmogorov's dimensional analysis of turbulence (1941) produced the -5/3 energy cascade law without solving the Navier-Stokes equations.

**Why best:** Dimensional analysis is model-independent — it constrains the form of physical laws using only the principle that physical equations must be dimensionally homogeneous. It reveals scaling laws and collapses data from different experimental conditions onto universal curves.

Sources: Buckingham Phys Rev 4:345–376 (1914); Bridgman (1922); Barenblatt (1996) ch. 1–3; White "Fluid Mechanics" 8th ed. (2016) ch. 5.

## Steps

1. **List all relevant physical variables** — identify every variable that could affect the phenomenon: n variables total (e.g., for drag: force F, velocity v, density ρ, viscosity μ, length L → n=5).

2. **Write dimensions for each variable** — use fundamental dimensions: M (mass), L (length), T (time), θ (temperature), I (electric current). Express each variable as M^a L^b T^c ... (e.g., dynamic viscosity μ: M L⁻¹ T⁻¹).

3. **Count independent dimensions** — determine k = number of independent fundamental dimensions appearing (usually k = 3 for mechanical problems: M, L, T).

4. **Apply Buckingham Pi theorem** — the number of independent dimensionless groups is p = n − k. These groups (π₁, π₂, ..., πₚ) completely describe the phenomenon.

5. **Choose k repeating variables** — select k variables that together contain all fundamental dimensions and are not dimensionless themselves. These will appear in every π group. Good choices: ρ, v, L for fluid mechanics.

6. **Form each π group** — for each remaining variable, combine it with the k repeating variables raised to unknown powers and solve the exponent system to make the product dimensionless.

7. **Simplify and name** — check if each π group matches a named dimensionless number (Re, Fr, Nu, etc.); if so, use the standard name and definition to ensure comparability with literature.

8. **Write the functional relationship** — express: π₁ = f(π₂, π₃, ..., πₚ). For p=1, the result is π₁ = constant (a complete similarity law). For p=2, π₁ = f(π₂) — a curve to be determined experimentally.

9. **Apply to scaling** — if two systems have equal values of all π groups, they are dynamically similar and will exhibit the same behavior (a model can predict the full-scale system).

10. **Verify dimensional consistency** — check your final equation/relationship: every additive term must have identical dimensions; this is a necessary (not sufficient) condition for a correct equation.

## Rules

- The choice of repeating variables affects the form of π groups but not their number or the physical content — choose for algebraic convenience.
- Dimensional analysis cannot determine numerical constants (e.g., the drag coefficient CD must be measured); it only determines the functional form.
- All variables that matter must be included; missing a relevant variable gives an incorrect set of π groups.
- Dimensional analysis cannot distinguish between different functional forms with the same dimensionless groups (e.g., Re¹ vs. Re²); experiments or theory are needed to resolve this.

## Common Mistakes

- **Including too many variables** — adding irrelevant variables increases p and obscures the result; include only physically motivated quantities.
- **Missing a variable** — data will not collapse when plotted against the dimensionless groups if a relevant variable is absent.
- **Using dependent dimensions** — if M, L, T, and F (force) are all used as fundamental dimensions, F = MLT⁻² is a dependency; use only independent fundamental dimensions (k decreases by 1).
- **Confusing geometric and dynamic similarity** — equal π groups (dynamic similarity) is necessary for model-to-prototype scaling; geometric similarity alone is insufficient.

## When NOT to Use

- When all variables and governing equations are known and can be solved analytically (use direct solution)
- For quantum mechanical phenomena where the action quantum ℏ introduces a natural scale that makes classical dimensional analysis incomplete
- For systems with many relevant variables (n−k > 4) where dimensional analysis is still valid but experimental mapping of the π space becomes impractical
