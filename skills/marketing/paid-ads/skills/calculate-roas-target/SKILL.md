---
name: calculate-roas-target
description: Use when setting ROAS (Return on Ad Spend) targets for paid advertising campaigns based on business economics and margin requirements
source: Google Ads Smart Bidding methodology; Meta advertising ROAS optimization guidelines; Marketing Science "Return on Advertising Spend" (IPA Databank research)
tags: [roas, paid-ads, roi, bidding, performance-marketing]
verified: true
---

# Calculate ROAS Target

Derive the minimum and target ROAS required for a paid advertising campaign to be profitable, given business margins, customer economics, and channel objectives.

## Why This Is Best Practice

**Adopted by:** Google Smart Bidding Target ROAS strategy, Meta Advantage+ campaign optimization, all performance marketing agencies
**Impact:** Campaigns with correctly set ROAS targets achieve 30–50% better efficiency vs. manually managed bids (Google Smart Bidding case studies); incorrect ROAS targets are the #1 cause of profitable campaigns being paused or scaled incorrectly
**Why best:** ROAS alone is meaningless without grounding it in business economics — a 300% ROAS may be highly profitable for one product and loss-making for another. The calculation forces alignment between marketing and finance.

Sources: Google "Smart Bidding Best Practices" (2023); Binet & Field "Media in Focus" IPA (2017); IPA Databank long-term effectiveness research

## Steps

1. **Gather input variables** — collect: average order value (AOV), gross margin %, customer acquisition cost (CAC) target, and LTV if repeat purchase business.
2. **Calculate break-even ROAS** — Break-even ROAS = 1 ÷ Gross Margin. Example: 40% margin → break-even ROAS = 1 ÷ 0.40 = 2.5 (250%). This is the floor; operating below it means every sale loses money.
3. **Add overhead allocation** — if blended overhead (fulfillment, customer service, returns, platform fees) averages 15% of revenue, add to margin cost: effective margin = gross margin - overhead rate. Recalculate break-even ROAS.
4. **Set profit target margin** — decide the target net profit percentage from paid ads (e.g., 10% net profit target). Adjust ROAS target: target ROAS = 1 ÷ (gross margin - overhead - target net margin).
5. **Adjust for LTV multiplier** — for subscription or repeat-purchase businesses: ROAS can be evaluated on LTV not first-order AOV. ROAS_LTV = ROAS_transaction × (LTV ÷ AOV). This allows lower first-purchase ROAS.
6. **Differentiate by campaign type** — brand campaigns (high intent, lower CAC) warrant higher ROAS targets; prospecting campaigns (lower intent, higher CAC) warrant lower ROAS targets; retargeting sits between.
7. **Set minimum, target, and stretch ROAS** — minimum = break-even; target = profitable at margin goal; stretch = ideal efficiency target for budget scaling decisions.
8. **Communicate ROAS as a bidding signal** — input target ROAS into Google/Meta Smart Bidding; the algorithm will optimize bids to hit this target across the portfolio of auctions.
9. **Monitor ROAS by segment** — track ROAS by: campaign type, audience, geography, device, and product category; average ROAS masks high and low performers that warrant different targets.
10. **Review monthly** — recalculate when AOV, margin, or overhead changes; seasonal ROAS variation is normal and targets should be seasonally adjusted.

## Rules

- Break-even ROAS is the non-negotiable floor — operating below it consistently means paid advertising is destroying margin.
- ROAS targets must vary by campaign stage — new customer acquisition warrants a lower ROAS target than retargeting.
- Never evaluate ROAS in isolation — pair with volume (impressions, clicks, conversions) to ensure the ROAS target is not too restrictive to generate meaningful volume.
- LTV-adjusted ROAS is only valid if retention metrics are reliably measured — using assumed LTV without actual data inflates targets.
- Blended ROAS (all paid channels combined) is a management metric; channel ROAS is the optimization metric.

## Common Mistakes

- **Using revenue ROAS without margin adjustment** — a 400% ROAS on a 20% margin product is only 80% gross margin ROAS — loss-making.
- **Setting one ROAS target for all campaigns** — brand, prospecting, and retargeting have structurally different economics; one target misfires for at least two of them.
- **Chasing ROAS at the expense of volume** — a highly restrictive ROAS target (e.g., 800%) may be technically achieved but at such low volume that it is irrelevant.
- **Not adjusting for seasonality** — Q4 ROAS naturally spikes for e-commerce; a static annual target prevents appropriate budget scaling during peak.
- **Confusing ROAS with ROI** — ROAS is revenue ÷ ad spend; ROI accounts for all costs. Always calculate both.

## When NOT to Use

- Brand awareness campaigns with no direct conversion goal — use reach and CPM efficiency metrics instead
- Content marketing or organic channels — ROAS applies to paid advertising only
- Early-stage testing with insufficient conversion data to set reliable ROAS benchmarks
