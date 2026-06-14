---
name: design-cash-flow-forecast
description: Use when building a cash flow projection for a business — e.g., "will we run out of cash?", "how much runway do we have?", "13-week cash flow model", "cash flow forecast for fundraising"
source: Association for Financial Professionals (AFP) cash management guidelines; McKinsey "Cash is king" operations research; Sageworks/Pepperdine private capital market research; FP&A best practices (APQC)
tags: [finance, corporate, cash-flow, forecasting, financial-planning, runway, treasury]
verified: true
---

# Design Cash Flow Forecast

Build a 13-week rolling cash flow forecast to predict cash position, identify shortfalls early, and inform financing decisions.

## Why This Is Best Practice

**Adopted by:** The 13-week cash flow model is the standard for distressed company analysis (used by restructuring advisors Alvarez & Marsal, FTI Consulting) and for venture-backed startups tracking runway. CFOs at public companies maintain rolling 13-week models; boards require them in covenant-driven credit agreements.
**Impact:** CB Insights (2023) reports that 38% of startup failures cite "running out of cash" — nearly all of which is predictable 8–12 weeks in advance with a proper forecast. McKinsey research on corporate bankruptcies found that companies with a real-time cash forecast outsurvived liquidity crises at 3× the rate of those without one.
**Why best:** Income statement (P&L) lags reality — revenue is recognized before cash arrives, expenses are accrued before payment. Only a cash flow forecast shows when the bank account will be empty. 13 weeks is the standard horizon because it covers one quarter (enough for operational visibility) while remaining reliable (90-day projections have meaningfully more accuracy than annual).

## Steps

1. **Define the structure** — Three sections: (1) Cash receipts, (2) Cash disbursements, (3) Net cash flow and rolling balance. Weekly columns for 13 weeks. Build in Excel or Google Sheets with the balance auto-updating.
2. **Project cash receipts** — Map every cash inflow by timing of actual receipt (not invoice date):
   - Customer payments: use accounts receivable aging + historical DSO (Days Sales Outstanding). If DSO = 45 days, an invoice sent today arrives in week 6.
   - Subscription/recurring revenue: exact timing by billing date.
   - One-time receipts: grants, tax refunds, asset sales, investment rounds (date-certain only).
3. **Project cash disbursements** — Map every outflow by actual payment date:
   - Payroll: exact dates (bi-weekly vs. semi-monthly cycles matter).
   - Vendor payments: accounts payable aging + payment terms (net-30, net-60).
   - Rent, utilities, insurance: fixed date each month.
   - Debt service: exact payment dates.
   - Taxes: quarterly estimated payments, payroll taxes.
   - One-time: equipment, professional services, capital expenditures.
4. **Calculate net cash flow and ending balance** — Each week: Beginning balance + receipts − disbursements = ending balance. The ending balance becomes next week's beginning balance.
5. **Identify minimum cash threshold** — Operating minimum = 4–8 weeks of operating expenses (per AFP guidelines). Flag any week where projected balance falls below this threshold — that is a cash shortfall event requiring action.
6. **Add scenario analysis** — Best case (all receivables collected on time), base case (10% collection delay), stress case (30% delay + 10% revenue miss). Stress case reveals true liquidity exposure.
7. **Update weekly** — Replace projections with actuals every Monday. The model's value is in variance analysis: why did last week differ from forecast? Adjust forward projections accordingly.
8. **Use to drive financing decisions** — If stress case shows shortfall in week 8: begin credit line draw or fundraise now, not in week 7. Lenders and investors need 4–6 weeks minimum lead time.

## Rules

- Never confuse P&L with cash flow — profitable companies go bankrupt from cash flow timing mismatches; loss-making companies survive with strong cash management.
- Build from actual bank data and AR/AP aging, not from the income statement — P&L accruals mislead timing.
- The model is wrong; the process of updating it weekly is what creates value.
- Payroll is non-negotiable — model payroll dates exactly; missing payroll destroys a company faster than almost any other failure.

## Examples

**SaaS startup, $180k MRR, $150k/month burn, $500k cash:**
Week 1: $500k beginning. Receipts: $60k (week-1 of month billing). Disbursements: $75k (payroll $50k + vendors $25k). Ending: $485k.
Week 4: payroll again ($50k), AWS bill ($20k), rent ($8k). Net: −$78k week.
Week 13 projection (base case): $270k. Stress case (15% churn): $140k.
Runway: 3.4 months base, 2.0 months stress. Decision: begin Series A process immediately.

## Common Mistakes

- **Building monthly instead of weekly** — Monthly cash flow misses payroll timing and end-of-month AP bunching that can create mid-month crises invisible at monthly granularity.
- **Using revenue forecast instead of collections forecast** — A $500k Q4 pipeline doesn't help if none closes until December and payroll is in November.
- **Not updating with actuals** — A static forecast becomes fiction within 3 weeks. The discipline of weekly actuals is 80% of the model's value.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
