---
name: design-capital-structure
description: Use when advising on or designing the optimal capital structure for a company, including debt-equity mix, financing decisions, and leverage analysis
source: Modigliani-Miller theorems (1958, 1963); Myers "The Capital Structure Puzzle" JoF (1984); Brealey, Myers & Allen "Principles of Corporate Finance" (2019)
tags: [corporate-finance, capital-structure, leverage, financing]
verified: true
---

# Design Capital Structure

Design an optimal capital structure that minimizes the cost of capital, maximizes firm value, and maintains financial flexibility across economic cycles.

## Why This Is Best Practice

**Adopted by:** Capital structure theory underpins all corporate finance curricula (CFA, MBA programs); investment banks advise on capital structure in all M&A, IPO, and restructuring transactions; credit rating agencies (Moody's, S&P, Fitch) evaluate capital structure as the primary determinant of credit quality.
**Impact:** Optimal capital structure reduces WACC by 100–300bp vs. all-equity financing (via debt tax shield); excessive leverage increases probability of financial distress by 300% per Altman Z-score research; companies in industry median capital structure quartile outperform on 5-year TSR by 8–12%.
**Why best:** Modigliani-Miller (1963) proved that in a world with taxes, debt creates value through the interest tax shield — but too much debt creates financial distress costs that offset the benefit. The optimal structure balances these forces.

Sources: Modigliani & Miller (1958, 1963); Myers "Capital Structure Puzzle" Journal of Finance (1984); Myers & Majluf "Corporate Financing and Investment Decisions" JFE (1984); Brealey, Myers & Allen "Principles of Corporate Finance" 13th ed. (2019).

## Steps

1. **Analyze the business risk profile** — high-growth, R&D-intensive, and cyclical businesses carry high operating risk and should use less financial leverage (debt amplifies risk); stable, asset-heavy, predictable cash flow businesses can support more debt. Business risk is inversely related to optimal financial leverage.

2. **Calculate current capital structure** — compute: Debt/Equity ratio (D/E), Debt/EBITDA (leverage ratio), Interest Coverage Ratio (EBITDA/Interest Expense), Debt/Total Assets. Compare against industry medians from Damodaran's industry databases and peer group.

3. **Model the debt capacity** — estimate sustainable debt level where: (a) Interest Coverage Ratio remains ≥3.0x under a stress scenario (20% EBITDA decline), (b) Debt/EBITDA ≤ industry median, (c) the company maintains investment-grade credit quality or the target credit rating.

4. **Quantify the debt tax shield** — Value of Tax Shield = Corporate Tax Rate × Debt Level (under MM 1963). For a 25% tax rate and $100M of debt, the tax shield is worth $25M. Model the present value of the prospective interest tax shield using the company's pre-tax cost of debt.

5. **Estimate financial distress costs** — estimate the probability of financial distress (using Altman Z-score or credit model) and the cost of distress (typically 10–25% of firm value for direct and indirect costs). At the optimal structure, marginal tax shield benefit equals marginal distress cost.

6. **Apply the Pecking Order Theory (Myers 1984)** — companies prefer internal financing first, then debt, then equity. This predicts observed financing behavior: profitable companies use less debt (retained earnings available), growing companies use more debt (need external financing but avoid dilution). Validate against the company's financing history.

7. **Evaluate trade-off theory implications** — the static trade-off model identifies target leverage based on the tax-distress trade-off. Dynamic trade-off theory suggests companies rebalance capital structure gradually, not continuously. Set a target D/E range rather than a precise point.

8. **Model the WACC and EPS impact** — calculate WACC at current and proposed leverage levels (re-levering beta for each scenario). Model EPS impact of additional debt (interest expense reduces EBT but leverage amplifies EPS sensitivity to revenue changes). Present the WACC-minimizing leverage range.

9. **Assess market timing and refinancing risk** — optimal structure must be achievable in current credit markets. Assess: current credit market conditions (credit spreads, lender appetite), maturity profile (avoid cliff maturities), covenant flexibility, and refinancing risk under stress scenarios.

10. **Define the financing policy** — document: target leverage range (e.g., 2.0–2.5x Net Debt/EBITDA), credit rating target, dividend and buyback policy (excess cash return after maintaining target leverage), and conditions under which the company will deviate from the target.

## Rules

- Never set leverage based on peer comparison alone — peers may be sub-optimally capitalized; evaluate business risk independently.
- Maintain a liquidity buffer: minimum 12 months of operating cash needs accessible via cash or undrawn revolving credit.
- Stress-test the capital structure: if EBITDA declines 30%, can the company still service debt and remain covenant-compliant?
- Capital structure changes have signaling effects: debt increases signal confidence; equity issuance signals overvaluation (per Myers & Majluf).

## Common Mistakes

- **Over-leveraging cyclical businesses** — using average-cycle EBITDA to size debt in cyclical industries ignores trough scenarios where companies breach covenants and face distress.
- **Ignoring covenant constraints** — leverage ratios and coverage covenants constrain operational flexibility; tightest covenant (not the financial limit) is the binding constraint.
- **Treating debt and equity as symmetric** — debt comes with contractual obligations (interest, principal, covenants); equity is flexible. The asymmetry creates real operational constraints that financial models understate.
- **Forgetting off-balance-sheet obligations** — operating leases (post-ASC 842/IFRS 16 they are on-balance-sheet), pension liabilities, and contingent liabilities are economically debt-like and must be included in leverage calculations.

## When NOT to Use

- When the company is pre-revenue or early-stage (capital structure theory applies to cash-flow-generating entities; early-stage financing follows different logic — venture terms, not optimal leverage).
- When the company is in financial distress (focus on liquidity management and restructuring, not optimal capital structure).
- When the financing decision is a single project (use project finance structures rather than corporate-level capital structure optimization).
