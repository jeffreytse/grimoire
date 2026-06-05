---
name: run-heuristic-evaluation
description: Use when conducting a usability review of an interface to identify usability problems without running user tests, using Nielsen's 10 heuristics as the evaluation framework.
source: Nielsen "Usability Engineering" (1993); Nielsen & Molich "Heuristic evaluation of user interfaces" (1990); Nielsen Norman Group heuristic evaluation methodology; ISO 9241-110 (2020)
tags: [ux, usability, heuristic-evaluation, inspection, design-review]
verified: true
---

# Run a Heuristic Evaluation

Apply Nielsen's 10 usability heuristics as an expert inspection method to systematically identify usability problems in an interface before or without user testing.

## Why This Is Best Practice

**Adopted by:** Google, Microsoft, Apple, IBM, and virtually every UX team conducting design review, from enterprise software to consumer apps. Taught in every accredited HCI program globally.
**Impact:** Nielsen & Molich (1990) demonstrated that 3–5 evaluators using heuristics find 75% of major usability problems; Nielsen (1993) showed the cost-to-benefit ratio of heuristic evaluation is 48:1 — for every dollar spent, $48 of downstream usability cost is avoided. NNG research shows heuristic evaluation finds approximately 35% more problems per hour than informal design review.
**Why best:** Heuristic evaluation provides a structured vocabulary for usability problems that makes findings communicable, prioritizable, and actionable. Without heuristics, design reviews devolve into aesthetic preference debates. With them, teams can objectively categorize and severity-rate problems, enabling efficient remediation prioritization.

Sources: Nielsen, Jakob — "Usability Engineering" (1993); Nielsen, J. & Molich, R. — "Heuristic evaluation of user interfaces," CHI '90 Proceedings (1990); Nielsen Norman Group — "10 Usability Heuristics for User Interface Design" (updated 2020); ISO 9241-110:2020 — Ergonomics of human-system interaction.

## Steps

1. **Assemble 3–5 independent evaluators** — Nielsen (1993) established that a single evaluator finds approximately 35% of problems; 3 evaluators find 60%; 5 evaluators find 75%. Beyond 5, marginal returns diminish sharply. Each evaluator must conduct their review independently before group discussion to avoid anchoring bias.

2. **Define the evaluation scope** — Specify which user flows will be evaluated. A complete product heuristic evaluation is impractical; focus on the highest-risk flows: primary onboarding, core task completion, error recovery, and checkout or conversion paths. Document the scope so all evaluators cover the same ground.

3. **Brief evaluators on the target user context** — Evaluators must approach the interface as the target user, not as design experts. Provide: (a) a one-paragraph user persona, (b) the task scenarios to evaluate, (c) the user's technical sophistication level. An evaluator applying expert-user mental models to a novice-user interface will find the wrong problems.

4. **Conduct two passes through each flow** — First pass: navigate freely, building a sense of the overall interaction structure. Second pass: evaluate methodically against each heuristic, noting problems as they occur. Both passes are necessary — the first pass prevents premature fixation; the second pass is where systematic documentation happens.

5. **Evaluate against all 10 Nielsen heuristics** — For each screen or interaction, check: (1) Visibility of system status; (2) Match between system and real world; (3) User control and freedom; (4) Consistency and standards; (5) Error prevention; (6) Recognition rather than recall; (7) Flexibility and efficiency of use; (8) Aesthetic and minimalist design; (9) Help users recognize, diagnose, and recover from errors; (10) Help and documentation. Document every violation found.

6. **Record problems in a structured log** — Each problem entry must include: (a) heuristic violated (number and name); (b) location (screen name, element); (c) description of the problem; (d) screenshot or annotation; (e) severity rating (0–4 scale). Use a shared spreadsheet so all evaluators log to the same format.

7. **Rate severity using the Nielsen scale** — 0: Not a usability problem; 1: Cosmetic only — fix if time allows; 2: Minor — low priority; 3: Major — important to fix; 4: Usability catastrophe — imperative to fix before release. Severity is a function of: frequency (how often it occurs), impact (how seriously it affects task completion), and persistence (whether it is a one-time or repeated problem).

8. **Aggregate findings across all evaluators** — Combine all individual logs, merge duplicate findings, and calculate the frequency with which each problem was independently identified. Problems found by multiple evaluators are more severe in practice than individual-only finds. Sort the aggregated list by severity × evaluator agreement.

9. **Facilitate a group debrief session** — Bring all evaluators together to review the aggregated list. Discuss borderline severity ratings, resolve disagreements about heuristic categorization, and identify patterns (e.g., all H6 violations occur in the settings section — suggesting a systemic problem, not isolated issues). Group discussion also generates preliminary remediation ideas.

10. **Produce a prioritized findings report** — Deliver: (a) an executive summary (top 5 issues and business impact); (b) the full findings list sorted by severity; (c) screenshot annotations for each finding; (d) one recommended fix per issue. Frame findings in terms of user behavior and impact, not designer preference. "Users will not know whether their action succeeded" (H1 violation) is more actionable than "the feedback is insufficient."

## Rules

- Each evaluator must complete their review independently before seeing other evaluators' findings — group review before individual review produces anchoring bias and reduces total problems found.
- Every finding must cite the specific heuristic violated; findings without heuristic attribution cannot be objectively prioritized.
- Distinguish problems from design preferences — a violation must be traceable to a specific heuristic, not to personal aesthetic opinion.

## Common Mistakes

- **Using only one evaluator** — A single evaluator finds fewer than 40% of problems and is subject to personal blind spots. If budget allows only one evaluator, acknowledge the coverage limitation explicitly in the report.
- **Evaluating as an expert user** — HCI experts who forget to simulate the target user's mental model find problems that do not matter to that user and miss problems that do. Always anchor evaluation to the defined user scenario.
- **Skipping error states and edge cases** — Evaluating only the happy path misses the highest-severity heuristic violations, which cluster in error recovery flows. Deliberately trigger error states: submit empty forms, use invalid inputs, trigger timeouts.
- **Conflating heuristic evaluation with user testing** — These are complementary methods, not substitutes. Heuristic evaluation finds expert-identifiable violations quickly; user testing finds the problems that only real users in real context reveal. Use both, in that order.

## Examples

**H1 violation (System status):** A file upload interface shows a spinner but provides no progress indicator, no estimated time remaining, and no confirmation when the upload completes. Users cannot distinguish a slow upload from a failed one.

**H5 violation (Error prevention):** A form allows users to submit without completing required fields, then shows an error. A constrained design disables the submit button until all required fields pass validation, preventing the error rather than recovering from it.

**H9 violation (Error recovery):** An error message reads: "Error 404: Request failed." It does not explain what went wrong, why it happened, or what the user should do next. A heuristic-compliant message: "We couldn't find that page. It may have moved or been deleted. Return to home or search for what you need."

## When NOT to Use

- As a substitute for all user research — heuristic evaluation is an expert inspection method; it cannot reveal how real users in real contexts actually behave, misunderstand, or work around design problems.
- When evaluators lack sufficient domain familiarity — an evaluator who does not understand the user's context or the product's purpose will misclassify problems and assign incorrect severity ratings.
