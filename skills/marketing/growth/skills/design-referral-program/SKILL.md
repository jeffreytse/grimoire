---
name: design-referral-program
description: Use when designing a customer referral or word-of-mouth program for a product or service
source: Sean Ellis & Morgan Brown "Hacking Growth" (2017), Dropbox and Airbnb referral case studies, Viral Loops methodology
tags: [growth, referral, viral-coefficient, word-of-mouth, acquisition, retention]
verified: true
---

# Design Referral Program

Design a referral program that systematically converts satisfied customers into an acquisition channel.

## Why This Is Best Practice

**Adopted by:** Dropbox (60% growth from referrals), Airbnb (double-sided referral drove significant expansion), PayPal (early $10/$10 referral bonus)
**Impact:** Dropbox's two-sided referral program increased signups by 60% and produced a 3900% growth rate over 15 months

**Why best:** Referral programs work because referred customers have pre-existing trust — they convert at higher rates and retain longer than paid acquisition. The viral coefficient (K = invites sent × conversion rate) determines whether a program achieves compounding growth or merely supplements it. Programs fail most often from misaligned incentives or friction in the sharing flow.

## Steps

1. **Identify the referral moment** — Map the customer journey; find the moment of peak satisfaction (post-first-value, post-upgrade, post-milestone) — this is when to prompt referral.
2. **Define the incentive structure** — Choose single-sided (referrer only) or double-sided (referrer + referee). Double-sided outperforms by 25%+ (Viral Loops data). Match incentive to product value: cash, credits, feature unlocks.
3. **Calculate target viral coefficient** — K = (average invites per user) × (invite-to-signup conversion rate). K > 1 = viral growth; K = 0.3-0.9 = meaningful supplement.
4. **Design the sharing mechanic** — Unique referral link is the minimum; add social sharing buttons, email templates, and in-app "invite" prompts. Reduce friction to one click.
5. **Build the reward fulfillment flow** — Automate reward delivery on referral conversion; delays kill program trust. Define fraud prevention rules (no self-referral, device fingerprinting).
6. **Set program rules and caps** — Define eligible referrals, reward caps, and expiry terms. Communicate clearly to prevent support tickets.
7. **Instrument tracking** — Track: invites sent, invite open rate, conversion rate per channel, referral CAC vs. paid CAC, and referral LTV vs. organic LTV.
8. **Launch and iterate** — A/B test incentive amounts, CTA copy, and placement. Airbnb iterated 50+ variations before finding peak performance.

## Rules

- Never launch a referral program before achieving product-market fit — unhappy users won't refer and will amplify negative word of mouth.
- Make the referral flow accessible in 2 clicks from the product dashboard.
- Referral incentives must be meaningful relative to product price — a $2 credit on a $100/mo product doesn't move behavior.
- Always measure referral LTV separately — referred customers often have higher LTV and justify higher incentive spend.

## Examples

SaaS tool at $29/mo: Referrer gets 1 month free per successful referral; referee gets 20% off first 3 months. Referral prompt appears on the "project published" success screen. Unique link auto-generated; shareable via email, Twitter, LinkedIn. Program tracked in Stripe + Mixpanel. CAC via referral: $18 vs. $94 via paid ads.

## Common Mistakes

- Single-sided incentive only — referee has no reason to act; double-sided consistently outperforms.
- No fraud prevention — self-referral and fake account farming drain program budget.
- Burying the referral prompt in account settings — placement in the peak satisfaction moment is the primary driver of participation rate.
- Launching and forgetting — referral programs decay; refresh incentives and placement quarterly.

## When NOT to Use

- Do not design a referral program before achieving product-market fit — users who are not yet delighted will not refer, and the social proof damage from neutral or negative referrals outweighs any acquisition benefit.
- Do not use a referral program as a primary acquisition channel in regulated industries (e.g., financial services, healthcare) where incentivized referrals may constitute unlawful inducement under applicable compliance rules.
- Do not launch a referral program when the product has a natural privacy expectation that makes sharing awkward — users of confidential HR, legal, or therapy tools will not share their use publicly regardless of incentive size.
