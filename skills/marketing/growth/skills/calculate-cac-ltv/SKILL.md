---
name: calculate-cac-ltv
description: Use when calculating customer acquisition cost and lifetime value to evaluate channel efficiency or business health
source: David Skok "SaaS Metrics 2.0" (forentrepreneurs.com), Bessemer Venture Partners SaaS Benchmarks, Bill Gurley "Dangerous Seduction of the Lifetime Value Formula"
tags: [growth, saas-metrics, cac, ltv, unit-economics, financial-modeling, retention]
verified: true
---

# Calculate CAC and LTV

Accurately compute Customer Acquisition Cost and Lifetime Value to assess unit economics and channel viability.

## Why This Is Best Practice

**Adopted by:** Bessemer Venture Partners uses LTV:CAC as a primary SaaS health benchmark; David Skok's "SaaS Metrics 2.0" is the canonical reference for SaaS finance
**Impact:** Bessemer benchmarks: healthy SaaS requires LTV:CAC ≥ 3x and CAC payback period ≤ 12 months; companies below these thresholds face capital efficiency problems

**Why best:** CAC and LTV are the twin pillars of sustainable growth. Miscomputing either produces a false signal: under-counting CAC justifies overspending on acquisition; over-estimating LTV justifies unsustainable CAC. The ratio and payback period together determine how aggressively a business can invest in growth.

## Steps

1. **Define the cohort and time period** — Select a cohort (e.g., customers acquired in Q1 2025) and a time window (monthly, quarterly, annual). All inputs must use the same period.
2. **Calculate fully-loaded CAC** — Sum all sales and marketing expenses in the period (salaries, tools, ad spend, events, agency fees). Divide by the number of new customers acquired in the same period. `CAC = Total S&M Spend ÷ New Customers`.
3. **Segment CAC by channel** — Calculate blended CAC first, then break out by channel (paid search, content, referral, sales) to identify efficient vs. inefficient channels.
4. **Calculate ARPU** — Average Revenue Per User per month = Total MRR ÷ Total Active Customers.
5. **Calculate gross margin** — Gross Margin % = (Revenue − COGS) ÷ Revenue. Use this to compute contribution margin per customer.
6. **Calculate churn rate** — Monthly Churn = Churned Customers ÷ Customers at Start of Month. For SaaS, also compute Net Revenue Retention (NRR).
7. **Calculate LTV** — Simple model: `LTV = ARPU × Gross Margin % ÷ Monthly Churn Rate`. Advanced: use cohort retention curves and discounted cash flows for precision.
8. **Compute LTV:CAC ratio and payback period** — Ratio = LTV ÷ CAC. Payback months = CAC ÷ (ARPU × Gross Margin %). Compare to Bessemer benchmarks: ratio ≥ 3x, payback ≤ 12 months.

## Rules

- Always use fully-loaded CAC including salaries and overhead — marketing-only CAC is misleading and produces false confidence.
- Never use simple average churn for high-growth companies — cohort churn curves reveal the actual retention shape.
- Segment LTV by acquisition channel and customer tier; blended LTV hides wide variance.
- Treat LTV as a probabilistic estimate, not a fact — Bill Gurley warns that companies optimize for LTV models they built themselves.

## Examples

SaaS product: S&M spend = $120k/quarter, new customers = 200. CAC = $600. ARPU = $80/mo, gross margin = 75%, monthly churn = 2.5%. LTV = $80 × 0.75 ÷ 0.025 = $2,400. LTV:CAC = 4x. Payback = $600 ÷ ($80 × 0.75) = 10 months. Result: healthy unit economics; growth investment is justified.

## Common Mistakes

- Using marketing spend only in CAC denominator — excludes sales team cost, overstating efficiency by 2-3x.
- Using aggregate churn not cohort churn — masks that older cohorts churn faster; produces optimistic LTV.
- Comparing LTV to CAC without gross margin adjustment — comparing revenue to cost, not profit to cost.
- Ignoring CAC payback period — a 5x LTV:CAC ratio with 36-month payback can still cause a cash flow crisis.

## When NOT to Use

- Do not calculate LTV for a business with fewer than 6 months of retention data — there is insufficient cohort history to produce a statistically valid churn rate, making the LTV figure speculative.
- Do not apply this skill to marketplace or transactional businesses with highly irregular purchase frequency (e.g., real estate, weddings) where the standard ARPU ÷ churn model does not reflect actual revenue patterns.
- Do not use blended CAC/LTV analysis to justify channel-level budget decisions — channel-segmented unit economics are required; blended numbers can mask a profitable channel subsidizing an unprofitable one.
