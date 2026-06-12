---
name: run-usability-test
description: Use when validating a UI design, prototype, or live product with real users to find usability problems before or after launch
source: Jakob Nielsen "Usability Engineering" (1993); Nielsen Norman Group Think-Aloud Protocol; Nielsen's 5-user heuristic
tags: [usability-testing, ux-research, think-aloud, user-research, design-validation, ux]
verified: true
---

# Run Usability Test

Conduct a structured usability test that reveals real user problems in a prototype or live product with as few as 5 participants.

## Why This Is Best Practice

**Adopted by:** Nielsen Norman Group, Google UX, IDEO, Apple HIG team — usability testing is the core research method in every major UX practice
**Impact:** Nielsen's landmark research shows 5 users discover ~85% of usability problems; each additional 5-user round yields diminishing returns, making small iterative rounds more cost-effective than large studies

**Why best:** Usability testing replaces opinion with observation. Watching a real user fail to complete a task is more persuasive to stakeholders than any heuristic review or UX opinion.

## Steps

1. **Define research questions** — Write 3–5 specific questions you want answered: "Can users find the cancel subscription option without help?" not "Is the UI good?" Questions drive task design.
2. **Write task scenarios** — Create 4–8 realistic tasks that reflect actual user goals. Phrase tasks as goals, not instructions: "You want to invite a teammate to your workspace. Please do that." Never name UI elements in the task.
3. **Recruit representative participants** — Recruit 5 users per distinct audience segment. Screen for real usage context, not demographics. Use a screener survey.
4. **Set up the test environment** — Use screen recording + audio. For remote tests: Lookback, UserTesting, or Zoom with screen share. For in-person: a quiet room, observer in the back, facilitator at the keyboard.
5. **Explain the think-aloud protocol** — Brief the participant: "We're testing the product, not you. There are no wrong answers. Please say everything you're thinking as you go." Do not say "testing your skills."
6. **Facilitate without leading** — When a user gets stuck, ask "What are you looking for?" not "Did you try the menu?" Silence is a research tool — let users struggle enough to reveal the real problem.
7. **Observe and note** — Note every hesitation, misread label, wrong path, and verbal expression of confusion. Timestamp against the recording for review.
8. **Synthesize findings into a severity matrix** — After all sessions: list every problem, rate severity (1–4: cosmetic to critical), and count how many participants hit it. Fix in severity order.

## Rules

- Tasks must be written as user goals, never as instructions that describe the UI.
- Facilitators do not help users complete tasks — that eliminates the data.
- Run a pilot session first with a colleague to debug the task script and tech setup.
- Synthesize findings within 48 hours of the last session while observations are fresh.
- Share a video clip of the most critical finding — a two-minute clip changes stakeholder minds faster than a 20-page report.

## Examples

A checkout flow test with 5 users revealed that 4 of 5 participants assumed "Continue" meant "save and exit" rather than "proceed to payment" — a single word change in the button label resolved a 12% cart abandonment rate in the subsequent A/B test.

## Common Mistakes

- **Tasks that describe the UI** — "Click the Settings icon in the top right" removes all the discovery behavior you were trying to observe.
- **Facilitator over-helping** — Jumping in when a user hesitates gives you false success data; real users won't have a facilitator.
- **Testing too late** — Running usability tests on shipped code instead of prototypes means findings require expensive rework; test early and often.

## When NOT to Use

- Do not run a moderated usability test when the question is quantitative — if you need to know what percentage of users succeed or which of two designs performs better, an unmoderated study or A/B test provides statistically valid data that think-aloud sessions cannot.
- Do not use usability testing to validate business strategy or product-market fit — if users can complete tasks easily but do not want the product, that is a market research problem, not a usability problem.
- Do not recruit convenience participants (teammates, friends, or people in the office) when the target audience has specialized domain knowledge or context — testing a medical records interface with non-clinicians produces misleading results about real-world task failure.
