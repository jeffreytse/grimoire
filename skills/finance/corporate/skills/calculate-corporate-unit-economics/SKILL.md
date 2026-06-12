---
name: calculate-corporate-unit-economics
description: Use when evaluating the profitability and scalability of a business model at the per-customer or per-transaction level
source: David Skok's SaaS Metrics 2.0 (forentrepreneurs.com); a16z SaaS metrics framework; Bessemer Venture Partners Cloud Index methodology
tags: [finance, corporate, saas, unit-economics, ltv, cac, payback-period]
verified: true
---

# Calculate Unit Economics

Measure per-unit profitability to determine whether a business model is viable and scalable.

## Why This Is Best Practice

**Adopted by:** Y Combinator, a16z, Bessemer Venture Partners; required analysis in Series A+ due diligence; standard in SaaS and subscription businesses
**Impact:** David Skok's research shows that SaaS businesses with LTV:CAC < 3x consistently fail to reach profitability at scale; a16z portfolio data confirms CAC payback period is the single strongest predictor of capital efficiency.

**Why best:** Unit economics reveal whether growth is value-creating or value-destroying before the aggregate P&L shows the problem. A business acquiring customers profitably at small scale can raise prices and invest in growth; one with broken unit economics cannot be fixed by scaling.

## Steps

1. **Define the unit** — For SaaS: one customer. For marketplace: one transaction or one supplier. For e-commerce: one order. Be explicit; unit economics are meaningless without a clear unit definition.
2. **Calculate Customer Acquisition Cost (CAC)** — CAC = (Total Sales + Marketing Spend in period) / New Customers Acquired in same period; use fully-loaded cost including headcount, tools, and overhead allocated to S&M.
3. **Segment blended vs. channel CAC** — Calculate CAC by acquisition channel (paid search, content, referral, outbound); blended CAC hides which channels are profitable.
4. **Calculate Lifetime Value (LTV)** — For SaaS: LTV = (Average Revenue Per Account × Gross Margin %) / Customer Churn Rate. For transactional: LTV = Average Order Value × Gross Margin × Purchase Frequency × Retention Horizon.
5. **Calculate LTV:CAC ratio** — Target ≥ 3x for SaaS; < 1x means the business loses money on every customer indefinitely; 1–3x is marginal and requires improvement.
6. **Calculate CAC Payback Period** — Payback = CAC / (Monthly Recurring Revenue per customer × Gross Margin %); target < 12 months for efficient SaaS, < 18 months acceptable, > 24 months requires explanation.
7. **Model churn impact** — Recalculate LTV at 1.5× and 2× current churn; show how sensitive the business is to churn changes; net revenue retention (NRR) > 100% dramatically improves LTV.
8. **Project at scale** — Model CAC efficiency improvements from organic/referral channel growth; confirm LTV:CAC and payback improve, not worsen, at scale.

## Rules

- Always use gross margin (not revenue) in LTV calculations; ignoring COGS overstates unit profitability.
- Never average LTV:CAC across cohorts without segmenting by acquisition vintage; early cohorts are typically more profitable.
- Report CAC payback on a gross margin basis, not revenue basis.
- Include customer success and onboarding costs in CAC for enterprise SaaS; these are acquisition costs if they precede revenue recognition.
- Benchmark against industry standards: consumer SaaS payback < 6 months, B2B SaaS < 18 months, enterprise SaaS < 24 months.

## Examples

**B2B SaaS example:** Sales + marketing spend = $500,000/quarter. New customers = 50. CAC = $10,000. Average MRR/customer = $1,500. Gross margin = 70%. Monthly gross profit per customer = $1,050. CAC payback = $10,000 / $1,050 = 9.5 months. Monthly churn = 1.5%. LTV = $1,050 / 0.015 = $70,000. LTV:CAC = 7x. Excellent unit economics.

## Common Mistakes

- **Using revenue instead of gross profit in LTV** — A 30% gross margin business with "great" revenue-based LTV:CAC may actually have terrible unit economics.
- **Ignoring churn when modeling growth** — Leaky bucket: if monthly churn is 3%, you lose 30% of the base per year; growth targets must account for replacement of churned revenue.
- **Blending enterprise and SMB customers** — These cohorts have fundamentally different CAC, LTV, and churn; always segment before comparing to benchmarks.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
