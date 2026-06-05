---
name: design-email-marketing-campaign
description: Use when planning or building an email marketing campaign including segmentation, sequencing, copy strategy, and performance measurement
source: Mailchimp Email Marketing Benchmarks (annual); DMA (Data & Marketing Association) Email Marketing Guidelines; Litmus "State of Email" report
tags: [email, marketing, campaign, automation, copywriting]
verified: true
---

# Design Email Marketing Campaign

Plan and execute an email campaign with the right audience segmentation, message sequencing, copy framework, and measurement to maximize open rates, clicks, and conversions.

## Why This Is Best Practice

**Adopted by:** HubSpot, Mailchimp, Klaviyo — all publish email best practices grounded in their platform data; DMA guidelines reflect industry legal and ethical standards
**Impact:** Email delivers average ROI of $36–42 per $1 spent (Litmus/DMA 2022); segmented campaigns achieve 14% higher open rates and 100% higher CTR than non-segmented (Mailchimp benchmark data)
**Why best:** Litmus research shows that 41% of email opens occur on mobile; DMA segmentation data consistently shows personalization and relevance as the top drivers of engagement — both require a systematic campaign design process.

Sources: Mailchimp "Email Marketing Benchmarks" (2023); DMA "Email Marketing Industry Census" (2022); Litmus "State of Email" (2023)

## Steps

1. **Define campaign objective** — select one primary goal: promotional (revenue), nurture (engagement), onboarding (activation), re-engagement (win-back), transactional (service); objective determines content and sequence.
2. **Define and segment the audience** — segment by: behavior (purchased/not purchased, active/lapsed), demographics, funnel stage (lead/MQL/SQL/customer), preference data; more specific segments outperform broad lists.
3. **Plan the sequence** — determine number of emails, send intervals, and branching logic; onboarding: 5–7 emails over 14 days; promotional: 3–5 emails over 7–10 days; re-engagement: 3 emails over 21 days.
4. **Write subject lines and preview text** — subject line is 40% of open rate decision; write 5 variations and pick the strongest; use: curiosity gaps, specificity (numbers), personalization (first name), urgency (when genuine); preview text extends the subject line.
5. **Write email body copy** — use AIDA structure (Attention, Interest, Desire, Action) or Problem-Agitate-Solution; keep body copy to 200–400 words for promotional emails; one primary CTA per email.
6. **Design for mobile-first** — 41% of opens on mobile; use single-column layout, 14–16px body text, touch-friendly CTAs (min 44px height), compressed images (<100KB for fast load).
7. **Configure technical setup** — verify: SPF, DKIM, DMARC authentication (deliverability); list unsubscribe header (legal); working unsubscribe link (CAN-SPAM/GDPR); suppress existing unsubscribes and bounces.
8. **A/B test subject lines** — test subject line A vs. B on 20% of list each; send winner to remaining 60% after 2–4h; this consistently improves open rates 5–15%.
9. **Set send time** — analyze audience engagement data; B2B: Tuesday–Thursday 9–11am local time; B2C: varies by product — test Thursday evening and Saturday morning; avoid Monday AM and Friday PM.
10. **Measure and optimize** — track: open rate (20–30% healthy B2B), CTR (2–5% healthy), conversion rate (1–3% for promotional), unsubscribe rate (<0.5%), revenue per email; optimize underperforming sequences.

## Rules

- Every email list requires explicit opt-in consent (GDPR, CAN-SPAM, CASL) — purchased or scraped lists are illegal in most jurisdictions and destroy deliverability.
- One primary CTA per email — multiple competing CTAs reduce click rates; pick the most important action.
- Every email must have a visible, one-click unsubscribe — CAN-SPAM requires 10-day processing; GDPR requires immediate; honor immediately regardless.
- Plain-text version required — some clients and filters reject HTML-only emails; always include a plain-text alternative.
- List hygiene is mandatory — hard bounces must be removed immediately; suppress soft bounces after 3 attempts; remove inactive subscribers older than 12 months.

## Common Mistakes

- **One email to the entire list** — mass unsegmented emails produce average results; even basic segmentation by engagement tier (active/lapsed/inactive) dramatically improves performance.
- **No clear CTA** — emails ending with "let us know if you have questions" produce no measurable action; every email needs a specific requested action.
- **Images without alt text** — many clients block images by default; an email that is all images with no alt text is blank to a significant portion of recipients.
- **Sending at arbitrarily chosen times** — "we always send Tuesday at 10am" without testing against audience data ignores significant send-time variation by industry and persona.
- **No suppression of recent purchasers** — sending a discount promotional email to customers who paid full price 48h ago creates frustration and damages trust.

## When NOT to Use

- No opted-in list (do not build list via purchased contacts)
- Audiences where email engagement is structurally low (Gen Z preference for DM/SMS channels)
- Transactional confirmations — use dedicated transactional email systems, not marketing automation
