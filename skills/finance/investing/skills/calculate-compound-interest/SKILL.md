---
name: calculate-compound-interest
description: Use when projecting investment growth, understanding the time value of money, or illustrating the cost of delay — e.g., "how much will $10k grow?", "Rule of 72?", "how much do I need to save monthly to reach $1M?", "what's the cost of waiting 5 years to invest?"
source: Albert Einstein (attributed, origin disputed) on compound interest; Bernstein "The Four Pillars of Investing" (2002); Bogle "Common Sense on Mutual Funds" (1999); IRS future value tables
tags: [finance, investing, compound-interest, time-value-of-money, future-value, rule-of-72, growth-projection]
verified: true
---

# Calculate Compound Interest

Apply the future value formula and Rule of 72 to project investment growth, illustrate the cost of delay, and solve for required savings rates.

## Why This Is Best Practice

**Adopted by:** Compound interest calculations are foundational to every financial planning framework — used by CFPs, actuaries, and every personal finance curriculum (Dave Ramsey, Ramit Sethi, Vanguard investor education). The Rule of 72 is taught in every business school and CFA program as a mental math shortcut.
**Impact:** Compound interest is the mechanism behind every retirement plan — not investment selection. Vanguard research shows that time in market explains 88% of the variance in long-term wealth outcomes; investment selection explains less than 5%. Understanding compound growth quantitatively transforms abstract future goals into actionable present decisions.
**Why best:** Intuitive understanding of compound growth is systematically wrong — humans underestimate exponential growth. Calculating actual numbers (not approximate mental models) reveals that the cost of 5 years of delay is not 5 years of contributions but exponentially more — converting abstraction into urgency.

## Steps

1. **Future Value of a lump sum:**
   `FV = PV × (1 + r)^n`
   Where PV = present value, r = annual rate (decimal), n = years.
   Example: $10,000 at 8% for 30 years → $10,000 × (1.08)^30 = $100,627.

2. **Future Value of regular contributions (annuity):**
   `FV = PMT × [(1 + r)^n − 1] / r`
   Where PMT = monthly contribution (use r/12 and n×12 for monthly).
   Example: $500/month at 7% for 30 years → $566,765.

3. **Required monthly contribution to reach a goal:**
   `PMT = FV × r / [(1 + r)^n − 1]`
   Example: Need $1,000,000 in 30 years at 7% → $886/month.

4. **Apply the Rule of 72 for quick doubling estimates:**
   Years to double = 72 ÷ annual return %.
   At 6%: 72 ÷ 6 = 12 years to double.
   At 8%: 72 ÷ 8 = 9 years to double.
   At 10%: 72 ÷ 10 = 7.2 years to double.

5. **Calculate the cost of delay:**
   Start at 25 vs. 30, same $500/month at 7%: at 25 → $2.4M by 65; at 30 → $1.7M by 65.
   Cost of 5-year delay: $700,000 — not 5 × $6,000 ($30,000) in contributions, but $700,000 in forgone compounding.

6. **Account for inflation (real return):**
   Real return = nominal return − inflation. Use 7% nominal → 4–4.5% real (net of ~2.5% inflation).
   Future purchasing power: $1,000,000 nominal in 30 years = ~$475,000 in today's dollars.

7. **Model fee impact:**
   1% fee drag on $500/month at 8% for 30 years: fee scenario (7%) → $566,765 vs. no-fee (8%) → $680,239.
   Cost of 1% annual fee: $113,474 — illustrating why low-cost index funds matter.

## Rules

- Always use consistent compounding periods — monthly contributions need monthly rate (annual rate ÷ 12) and monthly periods.
- Use real (inflation-adjusted) returns for goal-setting; nominal returns for account value projections.
- The Rule of 72 is most accurate between 6–10% rates; at extremes (>20% or <3%) use the exact formula.
- Fees, taxes, and inflation each apply their own Rule of 72 in reverse — a 3% expense ratio doubles the cost of investing every 24 years.

## Examples

**How much to save to retire at 60 with $2M, starting at 30?**
30 years at 7% return.
PMT = $2,000,000 × (0.07/12) / [(1 + 0.07/12)^360 − 1] = $1,772/month.
At 7% real: $1,772/month adjusted for inflation.
Rule of 72 check: money doubles every 10.3 years (72÷7). From 30→60: three doublings. $1k today → $8k at 60. Need $2M ÷ 8 = $250k equivalent contribution base. Ballpark consistent.

## Common Mistakes

- **Using annual rate without adjusting for monthly compounding** — Applying 7% annual rate directly to monthly contribution calculations overstates growth; use 7%/12 = 0.583% per month.
- **Ignoring inflation** — $1M in 30 years at 3% inflation is worth $412,000 today. Goal-setting without inflation adjustment produces false confidence.
- **Treating projected returns as guaranteed** — 7–8% historical average includes years of −40%. The calculation is a planning tool, not a promise; actual results will deviate.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
