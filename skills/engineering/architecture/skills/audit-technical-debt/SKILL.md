---
name: audit-technical-debt
description: Use when assessing, cataloging, or prioritizing technical debt in an existing codebase or system
source: Ward Cunningham (originator, 1992 OOPSLA); SQALE model (Letouzey & Coq, 2012); Martin Fowler "TechnicalDebtQuadrant" (martinfowler.com)
tags: [technical-debt, code-quality, refactoring, audit, architecture, sqale]
---

# Audit Technical Debt

Systematically identify, classify, and prioritize technical debt to make informed remediation decisions.

## Why This Is Best Practice

**Adopted by:** SonarQube (SQALE model), Google (internal debt tracking), Stripe engineering blog documents debt reduction programs
**Impact:** Google's internal studies found that unmanaged technical debt doubles feature development time within 18 months; SonarQube's SQALE index quantifies remediation cost per issue.

Ward Cunningham's metaphor frames debt as a deliberate trade-off, not a failure. Fowler's quadrant (deliberate/inadvertent × reckless/prudent) distinguishes debt worth paying from debt worth carrying. Auditing without this classification leads to either over-investment in low-value cleanup or ignored high-interest debt.

## Steps

1. **Run static analysis** — Use SonarQube, CodeClimate, or language-specific tools to surface code smells, duplication, complexity (cyclomatic), and security hotspots. Export the SQALE technical debt ratio.
2. **Classify by Fowler's quadrant** — For each significant debt item: is it deliberate or inadvertent? Prudent or reckless? Deliberate-prudent debt (conscious shortcuts) is highest priority to document; reckless debt is highest priority to fix.
3. **Assess interest rate** — How much does this debt slow down current work? High-churn files with high complexity are "high-interest" debt.
4. **Estimate remediation cost** — Use SQALE effort estimates or team estimation; group items into T-shirt sizes (S/M/L/XL).
5. **Prioritize by ROI** — Rank by: (interest rate × frequency touched) ÷ remediation cost. Fix high-interest, frequently-touched, cheap-to-fix debt first.
6. **Create a debt register** — Document each significant item: location, classification, estimated cost, owner, target resolution quarter.
7. **Allocate capacity** — Reserve 20% of sprint capacity (the "Boy Scout Rule" budget) for continuous debt reduction.

## Rules

- Never audit debt without also auditing business value — debt in a deprecated module should be deleted, not refactored.
- Document deliberate-prudent debt at the point of incurrence with a `TODO(DEBT):` comment and ticket reference.
- Treat the debt register as a living document — review quarterly.
- Do not attempt to pay all debt at once; targeted, incremental reduction outperforms big-bang rewrites.

## Examples

High-priority debt item:
- File: `src/checkout/payment_processor.py` — cyclomatic complexity 47, touched in 80% of sprints
- Classification: Inadvertent-reckless (no tests, no documentation)
- Remediation: Extract into PaymentGateway + PaymentValidator classes — estimated 3 days
- Priority: Critical (high interest, high churn, medium cost)

## Common Mistakes

- **Conflating all debt as bad** — deliberate-prudent debt (shipping an MVP shortcut) is a valid business decision.
- **Measuring only lines of code or coverage** — misses architectural debt (wrong abstractions, coupling) that static metrics don't surface.
- **No business buy-in** — auditing without budget allocation produces a list that is never acted upon.
