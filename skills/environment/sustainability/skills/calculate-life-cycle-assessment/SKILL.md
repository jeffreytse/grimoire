---
name: calculate-life-cycle-assessment
description: Use when evaluating the environmental impacts of a product, process, or service across its entire life cycle from raw material extraction to end of life.
source: ISO 14040:2006 and ISO 14044:2006 LCA standards; PRé Sustainability SimaPro methodology; European Commission ILCD Handbook (JRC 2010)
tags: [lca, life-cycle, environmental-impact, product, sustainability]
verified: true
---

# Calculate Life Cycle Assessment

Quantify and compare the environmental impacts of a product or service across its full life cycle using ISO-compliant methodology.

## Why This Is Best Practice

**Adopted by:** European Commission (mandatory for EU Environmental Product Declarations, EU Taxonomy green investments); US DOE, EPA, and NIST product sustainability programs; leading consumer goods companies (Procter & Gamble, Unilever, BASF) for eco-design; Global Footprint Network; Product Environmental Footprint (PEF) category rules

**Impact:** LCA-guided redesign reduced cradle-to-gate emissions by 15–35% in documented case studies (ecoinvent database meta-analysis); Interface Inc. used LCA to achieve 96% reduction in product carbon footprint over 25 years; automotive lightweighting decisions validated by LCA deliver 6–8% lifetime GHG reduction per 10% weight reduction

**Why best:** LCA prevents burden-shifting (improving one life stage at the cost of another), reveals hotspots invisible to single-metric approaches, and provides legally defensible evidence for environmental claims under ISO 14021/14025.

Sources: ISO 14040:2006 "Environmental management — Life cycle assessment — Principles and framework"; ISO 14044:2006 "Requirements and guidelines"; EC JRC "International Reference Life Cycle Data System (ILCD) Handbook" (2010); ecoinvent v3 background database

## Steps

1. **Define goal and scope** — State the study purpose, intended audience, and whether results will be used for comparative assertions disclosed to the public (triggers critical review requirement). Define functional unit (e.g., "1 kg of protein delivered to consumer").

2. **Establish system boundary** — Specify life cycle stages included: cradle-to-grave (full), cradle-to-gate (production only), gate-to-gate (processing), or cradle-to-cradle (with end-of-life recovery). Document exclusions with justification; excluded flows must be <1% of total impact in each category.

3. **Build process flow diagram** — Map all unit processes (extraction, manufacturing, transport, use, disposal) and their material and energy linkages. Identify cut-off points and multifunctional processes.

4. **Collect foreground inventory data** — Gather primary data from the product system: bill of materials, energy consumption per process, water use, direct emissions, transport distances and modes, and waste generated.

5. **Select background database** — Use ecoinvent v3, GaBi, or USLCI for background processes (electricity grids, commodity chemicals, transport). Match geographic and temporal scope to the product system.

6. **Handle multifunctionality** — When a process produces multiple outputs, apply in priority order: system expansion (preferred), substitution, or allocation by physical or economic parameter. Document choice.

7. **Compile Life Cycle Inventory (LCI)** — Aggregate all elementary flows (emissions to air, water, soil; resource extractions) across the system boundary using LCA software (SimaPro, OpenLCA, GaBi).

8. **Calculate Life Cycle Impact Assessment (LCIA)** — Apply characterization factors from a recommended method: ReCiPe 2016, EF 3.1 (EU Product Environmental Footprint), CML-IA, or TRACI (North America). Minimum categories: climate change, ozone depletion, acidification, eutrophication, resource depletion.

9. **Interpret results and conduct sensitivity analysis** — Identify hotspots (>20% contribution to any impact category); test sensitivity of key assumptions (e.g., ±20% on electricity carbon intensity, alternative end-of-life scenarios).

10. **Peer review and communicate** — Commission critical review by independent LCA expert per ISO 14044 §6 for comparative assertions or external disclosure; prepare Environmental Product Declaration (EPD) per ISO 14025/EN 15804 if required.

## Rules

- The functional unit must be measurable and reflect actual use — never use vague units like "per product."
- Data quality must be assessed using the ecoinvent Data Quality Guidelines or equivalent (geographic, temporal, technological representativeness).
- Never draw conclusions beyond what the defined scope supports — a cradle-to-gate study cannot compare product end-of-life performance.
- All comparative LCAs intended for public disclosure require independent critical review per ISO 14044.
- Report all impact categories calculated — cherry-picking favorable categories is misleading and violates ISO 14044 reporting requirements.

## Common Mistakes

- **Narrow system boundary to favor one product** — excluding upstream supply chain impacts or use-phase energy inflates comparative performance; ISO 14044 requires transparent justification of all exclusions.
- **Using wrong functional unit** — comparing 1 kg of Product A vs. 1 kg of Product B when they deliver different service durations (e.g., conventional vs. LED bulbs) produces meaningless results.
- **Treating all electricity as identical** — using average national grid without considering actual energy source or time of use overstates or understates impacts, especially for renewable-powered facilities.
- **Ignoring allocation in co-products** — applying economic allocation to waste streams or co-products without testing physical allocation as an alternative hides methodological sensitivity.

## When NOT to Use

- When a rapid screening tool is sufficient — use a spend-based or hotspot screening before committing to a full LCA.
- When primary data collection is not feasible and background database coverage is poor for the product system — results will have unacceptably high uncertainty.
- When the decision can be made based on a single dominant impact and data for that impact category alone is available — targeted carbon footprint (ISO 14067) may be more efficient.
