---
name: apply-sora-ame-kasa
description: Use when structuring a business recommendation, analysis, or proposal — to ensure every argument moves clearly from observable fact (Sky) through logical interpretation (Rain) to a specific recommended action (Umbrella), eliminating ambiguity about what decision is being requested.
source: 'Teruyuki Monnai, "ロジカル・シンキング" (Logical Thinking), Toyo Keizai (2001, 600,000+ copies); McKinsey Japan analyst training; Barbara Minto, "The Pyramid Principle", Pitman (1987) — conceptual parallel'
tags: [structured-thinking, recommendation, communication, consulting, logic, business-writing]
---

# Apply Sora-Ame-Kasa

Structure every business recommendation as three linked layers — Sky (the observable fact), Rain (the logical interpretation), Umbrella (the specific recommended action) — so the audience immediately understands what was observed, what it means, and exactly what you are asking them to decide.

## Why This Is Best Practice

**Adopted by:** Standard framework in Japanese consulting and business communication: McKinsey Japan, BCG Japan, Bain Japan, Accenture Japan, and Deloitte Tohmatsu all include it in new hire training. Taught in Japanese MBA programs at Hitotsubashi, Waseda, and Keio Business Schools. Monnai's "ロジカル・シンキング" (2001), which codifies this framework, has sold 600,000+ copies and is the single most-cited Japanese business communication text in management training curricula across major Japanese corporations (Toyota, Sony, Softbank training syllabi). The framework is the Japanese formulation of the same observation→interpretation→action structure that Barbara Minto's Pyramid Principle encodes for Western consulting.
**Impact:** McKinsey Japan data (cited in Monnai 2001): structured Sora-Ame-Kasa recommendations reduced client rework cycles by over 30% compared to unstructured presentations. In Japanese enterprise communication — where indirect communication norms create systemic ambiguity about whether a speaker is sharing information or requesting a decision — Sora-Ame-Kasa eliminates this ambiguity by making the Kasa (action request) explicit and structurally non-optional. Deloitte Tohmatsu internal training materials document that junior consultants who apply the framework consistently receive "actionable recommendation" ratings 2× higher than those who present analysis without it.
**Why best:** Most business communication fails at the Ame→Kasa step: analysts present compelling data (Sora) and insightful analysis (Ame) but leave the action implicit, forcing the audience to infer the recommendation. This produces "interesting but what do you want me to do?" reactions. The framework makes Kasa structurally mandatory — you cannot finish the communication without stating a specific action. The alternative (unstructured narrative) allows the presenter to end at analysis without ever committing to a recommendation, which is the single most common failure mode in consulting deliverables.

Sources: Teruyuki Monnai "ロジカル・シンキング" (Toyo Keizai, 2001); McKinsey Japan training; Barbara Minto "The Pyramid Principle" (Pitman, 1987)

## Steps

### 1. Identify your Sora (Sky — Observable Fact)

State one concrete, undeniable, specific fact or observation. Sora must pass the "can anyone in the room dispute this?" test.

```
✅ Sora examples (specific, undeniable):
  "Monthly active users declined from 2.1M to 1.8M over the past 90 days"
  "Our Q3 gross margin fell from 68% to 59%"
  "Customer support tickets increased 40% in the 2 weeks following the v2.3 release"
  "The competitor launched a comparable product at 30% lower price point last month"

❌ Not Sora (contains interpretation or opinion):
  "Our product is losing market position" → contains interpretation (Ame)
  "The team is not performing" → vague, not undeniable
  "We have a pricing problem" → already an interpretation
```

If you cannot state Sora without using the word "because" or "so" — stop. You have skipped to Ame.

### 2. Derive your Ame (Rain — Logical Interpretation)

Answer: "So what does this fact mean?" Ame is your analytical judgment — the implication that follows from the Sora observation.

```
✅ Ame examples (logical inference from Sora):
  Sora: "MAU dropped from 2.1M to 1.8M over 90 days"
  Ame: "The decline is concentrated in users acquired via the paid channel (cohort analysis
        shows organic users are stable), suggesting our paid acquisition targeting has
        degraded — not product quality"

  Sora: "Q3 gross margin fell from 68% to 59%"
  Ame: "COGS increase is driven entirely by cloud infrastructure costs (+$2.1M), which
        scaled faster than revenue — suggesting our pricing model does not account for
        per-user compute costs above threshold usage"
```

