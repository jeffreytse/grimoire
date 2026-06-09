# The Grimoire Philosophy

*The world's knowledge is in your AI. The world's practice is not.*

AI assistants have ingested every textbook, every paper, every article ever written. They understand fields. They do not practice them. Practice is what happens after 10,000 hours. Practice is what a senior surgeon does without thinking. Practice is what a staff engineer knows not to do. Practice is what grimoire encodes.

**In short:** Grimoire encodes the world's best practitioners' actual methods — proven at scale, cited, and actionable. Every skill is a verb you can use right now. Built by anyone who has mastered their craft. Free, open source, for every field, forever.

---

## The world's best practices belong to everyone.

A McKinsey engagement costs $1M. A senior lawyer bills $800/hr. A structural engineer isn't available at 2am. The practices they follow — proven at the highest levels — are not proprietary. They belong to the world. Grimoire makes them free.

## Skills are verbs.

Not descriptions of what experts know — the exact steps they take, in the exact situation they face, proven at scale. If you can't act on it in the next five minutes, it isn't a skill.

> **Story:** A junior engineer is assigned their first code review. They Google "how to review code." They get a Medium post: *be constructive, focus on the big picture, be kind.* They read it twice. They agonize for 45 minutes, leave vague comments like "this could be cleaner," and ship nothing actionable. The author doesn't know what to fix. The reviewer feels like they failed.
>
> Next time, they open the `review-pull-request` skill. Step 1: check security boundaries — 5 minutes. Step 2: look for N+1 queries — 5 minutes. Step 3: flag naming inconsistencies — 5 minutes. Done in 20. Three specific findings, each with a line number and a fix.
>
> That's the difference between knowledge and practice.

## Skills are atoms.

A composite practice is always derivable from its atomic components. SOLID principles
and KISS are both correct — their combination is not a third skill, it is what the LLM
does at runtime when both apply.

LLMs compose at runtime. Given a situation, the AI reads the context, picks the relevant
atoms from the library, and applies them in combination. Pre-composing two skills into
one rigid unit reduces this flexibility: the composite can only be applied as a whole,
cannot substitute individual atoms, and cannot recombine with other skills.

Atoms are a library. Composites are a script.

> **Story:** A contributor submits `implement-nutrition-and-training-plan`. It covers
> macros, periodization, progressive overload, and meal timing. Thorough. Two weeks of
> work.
>
> The reviewer points out: `calculate-macros` already exists. `design-training-program`
> already exists. `apply-progressive-overload` is in progress.
>
> The composite isn't wrong — it's redundant. An athlete bulking for powerlifting needs
> different atoms than a marathon runner cutting weight. The composite applies as a
> preset; the atoms compose to the exact situation.
>
> The contributor deletes the composite and files the two missing edge cases as updates
> to the existing skills. The library is stronger.
>
> That is what atomicity means in practice.

One skill, one concept. Describe any situation — the LLM picks the right atoms.

## Every claim must be proven.

One skill. One concept. Adopted by most top-tier institutions in the field, with measurable impact and a named source. If you can't prove it, you can't ship it.

> **Story:** A contributor submits a skill: *always wrap database calls in retry logic with exponential backoff.* It sounds right. Every experienced engineer nods.
>
> The reviewer asks: which institutions mandate this? With measurable outcomes? The contributor cites their own startup. Rejected — one team's convention isn't evidence.
>
> Three months later, they return. AWS SDK implements it by default. The Google SRE workbook dedicates a chapter to it. Netflix published failure-rate data showing a 40% reduction in cascading timeouts after Hystrix adopted it. The skill ships.
>
> The wait wasn't bureaucracy. It was the quality filter doing its job.

## Consensus is the floor.

If the world's best professionals are split, grimoire acknowledges the debate — and encodes the majority position. When there is no consensus, there is no best practice to ship.

> **Story:** A contributor submits a skill: *always use a monorepo.* Google does it. Reasonable.
>
> The reviewer checks the other side. Amazon runs thousands of microrepos. Meta started monorepo, partially reversed. The academic literature is inconclusive. The contributor wants grimoire to pick a side. Grimoire can't — the world's best practitioners are genuinely split.
>
> The skill gets held. Not rejected — held, pending evidence.
>
> Meanwhile, a different skill ships: `migrate-to-monorepo` — the exact steps for teams who have *already decided* to switch. A verb, not a verdict. The debate lives in the discussion thread. The practice waits for the evidence.

## Anyone who has mastered their craft can contribute.

A nurse. A jazz musician. A securities lawyer. A structural engineer. Grimoire is not a developer project. It is a project for everyone who has spent 10,000 hours in a field and has something to say about how it's really done.

> **Story:** A travel ICU nurse — 14 years, six hospitals, three countries — has handed off critically ill patients at shift change hundreds of times. She knows the exact sequence: vitals trend first, then pending labs, then family status, then the one thing most likely to go wrong before morning. She has watched nurses skip steps and patients deteriorate. She has the practice cold.
>
> No developer has this knowledge. No paper captures her exact sequence. The AI knows the theory of handoffs. It does not know what she knows.
>
> She opens a PR. Her first GitHub commit. The skill she writes — `conduct-icu-shift-handoff` — is more useful to 50,000 ICU nurses than anything in the engineering subdomain.
>
> Grimoire is for her.

## The skill outlasts the AI.

Plain Markdown. No lock-in. No proprietary format. These skills will outlive every AI assistant currently running.

---

## The Intellectual Foundations of Best Practices

The philosophy of best practices is a synthesis of three intellectual traditions — pragmatism, empiricism, and systems thinking — that together explain why proven methods outperform both pure theory and individual intuition.

### Pragmatism

The excellence of an idea or method is not judged by the grandeur of its theory, but by its ability to produce results in real-world action. Effectiveness is truth.

Rooted in the pragmatism of C.S. Peirce and John Dewey, this tradition holds that the test of any idea is its practical consequence. Grimoire encodes this directly: a skill that cannot produce a concrete outcome in the next five minutes is not a skill.

### Empiricism and Inductive Logic

Best practices are not invented — they are distilled. Through extensive practical activity, observation, and inductive reasoning, communities of practice converge on what consistently works. The mechanism is a closed loop:

**experience → refinement → standardization → verification → iteration**

This is why grimoire requires named sources and measurable outcomes. A practice that cannot be traced back to observed results is not a best practice — it is a hypothesis.

### Continuous Improvement and Systems Thinking

There is no absolute "best" — only "better" in a specific context at a specific time. Best practices are dynamically evolving. The Kaizen tradition and PDCA cycle (Plan–Do–Check–Act) formalize this: every standard is the current best known approach, not the final answer.

Systems thinking complements this by breaking complex problems into standardized, testable processes — reducing reliance on individual memory and judgment, and making quality reproducible at scale.

Together, these three traditions explain why grimoire works: it encodes only what has been proven (empiricism), in actionable form (pragmatism), while remaining open to revision as evidence accumulates (continuous improvement).
