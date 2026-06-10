---
name: calculate-startup-unit-economics
description: Use when computing or auditing a business's unit-level profitability including CAC, LTV, payback period, and contribution margin
source: Andreessen Horowitz (a16z) unit economics framework; Skok "SaaS Metrics 2.0" (2012); Bessemer Venture Partners cloud SaaS metrics guide
tags: [entrepreneurship, unit-economics, saas-metrics, finance, startups]
verified: true
---

# Calculate Unit Economics

Compute and interpret the unit-level profitability metrics that determine whether a business can scale sustainably.

## Why This Is Best Practice

**Adopted by:** a16z, Bessemer Venture Partners, Sequoia, SaaS founders, and marketplace operators
**Impact:** Bessemer's State of the Cloud report shows LTV:CAC ratio is the single metric most correlated with Series A/B funding success. Businesses with LTV:CAC > 3x are 4x more likely to reach $10M ARR than those below 1x. Skok found that SaaS companies that track payback period reduce churn by identifying at-risk cohorts 60 days earlier.
**Why best:** Unit economics reveal whether each incremental customer creates or destroys value, independent of growth rate. Unprofitable unit economics at small scale never improve at large scale — they accelerate losses.

Sources: Skok "SaaS Metrics 2.0" (2012, For Entrepreneurs); Bessemer Venture Partners "State of the Cloud" (2023); a16z "16 Startup Metrics" (2015)

## Steps

1. **Define the "unit"** — clarify what one unit is: one customer, one subscriber, one transaction, one driver, one seat. All calculations depend on consistent unit definition.

2. **Calculate Customer Acquisition Cost (CAC)** — divide total sales and marketing spend in a period by the number of new customers acquired in that period. Use fully-loaded costs: salaries, tools, ad spend, events.

   `CAC = Total S&M Spend / New Customers Acquired`

3. **Calculate Average Revenue Per Unit (ARPU)** — compute monthly recurring revenue (MRR) or transaction revenue per unit. For multi-tier products, calculate blended ARPU across all tiers.

4. **Calculate Gross Margin per unit** — subtract cost of goods sold (COGS) per customer from ARPU. COGS includes hosting, support, customer success, and third-party fees attributable to serving that customer.

   `Gross Margin % = (ARPU − COGS per customer) / ARPU`

5. **Estimate Customer Lifetime** — use one of: (1) 1 ÷ monthly churn rate, (2) cohort analysis of actual retention curves, or (3) contract length for B2B SaaS. Cohort analysis is most accurate.

6. **Calculate Lifetime Value (LTV)** — multiply gross margin per unit per month by customer lifetime in months.

   `LTV = (ARPU × Gross Margin %) / Monthly Churn Rate`

7. **Calculate LTV:CAC ratio** — the primary health indicator. Benchmark: > 3x is healthy, > 5x is excellent, < 1x means the business loses money on every customer.

   `LTV:CAC = LTV / CAC`

8. **Calculate CAC Payback Period** — how many months to recover the cost of acquiring one customer.

   `Payback Period = CAC / (ARPU × Gross Margin %)`   Benchmark: < 12 months for SaaS, < 6 months for marketplace.

9. **Segment by acquisition channel** — compute CAC separately for each channel (paid search, content, outbound, partner). Channels with 3x+ worse CAC than average are destroying unit economics even if volume looks good.

10. **Project at scale** — model whether CAC will increase or decrease as the business grows (most CAC increases with scale as cheap channels saturate). Project LTV improvement from expansion revenue (upsell, cross-sell). Determine if unit economics improve or deteriorate at 10x current scale.

## Rules

- Always use fully-loaded CAC — excluding salaries understates acquisition cost by 40–70%.
- Never blend new customer and expansion revenue in LTV without separating the cohorts.
- Payback period matters more than LTV for capital efficiency — a 36-month payback requires 3 years of working capital per customer.
- Churn must be calculated on a cohort basis, not as a period metric — period churn underestimates actual retention loss.
- Recalculate every quarter as product and go-to-market evolve.

## Common Mistakes

- **Excluding customer success from COGS** — high-touch B2B SaaS often has 30–50% of true COGS in customer success; excluding it inflates gross margin.
- **Using average churn instead of cohort churn** — growing companies have high proportions of new customers with low tenure, making average churn look better than actual cohort retention.
- **Ignoring CAC by channel** — blended CAC hides channels with negative unit economics subsidized by efficient channels.
- **Treating logo churn and revenue churn identically** — net revenue retention above 100% (expansion > churn) creates negative effective churn; this changes LTV calculation entirely.

## When NOT to Use

- Pre-revenue ideation stage (no data to calculate with — use assumptions explicitly labeled as hypotheses)
- Non-scalable service businesses without repeatable customer acquisition (use contribution margin per project instead)
- One-time transaction businesses without lifetime or repeat purchase component (use contribution margin per transaction instead)
