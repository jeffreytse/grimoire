---
name: apply-feynman-technique
description: Use when learning a new concept, diagnosing why understanding feels shallow, or preparing to explain complex material — to verify genuine comprehension by explaining the topic in plain language, identify exactly where the explanation breaks down, and fill those gaps until the explanation is clear and jargon-free.
source: Feynman (Gleick "Genius", Norton 1992); Oakley "A Mind for Numbers" (Tarcher/Penguin, 2014); Chi et al. "Eliciting Self-Explanations Improves Understanding" (Cognitive Science, 1994); Dunlosky et al. "Improving Students' Learning With Effective Learning Techniques" (Psychological Science in the Public Interest, 2013); Nestojko et al. "Expecting to Teach Enhances Learning" (Memory & Cognition, 2014)
tags: [learning, comprehension-verification, self-explanation, gap-identification, teaching, study-technique, knowledge-depth]
related: [apply-spaced-repetition, apply-retrieval-practice, apply-socratic-questioning]
---

# Apply Feynman Technique

Expose the exact boundaries of your understanding by attempting to explain a concept in plain language — gaps surface as vague language, circular definitions, and the impulse to reach for the textbook — then eliminate each gap until the explanation is smooth.

## Why This Is Best Practice

**Adopted by:** Barbara Oakley's "Learning How to Learn" (Coursera) — the most-enrolled MOOC in history with 4 million+ learners across 200+ countries — makes the Feynman Technique a core module. Khan Academy's pedagogical philosophy (Sal Khan) is grounded in the same principle: you understand something when you can explain it simply. Cal Newport's "How to Become a Straight-A Student" (2006) codifies it as the primary deep-understanding technique used by top university students. Used in problem-based learning (PBL) curricula at Harvard Medical School, McMaster, and Maastricht, where self-explanation is a required component of case-based rounds.

**Impact:** Chi et al. (1994) showed students who self-explained while studying scored 40% better on subsequent problem-solving tests than students who re-read the same material (n=40, physics). Dunlosky et al. (2013) systematic review of 10 popular study techniques rated elaborative interrogation — the cognitive mechanism behind the Feynman Technique — as "high utility", one of only 2 out of 10 techniques to receive that rating; re-reading (the most common technique) received "low utility". Nestojko et al. (2014) found students who expected to teach material scored 28% higher on recall and 19% higher on transfer tests than students who studied for a personal test.

**Why best:** Re-reading and highlighting produce familiarity, not comprehension — learners mistake recognition for understanding. Summarising preserves the original phrasing and hides gaps behind borrowed terminology. The Feynman Technique forces *generation* of a novel explanation in plain language, which surfaces exactly where understanding is missing. Unlike pure retrieval practice (flashcards, practice tests), it targets conceptual depth and the ability to reason with ideas rather than recall facts — making it the superior technique for complex, interrelated, or heavily abstract concepts.

Sources: Chi et al. (Cognitive Science, 1994); Dunlosky et al. (Psychological Science in the Public Interest, 2013); Nestojko et al. (Memory & Cognition, 2014); Oakley (2014)

## Steps

### 1. Write the concept name at the top of a blank page

Name only — no definition yet. This commits you to explaining one concept. If the name requires a qualifier to be specific ("Nash equilibrium in repeated games"), include it.

### 2. Explain the concept as if teaching a curious 12-year-old

Write a complete explanation without:
- Jargon or technical terms (if a term is unavoidable, define it before using it)
- Circular definitions ("recursion is when something is recursive")
- Vague quantifiers ("some", "often", "many cases")

Use:
- Concrete analogies ("a Nash equilibrium is like two drivers choosing lanes — neither wants to switch once they've committed, because switching makes things worse")
- Cause-and-effect chains
- Specific examples with actual numbers or names

Write continuously — do not stop to check notes or source material.

### 3. Identify exactly where the explanation breaks down

Re-read what you wrote. Mark every point where you:

| Symptom | What it reveals |
|---------|----------------|
| Reached for jargon without defining it | Borrowed vocabulary, not owned understanding |
| Wrote "it depends" without specifying on what | Incomplete causal model |
| Gave a circular definition | No underlying mental model exists |
| Could not produce a concrete example | Abstract knowledge not grounded |
| Explanation contradicts itself | Two conflicting partial models |
| Felt the urge to "check the book" | Gap in knowledge, not retrieval failure |

**These are your gaps.** List each one explicitly: "I cannot explain why X leads to Y without using [jargon]."

### 4. Return to source material — targeted, not broad

For each gap identified in Step 3:
- Read *only* the sections addressing that specific gap
- Do not re-read chapters you already understood — that is re-reading, not gap-filling
- Take notes in your own words as you read, not verbatim quotes

### 5. Re-explain, replacing jargon with analogies

Rewrite the explanation from Step 2, this time:
- Replace every technical term with a plain-language analogy or description
- If you cannot replace a term, the understanding behind it is incomplete — go back to Step 4
- Stress-test analogies: where does the analogy break down? Knowing its limits is part of understanding

A working analogy does not have to be perfect — it has to be accurate within the scope you need.

### 6. Repeat until smooth

The explanation is complete when:
- A curious non-expert could follow it without prior background
- You can answer "why?" at each step without returning to notes
- You can generate a new example on demand, not just repeat the one you memorised
- You can explain where your analogy breaks down (showing you understand its limits)

Typical cycles: 2–3 passes for a single concept; 4–6 for dense interrelated topics.

## Rules

- Never explain by restating the definition — paraphrase is not understanding
- If you cannot generate an example, the concept is not learned yet
- Jargon is a placeholder, not an explanation — always chase it down to plain language
- One concept per session; if explanation reveals sub-concepts you also cannot explain, queue them as separate sessions
- The goal is not to produce a simple explanation — it is to discover what you don't know

## Examples

**1. Software engineer learning Byzantine fault tolerance**
Step 2 draft: "It's when some nodes fail in a malicious or arbitrary way..." — reaches for "Byzantine" to define Byzantine. Gap identified. Step 4: read the two-generals problem analogy. Step 5 rewrite: "Imagine 5 generals voting on whether to attack. Up to 1 general might be a traitor who sends conflicting votes to different generals. Byzantine fault tolerance means the system still reaches the right decision even with 1 liar — it needs at least 3×(liars) + 1 = 4 honest generals to guarantee it."

**2. MBA student learning option pricing**
Step 2 draft: "The Black-Scholes formula prices an option based on volatility and time..." — cannot explain *why* volatility increases option value. Gap identified. Step 4: re-read the concept of uncertainty creating upside without downside (options have asymmetric payoffs). Step 5 rewrite: "An option is like insurance. Higher volatility = more chances the stock shoots up past your strike price. You capture the upside but your loss is capped at the premium. Higher chaos = more valuable insurance."

**3. Student learning Newton's third law**
Step 2 draft: "For every action there is an equal and opposite reaction" — verbatim textbook. Gap: cannot explain why a wall doesn't move when you push it. Step 4: read about force pairs and mass. Step 5: "When you push a wall, the wall pushes back equally. But the wall is connected to the Earth (billions of kg), so the acceleration is unmeasurably small (F=ma). You feel the reaction; you just can't see the wall move."
