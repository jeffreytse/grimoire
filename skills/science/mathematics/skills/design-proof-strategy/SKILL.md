---
name: design-proof-strategy
description: Use when selecting an approach for constructing a mathematical proof or verifying a conjecture
source: George Polya "How to Solve It" (1945), Daniel Solow "How to Read and Do Proofs" (2013), Paul Halmos "How to Write Mathematics" (1970)
tags: [mathematics, proof, logic, induction, contradiction, direct-proof, problem-solving]
verified: true
---

# Design Proof Strategy

Select the appropriate proof technique for a mathematical claim and build the proof structure before writing it out.

## Why This Is Best Practice

**Adopted by:** Polya's heuristics are foundational to mathematics education worldwide; used in every major undergraduate real analysis and discrete mathematics course
**Impact:** Polya's "How to Solve It" has sold over 1 million copies and is cited as the single most influential book in mathematics pedagogy; its heuristics demonstrably reduce the time to find valid proofs in educational settings

**Why best:** Choosing the wrong proof strategy wastes significant time. A direct proof that cannot be completed is often a signal that contradiction or contrapositive would work. Understanding which technique fits which claim structure prevents the common failure of pushing one approach until the proof collapses. Polya's heuristics — draw a picture, work backward, find an analogous problem — consistently surface the right strategy faster than random exploration.

## Steps

1. **Understand the claim** — Rewrite the statement formally: identify what is given (hypotheses) and what must be shown (conclusion). A proof proves P → Q; be clear on what P and Q are.
2. **Consider direct proof first** — If Q follows naturally from P by algebraic manipulation, definition application, or known theorems, attempt direct proof. Most computational proofs are direct.
3. **Consider contrapositive if direct fails** — Proof by contrapositive proves ¬Q → ¬P instead of P → Q (logically equivalent). Useful when the negation of Q provides strong constraints that make the proof tractable.
4. **Consider proof by contradiction** — Assume the claim is false (P ∧ ¬Q) and derive a contradiction. Useful for claims about uniqueness ("there is only one..."), impossibility ("there is no..."), and irrational/transcendental numbers.
5. **Consider mathematical induction for universal claims over integers** — For claims of the form "for all n ≥ n₀, P(n)": prove base case P(n₀), then prove P(n) → P(n+1). Strong induction (assuming P for all k < n) is available when the inductive step requires more than just P(n-1).
6. **Consider construction for existence proofs** — "There exists x such that P(x)" is often best proved by constructing the x explicitly. Alternatively, use probabilistic or pigeonhole arguments for non-constructive existence.
7. **Identify and apply known lemmas** — Before building from scratch, search for known theorems that imply the result. Recognize structural patterns (divisibility, group theory, topology) and apply the appropriate lemma.
8. **Write a proof outline before the proof** — State the strategy in one sentence. Write the major logical steps as a numbered list. Then fill in the details. A proof without an outline often wanders and introduces unjustified steps.

## Rules

- Never begin writing a proof before identifying the strategy — prose mathematics written without a plan produces circular or incomplete arguments.
- A failed direct proof is information — what step did it fail at? That failure reveals why an indirect approach is needed.
- Every "it is clear that" or "obviously" in a proof is a warning flag — write out what is actually clear to verify no implicit assumption is being smuggled in.
- For induction, state the inductive hypothesis explicitly before using it — "assume P(n) is true for some n ≥ n₀" must appear verbatim.

## Examples

Claim: √2 is irrational. Strategy selection: Irrationality is an impossibility claim ("there exist no integers p, q such that √2 = p/q in lowest terms") → proof by contradiction is the natural fit. Assume √2 = p/q in lowest terms → 2q² = p² → p² is even → p is even → p = 2k → 2q² = 4k² → q² is even → q is even. But p and q were assumed coprime (contradiction). QED.

## Common Mistakes

- Attempting direct proof on every claim — some claims are not directly provable and will resist indefinitely; recognizing when to switch strategies is a core skill.
- Writing the inductive step without stating the inductive hypothesis — assumes what must be proved; a logical error that invalidates the proof.
- Circular reasoning in contradiction proofs — deriving the contradiction using the original claim (not its negation) proves nothing.
- Confusing "proof by example" with a universal proof — a single example demonstrates a claim for one case; it does not prove the universal.

## When NOT to Use

- When the goal is numerical computation or approximation rather than deductive proof — selecting a proof strategy is irrelevant if the task requires a calculated result rather than a logical derivation.
- When the claim is an open conjecture at the research frontier where no known technique is sufficient — applying standard proof strategies to unsolved problems wastes time better spent on exploratory heuristics, special cases, and literature review.
- When the claim's truth value is first being established empirically via counterexample search — verifying a conjecture is false requires producing one counterexample, not designing a proof strategy for it.
