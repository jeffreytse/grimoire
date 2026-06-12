---
name: write-financial-model
description: Use when building a three-statement financial model for a company or business unit
source: 'McKinsey "Valuation: Measuring and Managing the Value of Companies" (7th ed.); Macabacus financial modeling standards; BIWS (Breaking Into Wall Street) modeling curriculum'
tags: [finance, corporate, financial-modeling, three-statement, forecasting, sensitivity]
verified: true
---

# Write Financial Model

Construct an integrated three-statement financial model linking income statement, balance sheet, and cash flow statement.

## Why This Is Best Practice

**Adopted by:** Investment banks (Goldman Sachs, JP Morgan), management consulting firms (McKinsey, Bain), corporate finance teams at Fortune 500 companies
**Impact:** McKinsey "Valuation" research demonstrates that integrated three-statement models reduce forecasting errors by forcing accounting consistency; broken links between statements are the primary source of model errors in practice.

**Why best:** Three-statement integration ensures that every assumption flows consistently through all financial statements. A model where cash flow does not reconcile to the balance sheet is fundamentally unreliable. Sensitivity analysis built into the structure allows real-time scenario toggling without rebuilding the model.

## Steps

1. **Set up model architecture** — Separate tabs: Assumptions, Income Statement, Balance Sheet, Cash Flow Statement, Valuation, Sensitivity; color-code: blue = hardcoded inputs, black = formulas, green = links from other sheets.
2. **Build the income statement** — Start with revenue drivers (volume × price or segment-level); model COGS, gross margin, OpEx line items, EBITDA, D&A, EBIT, interest expense, taxes, net income.
3. **Model the balance sheet** — Link retained earnings to net income; model working capital (AR = DSO × Revenue / 365; Inventory = DIO × COGS / 365; AP = DPO × COGS / 365); model PP&E schedule (beginning + CapEx − D&A = ending).
4. **Build the cash flow statement** — Use indirect method: start with net income; add back D&A; adjust for working capital changes; subtract CapEx; add/subtract financing activities; ending cash must equal balance sheet cash.
5. **Create the debt and interest schedule** — Model revolving credit, term loans, and new issuances; interest expense feeds back to income statement; this is the circular reference — use iterative calculation or break with prior-period average balance.
6. **Build a scenarios tab** — Define at minimum three scenarios (base, upside, downside) by toggling revenue growth, margins, and CapEx assumptions; use INDEX-MATCH or scenario manager.
7. **Add sensitivity tables** — Two-way data tables for key outputs (EV, equity value, EPS) across revenue growth and EBITDA margin; mandatory before any presentation.
8. **Audit and stress-test** — Run error checks: does cash flow match balance sheet delta? Is retained earnings roll-forward correct? Run the downside scenario to confirm the model does not break.

## Rules

- Every hardcoded number must live in the Assumptions tab only; never bury inputs inside formulas.
- The balance sheet must balance (Assets = Liabilities + Equity) in every period; this is a non-negotiable integrity check.
- Avoid circular references where possible; if required (interest on average debt), enable iterative calculations with maximum iterations set to 100.
- Label units (thousands, millions) prominently; unit errors are a leading cause of model disasters.
- Protect formula cells; only input cells should be editable in shared models.

## Examples

**SaaS company model:** Revenue drivers = ARR × net revenue retention + new ARR bookings. Gross margin = 72% (stable). S&M = 40% of revenue declining to 30% over 5 years as CAC payback improves. R&D = 20% of revenue. G&A = 10% declining to 7%. EBITDA margin: year 1 = 2%, year 5 = 23%. D&A = 5% of revenue. CapEx = 3%. Working capital: DSO 45 days, DPO 30 days, no inventory.

## Common Mistakes

- **Hard-coding calculated values** — Never paste values over formulas; future scenario changes will silently break the model.
- **Ignoring minority interest and non-controlling stakes** — These affect equity value; model them explicitly if material.
- **Building a model no one else can audit** — Overcomplicated formulas defeat the purpose; if you cannot explain every cell, refactor.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
