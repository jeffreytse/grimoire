---
name: design-onboarding-flow
description: Use when designing or redesigning a user onboarding flow for a SaaS or digital product where activation rate, time-to-value, or early churn are problems.
source: Intercom "The Onboarding Manifesto" (2020); Bush "Product-Led Growth" (2019); Appcues "The State of Product Onboarding" (2023); Ramli John "Delete Your Drip" methodology; Samuel Hulick "The User Onboarding Academy"
tags: [ux, onboarding, product-led-growth, activation, saas]
verified: true
---

# Design a User Onboarding Flow

Build an onboarding experience that guides new users to their first meaningful outcome as quickly as possible, establishing the habits and understanding that drive long-term retention.

## Why This Is Best Practice

**Adopted by:** Slack, Notion, Figma, Linear, Duolingo, and most product-led growth (PLG) companies — all of which treat onboarding as a growth function, not a support function.
**Impact:** Appcues (2023) reports that improving onboarding completion rate by 25% correlates with a 34% improvement in 90-day retention; Intercom's onboarding research shows that users who reach their first "aha moment" within 7 days have 3× the 6-month retention rate of those who do not; Wes Bush (2019) documents that PLG companies with systematic onboarding grow 2× faster on average than sales-led equivalents in the same category.
**Why best:** Onboarding is the moment when user intent is highest and cognitive investment in learning is most justified. A systematic flow converts that intent window into habit formation before the user's attention returns to competing priorities. Unstructured onboarding wastes the highest-value moment in the customer relationship.

Sources: Intercom — "The Onboarding Manifesto" (2020); Bush, Wes — "Product-Led Growth: How to Build a Product That Sells Itself" (2019); Appcues — "The State of Product Onboarding" (2023); Ramli John — "Delete Your Drip" (2021); Hulick, Samuel — "The User Onboarding Academy" (2019).

## Steps

1. **Define the activation milestone before designing anything** — Activation is the first moment when the user gets meaningful value from the product. For Slack: sending the first message to a teammate. For Figma: completing the first frame. For Duolingo: finishing the first lesson. Define activation specifically and measurably — not "user engaged" but "user completed their first [specific action] within [timeframe]." Every onboarding decision should be evaluated against its effect on reaching this milestone.

2. **Map the critical path from signup to activation** — List every step between account creation and the activation milestone. Identify which steps are value-creating (the user is doing something meaningful) and which are friction (the user is filling out profile fields, waiting for emails, reading documentation). The onboarding design problem is: eliminate friction, compress the critical path, make every remaining step feel valuable.

3. **Apply the Jobs-to-Be-Done lens** — Users hire your product to make progress in a specific context. The onboarding flow must speak to that job, not to your product's feature set. Instead of "Here are our 12 features," ask "What's the first job you want to accomplish?" and route users directly to the flow that delivers it. Segmentation at signup (by role, goal, or use case) allows personalized paths that reduce time-to-value by removing irrelevant steps.

4. **Design the empty state as an onboarding trigger** — The first view most users see is an empty dashboard. This is the most important onboarding screen. An empty state that says "No projects yet. Create one." is a missed opportunity. A designed empty state demonstrates what the product looks like in use, provides one clear next action, and contextualizes the value the user is about to unlock.

5. **Use progressive disclosure to reduce cognitive load** — Do not show all features during onboarding. Show only the features required to complete the first activation milestone. Reveal additional features contextually as the user demonstrates readiness through behavior (completed first project → now surface collaboration features). Front-loading features overwhelms users and delays activation.

6. **Design a product tour that teaches by doing, not by explaining** — Tooltip-based tours that describe features ("This is the dashboard") are less effective than action-based flows that guide users through doing ("Click here to create your first project"). Users learn through action; description alone does not build muscle memory or habit. If a feature requires explanation to be useful, reconsider the feature design before documenting it in onboarding.

7. **Build an email sequence that extends onboarding** — In-app onboarding covers sessions 1–2; email covers days 2–14. Design a behavioral email sequence triggered by activation state, not by time: "You created your account but haven't connected your first integration — here's a 2-minute guide" is more effective than a day-3 drip email regardless of what the user has done. Ramli John's "Delete Your Drip" methodology: replace time-based drip with behavior-triggered sequences.

8. **Identify and remove friction checkpoints** — Friction checkpoints are steps that cause users to pause, hesitate, or exit. Common examples: mandatory credit card before trial, mandatory profile completion before first use, mandatory team invitation before solo use. Audit your onboarding funnel with drop-off analytics. Any step with >15% drop-off rate is a friction checkpoint requiring investigation and redesign.

9. **Instrument the entire onboarding funnel** — Track completion rate for every step in the flow. The funnel analysis reveals where users exit (optimize those steps) and where users fly through (consider compressing those steps with others). Without instrumentation, onboarding improvements are guesses. Minimum instrumentation: step completion events, activation event, day-1/day-7/day-30 retention segmented by onboarding path taken.

10. **Run ongoing experiments on the highest-impact step** — After launch, continuously test variations on the step with the highest drop-off rate. A/B test one element at a time: copy, CTA placement, number of required fields, presence of social proof, tutorial format (video vs. interactive vs. checklist). The onboarding flow is never "done" — it is a system under continuous optimization.

## Rules

- Never require profile completion before demonstrating product value — users do not know enough about the product to make accurate profile decisions at signup; let them get value first, then ask for context.
- The onboarding checklist must be completable in one session for simple products and within 3 days for complex products; longer onboarding checklists are abandoned.
- Every onboarding email must include a direct link back to the specific incomplete step — not the dashboard home — to minimize re-orientation cost.

## Common Mistakes

- **Designing onboarding for the average user** — Onboarding designed for no one in particular serves no one well. Segment users by role or goal at signup and deliver targeted paths; a developer integrating an API and a marketing manager scheduling campaigns need entirely different first experiences.
- **Confusing product education with onboarding** — Onboarding is not a product tour. Its purpose is to get the user to their first value moment, not to explain how every feature works. Feature education happens after activation, through contextual guidance and documentation — not during the first session.
- **Over-engineering the checklist** — A 12-item onboarding checklist is not an onboarding flow, it is a homework assignment. Limit checklists to 3–5 items maximum; each item should be a single action that directly contributes to activation.
- **Ignoring mobile onboarding** — If significant traffic comes from mobile, the onboarding flow must be tested and optimized on mobile separately. Desktop onboarding designs frequently break on mobile due to tooltip positioning, form usability, and step complexity.

## Examples

**Activation-first design:** Notion's onboarding skips profile setup entirely and drops new users directly into a template gallery — the first action is duplicating a template to their workspace, which demonstrates the product's value in 60 seconds.

**Behavioral email:** Intercom's onboarding email sequence sends an "install the messenger" prompt only to users who have not completed the integration step after 48 hours — users who completed it receive a "next step" email instead, showing behavioral triggering reduces irrelevant email volume and increases completion rate.

## When NOT to Use

- Enterprise products where onboarding is primarily handled by a customer success team — the design challenge is the CS playbook, not the self-serve UI flow.
- Single-use or infrequent-use products where "habit formation" is not relevant — a tax filing tool does not need habit-forming onboarding; it needs clarity and speed.
