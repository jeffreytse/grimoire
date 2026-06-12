---
name: design-spaced-repetition
description: Use when designing a curriculum, flashcard deck, or study system intended to maximize long-term retention of facts, concepts, or procedures
source: Hermann Ebbinghaus "Über das Gedächtnis" (1885); Piotr Wozniak SM-2 algorithm (1987); Gwern Branwen spaced repetition review
tags: [spaced-repetition, memory, learning-science, curriculum-design, retention, anki]
verified: true
---

# Design Spaced Repetition

Build a review schedule and card system that exploits the spacing effect to maximize long-term retention with minimum study time.

## Why This Is Best Practice

**Adopted by:** Medical schools (Step 1 boards prep), language learning platforms (Duolingo, Anki), military training curricula, and over 5M active Anki users
**Impact:** Ebbinghaus showed retrieval at spaced intervals reduces the forgetting curve from 50% recall at 1 day to 80%+ recall at 30 days; Cepeda et al. (2008) meta-analysis of 254 studies confirmed optimal spacing outperforms massed practice by 10–30% on delayed tests

**Why best:** Massed practice (cramming) produces short-term recall. Spaced repetition exploits the spacing effect and the testing effect to encode material into long-term memory by retrieving it just as forgetting begins.

## Steps

1. **Define the retention target** — Decide: what must be recalled, at what accuracy, and after how long? "Recall 200 medical terms at 90% accuracy after 6 months" drives the spacing schedule differently than "recall for next week's exam."
2. **Write minimum-information cards** — One card = one atomic fact. Bad: "Explain photosynthesis." Good: "What molecule accepts electrons in the light-dependent reactions? → NADP+." Atomic cards produce precise retrieval; complex cards produce guessing.
3. **Use cloze deletion for context** — For facts that need context: "The [enzyme] that joins Okazaki fragments is DNA ligase." Cloze preserves context while testing one specific item. Avoid clozing multiple items per card.
4. **Add imagery and mnemonic links** — Pair abstract facts with a vivid image or story. The Method of Loci (memory palace) combined with spaced repetition is the most effective retention technique in cognitive psychology.
5. **Implement an SM-2-compatible algorithm** — Use Anki (free, open-source SM-2 derivative) or SuperMemo. Set new cards per day to match acquisition rate: language learners typically 10–20 new words/day; medical students 50–100 cards/day during intensive study.
6. **Schedule based on performance** — SM-2 assigns the next review interval based on quality of recall (0–5 scale). Cards rated ≤2 are re-queued immediately; cards rated 4–5 expand their interval (1d → 4d → 10d → 25d…).
7. **Design deck hierarchy** — Organize cards into decks by subject and sub-topic. Tag by difficulty and source. Never mix fundamentals and advanced content in the same review session — prioritize foundational cards first.
8. **Review daily without skipping** — The system breaks down when reviews accumulate. Even 10 minutes daily maintains the schedule; skipping 3 days doubles the queue and breaks the interval calculations.

## Rules

- One atomic fact per card — if you need two clicks to recall an item, split the card.
- Cards must test retrieval, not recognition — avoid multiple choice; use fill-in-the-blank or production prompts.
- Do not add cards for material not yet understood — rote memory of misunderstood material is useless and hard to correct.
- Review during active recall, not passive re-reading — the testing effect requires genuine retrieval effort.
- Retire or suspend cards for information that is no longer needed; an overloaded deck reduces motivation.

## Examples

A medical student preparing for USMLE Step 1 uses Anki with the Zanki deck (30,000 cards). They add 100 new cards/day for 6 months, reviewed daily for 2–3 hours. On exam day, 95% of the deck has been seen ≥3 times. Average score among this methodology's users is 240+ vs. 230 national average.

## Common Mistakes

- **Complex multi-part cards** — "Describe the coagulation cascade" produces a guess and a re-queue; the specific broken-down steps produce lasting memory.
- **Adding without reviewing** — Creating 500 cards and not reviewing them for two weeks produces an unmountable backlog and abandoned decks.
- **Passive re-reading instead of retrieval** — Reading the answer before guessing eliminates the testing effect; the card must be answered before flipping.

## When NOT to Use

- When the learning goal is conceptual understanding or problem-solving skill rather than discrete fact recall — spaced repetition optimizes for retrieval of atomic items, not for developing the flexible reasoning needed to apply concepts to novel problems.
- When the material will be needed only within the next 48–72 hours and long-term retention is not required — massed practice (focused review) is more time-efficient for short-horizon recall than building and scheduling a card deck.
- When the learner does not yet understand the material being encoded — adding cards for content not yet comprehended produces rote strings that are recalled without meaning and are difficult to correct once incorrectly encoded.