Ame quality check:
- Does Ame follow logically from Sora, or does it require additional unstated assumptions?
- Is Ame falsifiable? (If there's no evidence that could disprove it, it's an opinion, not an interpretation)
- Is there only one Ame? If multiple interpretations are equally plausible, state the most defensible one and acknowledge alternatives in your backup slides

### 3. State your Kasa (Umbrella — Recommended Action)

State exactly what you are recommending. Kasa must be specific, time-bound, and owned.

```
✅ Kasa examples (specific, actionable, owned):
  "Pause paid acquisition spend by October 15 pending audit of targeting parameters;
   redirect $200K/month to retention campaigns targeting the at-risk cohort"

  "Cap infrastructure auto-scaling at 150% of current baseline and implement usage-based
   pricing tier for accounts above 500 API calls/day — ship pricing page update by Nov 1"

❌ Not Kasa:
  "We should consider improving our paid acquisition" → not specific
  "The team should look into this" → not owned, not time-bound
  "We need to address the margin issue" → describes the problem, not an action
```

Kasa test: can the decision-maker approve or reject this recommendation right now, with the information in the room? If they need to ask "what specifically?" — Kasa is incomplete.

### 4. Chain multiple Sora-Ame-Kasa units for complex arguments

For multi-part recommendations, each point in the argument is its own Sora-Ame-Kasa unit. The overall document has one top-level Kasa (the main ask), with each supporting section providing the Sora-Ame chain that justifies it.

```
Top-level Kasa: "Approve $1.2M budget reallocation from paid acquisition to product-led growth"

Supporting unit 1:
  Sora: "Paid CAC increased 3× in 18 months (from $42 to $127)"
  Ame:  "ROI on paid channel is now negative at current LTV"
  Kasa: "Eliminate paid channel spend by Dec 1"

Supporting unit 2:
  Sora: "NPS referral source = 34% of new signups, with zero marketing spend"
  Ame:  "Organic word-of-mouth is our most efficient acquisition channel"
  Kasa: "Invest $800K in referral program and in-product virality features"
```

### 5. Apply to written communication (emails, slide decks, memos)

```
Email structure:
  Subject: [Action required] Pause paid acquisition — recommendation for your decision
  
  [Sora] Our paid channel CAC has increased from $42 to $127 over 18 months,
  while organic CAC has remained stable at $18.
  
  [Ame] The paid channel is now negative-ROI at our current LTV of $95,
  and organic channels are underinvested relative to their efficiency.
  
  [Kasa] I am recommending we pause paid acquisition by October 15 and
  redirect the $200K/month budget to referral program development.
  
  Please confirm approval by October 10 so we can notify agency partners.

Slide structure:
  Title slide: The specific Kasa (recommendation) — not "Q3 Analysis"
  Slide 1: The Sora (data/observation)
  Slide 2: The Ame (what it means)
  Slide 3: The Kasa restated with implementation plan
```

## Rules

- Sora must be undeniable — if the audience can dispute the fact in the opening sentence, the entire chain loses credibility before it begins; verify Sora data before the meeting.
- Ame is the only place for judgment — Sora is facts only; Kasa is action only; analytical reasoning belongs exclusively in Ame.
- Kasa must be specific enough to approve or reject — if it requires follow-up clarification ("what specifically?"), it is not Kasa yet.
- One primary Kasa per communication — if you have three recommendations of equal weight, you do not have a recommendation; prioritize and lead with the most important Kasa.

## Common Mistakes

- **Starting with Ame** — "Our product is losing market position" sounds like Sora but is already an interpretation; readers will dispute it immediately. Always open with the undeniable observable fact.
- **Ending at Ame** — presenting analysis and stopping ("as you can see, the data suggests challenges") without stating a Kasa produces "interesting presentation" feedback, not decisions. Kasa is not optional.
- **Weak Kasa: "we should explore"** — "explore," "consider," "look into," and "evaluate" are not Kasa; they are deflections. Kasa names a specific decision with an owner and a date.
- **Skipping the Ame link** — jumping from Sora directly to Kasa ("MAU fell 15%, therefore increase marketing spend") makes the recommendation seem arbitrary; the Ame (interpretation) is the logic that justifies the action, and omitting it causes audiences to reject the Kasa even when it is correct.
