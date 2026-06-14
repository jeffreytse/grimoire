---
name: calculate-customer-lifetime-value
description: Use when calculating how much a customer is worth, pricing acquisition costs, or evaluating unit economics — e.g., "what's our LTV?", "how much can we spend to acquire a customer?", "LTV:CAC ratio?", "what's our payback period?"
source: Gupta, Lehmann & Stuart "Valuing Customers" (Journal of Marketing Research, 2004); SaaStr LTV/CAC framework; Andreessen Horowitz "16 Startup Metrics"; Harvard Business Review customer value research
tags: [finance, corporate, LTV, CLV, customer-lifetime-value, unit-economics, CAC, SaaS, payback-period]
verified: true
---

# Calculate Customer Lifetime Value

Compute LTV (lifetime value) and LTV:CAC ratio to evaluate unit economics, set acquisition budgets, and determine business model sustainability.

## Why This Is Best Practice

**Adopted by:** LTV:CAC is the primary unit economics metric used by VCs (Andreessen Horowitz, Sequoia), SaaS benchmarking firms (SaaStr, KeyBanc Capital Markets), and every CFO at a subscription business. It is a required metric for Series A+ investor presentations. Gupta et al. (2004) formalized the academic framework; the tech industry operationalized it as the dominant customer economics framework.
**Impact:** Companies that achieve LTV:CAC > 3 and payback period < 18 months grow sustainably; below these thresholds, growth destroys cash. CB Insights analysis of SaaS failures shows LTV:CAC misunderstanding (over-stated LTV or under-counted CAC) as a contributing factor in 60%+ of SaaS company failures.
**Why best:** LTV:CAC converts intuitive notions of "customer value" into a concrete capital allocation constraint. Without it, companies over-invest in acquisition (burning cash on unprofitable customers) or under-invest (leaving growth on the table when CAC payback is short). The ratio provides the decision rule: spend more if LTV:CAC > 3 and payback < 18 months; cut spend if below.

## Steps

1. **Calculate Average Revenue Per User (ARPU)** — Monthly ARPU = MRR ÷ total customers. For cohort analysis: use ARPU at time of acquisition, then model expansion/contraction.

2. **Calculate Gross Margin %** — LTV must be calculated on gross margin, not revenue. A $100/month customer at 30% gross margin has $30/month of contribution, not $100.
   Contribution margin/month = ARPU × gross margin %.

3. **Calculate Churn Rate** — Monthly churn = customers lost in month ÷ customers at start of month. Annual churn = 1 − (1 − monthly churn)^12. For cohort LTV, use cohort-specific retention curves.

4. **Calculate LTV (Simple formula):**
   `LTV = (ARPU × Gross Margin %) ÷ Monthly Churn Rate`
   Example: $100 ARPU × 70% gross margin = $70/month contribution. Monthly churn 3%. LTV = $70 ÷ 0.03 = $2,333.

5. **Calculate LTV (DCF formula for precision):**
   `LTV = Σ [Margin × (1 − Churn)^t ÷ (1 + d)^t]` for t = 1 to ∞
   Where d = discount rate (monthly). Simplifies to: `LTV = Margin / (Churn + Discount Rate)`.
   Use discount rate 10–15% annually (0.8–1.25% monthly) to reflect cost of capital.

6. **Calculate CAC** — Total sales & marketing spend in period ÷ new customers acquired in period. Include: salaries, commissions, ad spend, tools, events, allocated overhead. Common mistake: use only ad spend (underestimates true CAC by 2–5×).

7. **Calculate LTV:CAC ratio and payback period:**
   - LTV:CAC = LTV ÷ CAC. Target: ≥ 3 for SaaS/subscription; ≥ 5 is strong.
   - Payback period (months) = CAC ÷ (ARPU × Gross Margin %). Target: < 12 months for fast-growing SaaS; < 18 months for typical SaaS.

8. **Segment by customer type** — Enterprise vs. SMB vs. consumer LTVs differ dramatically. Calculate separately; avoid blended LTV that obscures which segment is profitable.

## Rules

- Always use gross margin, not revenue, in LTV — revenue-based LTV systematically overstates the unit economics.
- Fully-loaded CAC includes all S&M overhead; "blended CAC" that includes only ad spend is misleading for resource allocation.
- LTV is theoretical; actual realized value depends on churn not increasing as the business scales — validate with cohort data.
- If payback period > 18 months, growth requires external capital; if > 24 months, the model is structurally capital-intensive.

## Examples

**SaaS company metrics:**
ARPU: $200/month. Gross margin: 75%. Monthly churn: 2%.
Contribution/month: $200 × 75% = $150.
LTV = $150 ÷ 0.02 = $7,500.
CAC = $480k S&M spend ÷ 200 new customers = $2,400.
LTV:CAC = $7,500 ÷ $2,400 = 3.1. Acceptable.
Payback = $2,400 ÷ $150 = 16 months. Acceptable.
Action: can grow, but monitor churn — a rise from 2% to 3% drops LTV from $7,500 to $5,000 and LTV:CAC to 2.1 (below threshold).

## Common Mistakes

- **Using revenue instead of gross margin** — If COGS is 40%, revenue-based LTV overstates the real value by 67%. Investors will correct this; better to know it yourself first.
- **Ignoring expansion revenue** — NRR (Net Revenue Retention) > 100% means existing customers grow over time. Factor expansion MRR into LTV for businesses with upsell motions (Snowflake's NRR of 130% means LTV models excluding expansion dramatically understate value).
- **Static churn assumption** — Enterprise churn (0.5%/month) and SMB churn (3–5%/month) produce 6× different LTVs for the same ARPU. Segmented calculation is essential.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
