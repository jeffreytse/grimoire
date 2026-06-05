---
name: audit-marketing-funnel
description: Use when diagnosing drop-off points, conversion bottlenecks, and performance gaps across marketing and sales funnel stages
source: HubSpot "State of Marketing" annual report; Forrester customer journey research; Avinash Kaushik "Web Analytics 2.0" (2009)
tags: [funnel, analytics, conversion, marketing, audit]
verified: true
---

# Audit Marketing Funnel

Systematically measure, diagnose, and prioritize improvements across each funnel stage — from awareness to revenue — to maximize conversion rate and revenue per visitor.

## Why This Is Best Practice

**Adopted by:** HubSpot flywheel methodology, Salesforce revenue operations framework, CXL Institute CRO methodology
**Impact:** Funnel audits identify conversion improvements averaging 15–30% within 90 days (CXL Institute benchmark); fixing the worst-converting stage produces compounding improvement on all downstream stages
**Why best:** Kaushik's "See-Think-Do-Care" framework and funnel analysis reveal where value leaks — a 2× improvement at the leakiest stage can double revenue without increasing traffic spend.

Sources: Kaushik "Web Analytics 2.0" (2009) Ch. 5–6; HubSpot "State of Marketing" (annual); Forrester "Customer Journey Analytics" (2022)

## Steps

1. **Map the funnel stages** — define the stages for your business: Awareness → Interest → Consideration → Intent → Purchase → Retention → Advocacy; adapt to match your specific buyer journey.
2. **Define conversion events** — assign a measurable event to each stage transition: session → email signup (awareness→interest), email → demo request (interest→consideration), demo → trial (consideration→intent), trial → paid (intent→purchase).
3. **Collect baseline data** — pull 90-day data for each stage: volume (absolute), conversion rate to next stage, time-in-stage, and channel source breakdown; use GA4, HubSpot, Salesforce, or equivalent.
4. **Calculate stage conversion rates** — compute: stage CVR = (entries into next stage) ÷ (entries into current stage); identify which stage has the largest absolute drop-off (volume × lost conversion = lost revenue impact).
5. **Benchmark against industry standards** — compare to industry benchmarks: landing page CVR (2–5% B2B, 3–8% B2C), email open rate (20–30%), trial-to-paid (15–25% SaaS), checkout completion (65–80% e-commerce).
6. **Segment by channel and cohort** — break down conversion rates by acquisition channel, geography, device, and customer segment; identify which segments convert best and worst.
7. **Identify the biggest revenue opportunity** — calculate: Revenue Impact = Monthly Volume × Stage CVR Gap × AOV; prioritize the stage with highest revenue impact, not just lowest conversion rate.
8. **Diagnose root causes** — for each priority stage, analyze: qualitative (user research, session recordings, heatmaps, survey responses) and quantitative (bounce rate, exit pages, form abandonment, page load time).
9. **Generate and prioritize hypotheses** — write specific improvement hypotheses for each diagnosed issue: "Users abandon the pricing page because pricing is not visible above the fold — adding a pricing summary card will improve CTR to checkout by 20%."
10. **Test and measure** — implement changes via A/B test where possible; measure impact on the target stage CVR and downstream revenue; iterate monthly.

## Rules

- Fix the biggest leaking stage first — micro-optimizing a well-converting stage wastes effort; the worst stage has the highest leverage.
- Never audit without segmentation — average conversion rates hide high-performing and underperforming segments that need different strategies.
- Qualitative and quantitative data must both inform the audit — numbers identify where problems are; user research explains why.
- Revenue impact, not conversion rate, determines priority — a 1% improvement at a $10,000 AOV stage beats 10% improvement at a $50 AOV stage.
- Attribution data must be consistent across the audit — mixing last-click and data-driven attribution in the same funnel view creates false conversion pictures.

## Common Mistakes

- **Measuring only top-of-funnel volume** — driving more traffic into a leaking funnel wastes acquisition spend; fix conversion before scaling.
- **Ignoring retention and advocacy stages** — the funnel does not end at purchase; repeat purchase rate and referral rate are significant revenue multipliers.
- **Acting on too-short data windows** — 1-week funnel data has high variance; use 90-day minimum for baseline, 30-day for trend detection.
- **No segmentation** — a 3% overall checkout CVR may hide 8% for organic and 0.5% for paid social — identical treatment produces wrong interventions.
- **Hypothesizing without user evidence** — funnel data shows where to look; user sessions and interviews explain what to fix.

## When NOT to Use

- Pre-product businesses with no existing funnel data
- Campaigns running fewer than 30 days (insufficient data for reliable conversion rate calculation)
- Businesses with extremely low traffic volumes where statistical significance is unachievable
