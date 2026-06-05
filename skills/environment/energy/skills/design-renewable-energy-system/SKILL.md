---
name: design-renewable-energy-system
description: Use when planning, sizing, or evaluating a renewable energy system (solar, wind, storage, or hybrid) for a facility, campus, or organization.
source: NREL (National Renewable Energy Laboratory) system sizing guidelines; IEA Renewable Energy guidelines; IEEE 1547-2018 interconnection standard
tags: [renewable-energy, solar, wind, storage, grid]
verified: true
---

# Design Renewable Energy System

Plan and size a renewable energy system using NREL methodology and IEEE interconnection standards to maximize generation, ensure grid compliance, and optimize economics.

## Why This Is Best Practice

**Adopted by:** NREL's SAM (System Advisor Model) used by 100,000+ users worldwide for renewable system design; IEEE 1547-2018 mandatory interconnection standard in all US states; IEA guidelines adopted by 40+ member countries; RE100 initiative (400+ companies committing to 100% renewable electricity); FERC Order 2222 enabling distributed energy resources

**Impact:** NREL reports utility-scale solar LCOE fell 89% from 2010–2023 ($0.36/kWh to $0.04/kWh); corporate PPA signatories save 10–20% vs. retail electricity rates; behind-the-meter solar with storage reduces peak demand charges by 30–50% (Rocky Mountain Institute 2020)

**Why best:** Proper resource assessment and system sizing using measured irradiance/wind data prevents both undersizing (unmet clean energy goals) and oversizing (stranded capital), while IEEE 1547 compliance prevents costly utility interconnection rejections.

Sources: NREL "Best Practices for Solar Photovoltaic System Design" (2021); IEA "Renewables 2023" and technology roadmaps; IEEE 1547-2018 "Standard for Interconnection and Interoperability of Distributed Energy Resources"; NREL System Advisor Model (SAM) documentation

## Steps

1. **Assess energy demand profile** — Analyze 12–24 months of interval energy data (15-min or hourly) to determine annual consumption (kWh), peak demand (kW), load shape by season and time of day, and self-consumption potential.

2. **Conduct resource assessment** — Use NREL's National Solar Radiation Database (NSRDB) for solar irradiance (GHI, DNI, DHI) or AWS Truepower / NREL Wind Toolkit for wind speed data at hub height. Minimum 10-year TMY (Typical Meteorological Year) dataset.

3. **Evaluate site constraints** — Assess available roof/ground area, structural capacity (roof loading for solar panels, typically 3–4 psf), shading analysis (using tools like Helioscope or PVsyst), geotechnical conditions for ground mount, setback requirements, and flood zone status.

4. **Select technology and configuration** — Choose between: rooftop solar PV, ground-mount solar, carport, floating solar, wind (small-scale <100kW or utility-scale), or hybrid (solar + storage, wind + solar + storage). Match technology to site resource and load profile.

5. **Size the system using simulation** — Use NREL SAM or PVsyst to model system performance: input panel specifications, inverter efficiency curves, system losses (soiling 1–3%, wiring 2%, mismatch 2%, shading), and tilt/azimuth. Output: annual generation (kWh/yr), specific yield (kWh/kWp), and performance ratio (target >80%).

6. **Design energy storage if required** — Size battery storage based on self-consumption target or backup duration. Calculate: storage capacity (kWh) = peak load (kW) × backup hours ÷ depth of discharge (DoD, typically 80–90%). Select chemistry: LFP for stationary storage (cycle life >4,000), NMC for space-constrained applications.

7. **Assess grid interconnection requirements** — Submit preliminary interconnection application to utility. IEEE 1547-2018 requires: anti-islanding protection, voltage and frequency ride-through settings, reactive power capability (for systems >200kW), and advanced inverter functions. Allow 3–18 months for utility review.

8. **Evaluate financing and incentives** — Calculate system economics including: federal Investment Tax Credit (ITC, currently 30% base + adders under IRA 2022), state incentives, net metering or NEM 3.0 tariffs, accelerated depreciation (MACRS 5-year for solar), and PPA/lease vs. ownership. Model NPV, IRR, and simple payback.

9. **Develop procurement and EPC specifications** — Define technical specifications for panels (efficiency ≥21% for monocrystalline), inverters (string vs. central vs. microinverter), racking, monitoring system (plant-level and module-level), and performance guarantees (P50/P90 generation, linear degradation ≤0.5%/yr).

10. **Commission and monitor** — Verify system performance against modeled P50 generation in first 12 months using SCADA/monitoring. Address underperformance >10% vs. model. Track: capacity factor (%), PR, specific yield, inverter availability, and soiling losses.

## Rules

- Always use measured TMY resource data (NSRDB or equivalent), not rule-of-thumb irradiance values — 10% irradiance error creates 10% sizing error, often worth $50K–$500K on commercial projects.
- IEEE 1547-2018 advanced inverter settings must be pre-configured to utility specifications before energization — post-commissioning changes require utility approval and cause delays.
- Battery storage systems must include a fire suppression and thermal management plan per NFPA 855 (US) before permitting.
- All generation models must report both P50 (median) and P90 (90% exceedance probability) estimates — only P90 is appropriate for debt financing.
- Interconnection agreement terms (export limits, curtailment rights) must be confirmed before finalizing system size — utility curtailment can eliminate 10–30% of projected revenue.

## Common Mistakes

- **Oversizing without self-consumption analysis** — installing more solar than annual consumption generates excess that earns retail or wholesale rates far below project assumptions, destroying project economics.
- **Ignoring degradation in financial models** — solar panels degrade 0.4–0.7%/yr; 25-year financial models using year-1 generation overstate revenue by 8–15%.
- **Skipping shading analysis** — a single chimney or rooftop HVAC unit can reduce string-level output by 20–50% if not mitigated by microinverters or DC optimizers.
- **Assuming net metering without verification** — net metering tariffs change; California NEM 3.0 reduced export credits by 75%, fundamentally changing system sizing and storage justification.

## When NOT to Use

- When site resource assessment shows annual solar irradiance <3.5 kWh/m²/day or annual average wind speed <6 m/s at hub height — economics rarely pencil without extraordinary incentives.
- When utility interconnection queue wait times exceed 3 years and project timeline is shorter — consider offsite PPA or RECs as an interim strategy.
- When a Power Purchase Agreement from an existing utility-scale renewable plant is available at lower LCOE than on-site development — off-site procurement may be more cost-effective than on-site generation.
