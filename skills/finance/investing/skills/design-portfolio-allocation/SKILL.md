---
name: design-portfolio-allocation
description: Use when constructing or rebalancing an investment portfolio across asset classes
source: Markowitz (1952) "Portfolio Selection" Journal of Finance; CFA Institute Portfolio Management curriculum; Vanguard Research on asset allocation
tags: [finance, investing, portfolio, asset-allocation, diversification, risk]
verified: true
---

# Design Portfolio Allocation

Construct a risk-adjusted portfolio allocation grounded in Modern Portfolio Theory.

## Why This Is Best Practice

**Adopted by:** CFA Institute, institutional asset managers (Vanguard, BlackRock, endowments), pension funds globally
**Impact:** Vanguard research shows asset allocation explains ~90% of long-run return variability; Markowitz's mean-variance optimization won the 1990 Nobel Prize in Economics.

**Why best:** Diversification across uncorrelated asset classes reduces portfolio volatility without proportionally reducing expected returns — the only "free lunch" in finance. A structured allocation process prevents behavioral errors (chasing performance, panic-selling) by anchoring decisions to an Investment Policy Statement.

## Steps

1. **Establish risk tolerance** — Quantify willingness (questionnaire) and ability (time horizon, liquidity needs, income stability) to take risk; these may conflict — ability dominates.
2. **Write an Investment Policy Statement (IPS)** — Document return objective, risk tolerance, liquidity requirements, time horizon, tax situation, and any constraints; revisit annually.
3. **Select asset classes** — Include only asset classes that improve risk-adjusted return (Sharpe ratio); common palette: domestic equities, international equities (developed + EM), bonds (government + corporate), real assets (REITs, commodities), alternatives.
4. **Estimate capital market assumptions** — Use 10-year expected returns, volatilities, and correlations from reputable sources (CFA Institute, Research Affiliates, Vanguard); avoid using raw historical returns as forecasts.
5. **Run mean-variance optimization** — Solve for the efficient frontier; select the portfolio that maximizes Sharpe ratio or matches the target volatility; use resampled optimization or Black-Litterman to avoid corner solutions.
6. **Apply constraints** — Set minimum and maximum weights per asset class; require minimum diversification (no single asset > 30% unless justified).
7. **Define rebalancing policy** — Use threshold-based rebalancing (e.g., ±5% drift triggers rebalance) rather than calendar-based; minimizes unnecessary turnover.
8. **Implement and document** — Select low-cost index funds or ETFs for each allocation bucket; document implementation rationale.

## Rules

- Never allocate based on recent performance alone; mean reversion is a documented phenomenon.
- Always match portfolio duration and liquidity to the investor's time horizon.
- Maintain at least 3–6 months of expenses in cash or near-cash before investing.
- Rebalance systematically; do not let drift compound over multiple years.
- Tax-locate assets correctly: hold bonds in tax-deferred accounts, equities in taxable accounts where possible.

## Examples

**Moderate-risk investor, 30-year horizon:** 60% global equities (40% US, 20% international), 30% bonds (20% US aggregate, 10% TIPS), 10% real assets (5% REITs, 5% commodities). Expected return ~7%, volatility ~10%, Sharpe ~0.55. Rebalance when any allocation drifts >5% from target.

## Common Mistakes

- **Home-country bias** — Overweighting domestic equities reduces diversification; global market-cap weighting is the neutral starting point.
- **Ignoring correlation changes in crises** — Correlations rise during market stress; hold true diversifiers (Treasuries, gold) not just high-return assets.
- **Optimizing on historical data without adjustment** — Garbage in, garbage out; always adjust capital market assumptions for current valuations.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
