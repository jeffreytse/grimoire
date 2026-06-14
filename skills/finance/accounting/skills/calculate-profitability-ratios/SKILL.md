---
name: calculate-profitability-ratios
description: Use when evaluating a company's financial performance, comparing competitors, or assessing investment quality — e.g., "what's a good profit margin?", "how do I calculate ROIC?", "what does ROE tell me?", "interpreting gross vs. operating margin"
source: Graham & Dodd "Security Analysis" (1934); Damodaran "Applied Corporate Finance" (4th ed.); CFA Institute "Financial Statement Analysis"; McKinsey "Valuation" (Koller, Goedhart, Wessels, 7th ed.)
tags: [finance, accounting, profitability-ratios, ROE, ROIC, margins, financial-analysis, valuation]
verified: true
---

# Calculate Profitability Ratios

Compute and interpret the core profitability ratios — gross margin, operating margin, net margin, ROE, ROA, and ROIC — to assess business quality and capital efficiency.

## Why This Is Best Practice

**Adopted by:** Profitability ratio analysis is foundational in every CFA curriculum, MBA finance course, and investment banking training program. Warren Buffett famously screens for high ROE + low debt as his primary business quality filter. McKinsey's valuation framework centers ROIC (return on invested capital) as the single best predictor of long-term value creation.
**Impact:** McKinsey research shows companies with ROIC > WACC create value; those below destroy it, regardless of revenue growth. Consistent analysis of profitability ratios over time reveals business model trajectory — whether competitive advantage is strengthening or eroding — which neither revenue nor absolute profit shows alone.
**Why best:** Revenue growth is vanity; profitability is sanity; cash is reality. Margin ratios strip growth to reveal the underlying economics: a business growing at 30% but with declining gross margins is likely spending its way to growth. ROIC reveals whether reinvested capital generates returns above the cost of that capital — the definitive test of whether a business creates or destroys value.

## Steps

1. **Gross Profit Margin** — `(Revenue − COGS) ÷ Revenue`
   Measures how efficiently core products/services are produced.
   Benchmarks: SaaS 70–85%; retail 25–40%; manufacturing 20–35%; restaurants 60–70%.
   Trend matters: declining gross margin signals pricing pressure, input cost inflation, or mix shift to lower-margin products.

2. **Operating Profit Margin (EBIT Margin)** — `EBIT ÷ Revenue`
   EBIT = Gross Profit − Operating Expenses (R&D, S&M, G&A).
   Measures operational efficiency including overhead. Excludes financing (interest) and tax decisions.
   Benchmarks: mature SaaS 20–30%; consumer goods 10–15%; retail 3–8%; software company at scale 25–40%.

3. **Net Profit Margin** — `Net Income ÷ Revenue`
   Includes interest, tax, and one-time items. Most susceptible to accounting choices. Use EBIT margin for operating comparison; net margin for shareholder return analysis.

4. **Return on Equity (ROE)** — `Net Income ÷ Average Shareholders' Equity`
   Measures return on book value of equity. Buffett threshold: consistently > 15% without excessive leverage.
   DuPont decomposition: `ROE = Net Margin × Asset Turnover × Equity Multiplier (Leverage)`. High ROE from leverage is less valuable than high ROE from margins or turnover.

5. **Return on Assets (ROA)** — `Net Income ÷ Average Total Assets`
   Asset-efficiency measure. Useful for capital-intensive industries (manufacturing, banking). Benchmarks: banks 1–2%; tech 10–20%; manufacturing 5–10%.

6. **Return on Invested Capital (ROIC)** — `NOPAT ÷ Invested Capital`
   NOPAT = EBIT × (1 − tax rate). Invested Capital = Total Equity + Total Debt − Cash.
   The superior metric: measures return on all capital deployed, regardless of financing structure.
   ROIC > WACC = value creation. ROIC < WACC = value destruction even if profitable.
   High-quality benchmarks: Visa 35%+ ROIC; Apple 55%+ ROIC; commodity manufacturers 6–10% ROIC.

7. **Trend analysis (minimum 3 years)** — Calculate each ratio for 3–5 years. Ask: Are margins expanding (positive leverage) or compressing (cost pressure, competition)? Is ROIC converging toward or away from WACC? Compare to 3 closest peers.

8. **Identify ratio interactions** — High gross margin + low operating margin = cost problem in G&A or R&D. High net margin + low ROIC = capital-inefficient business. High ROIC + low ROE = under-leveraged, could use debt to enhance equity returns.

## Rules

- Always compare ratios to industry peers — a 5% net margin is excellent in retail, weak in SaaS.
- Adjust for non-recurring items before comparing — strip restructuring charges, write-downs, and one-time gains from operating margin.
- ROIC is more decision-useful than ROE for capital allocation — ROE can be inflated by leverage; ROIC cannot.
- A declining gross margin is the earliest warning sign of deteriorating competitive position; never ignore it.

## Examples

**Analysis of two competing software companies:**
Company A: 72% gross margin, 18% operating margin, ROIC 28%.
Company B: 65% gross margin, 12% operating margin, ROIC 14%.
WACC for both: 10%.
Verdict: Company A creates substantial value (ROIC 28% > WACC 10%); Company B creates marginal value. Company A's superior gross margin indicates stronger product differentiation or pricing power. Worth a premium multiple.

## Common Mistakes

- **Ignoring the DuPont decomposition** — A high ROE from 10× leverage is a red flag, not a quality signal; decompose before concluding.
- **Using net margin across industries** — Airlines at 2% net margin and software at 2% net margin are not equivalent; asset intensity, capital requirements, and risk profiles differ dramatically.
- **Comparing across accounting systems** — GAAP vs. non-GAAP margins are not directly comparable; adjust to a common basis before peer benchmarking.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
