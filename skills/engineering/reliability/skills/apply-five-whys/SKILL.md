---
name: apply-five-whys
description: Use when investigating the root cause of a defect, incident, or recurring failure and need to move beyond symptoms to the underlying systemic cause.
source: Taiichi Ohno, "Toyota Production System" (1978); Sakichi Toyoda (original inventor); widely adopted in lean manufacturing, SRE, and quality engineering
tags: [root-cause-analysis, incident-response, reliability, debugging, lean, quality, defect-reduction]
verified: true
---

# Apply Five Whys

Iteratively ask "why?" five times to trace a failure from its surface symptom to its root systemic cause.

## Why This Is Best Practice

**Adopted by:** Toyota (originator), Google SRE, Amazon, every organization operating under ISO 9001 quality management, and standard practice in lean manufacturing globally.
**Impact:** Toyota's implementation of Five Whys reduced defect recurrence rates by 70%+ by shifting focus from symptom suppression to systemic fix. Google SRE teams credit iterative root cause analysis as the primary driver of reducing repeat incidents; the Google SRE workbook mandates it as part of postmortem procedure.
**Why best:** Most incident responses fix the immediate symptom — the server that crashed, the query that timed out. Without root cause analysis, the same failure recurs in a different form weeks later. Five Whys breaks the symptom–fix–recurrence loop by forcing every investigation to reach a systemic or process-level cause before closing.

Sources: Taiichi Ohno "Toyota Production System" (1978); Google SRE Workbook Chapter 10; ISO 9001:2015 §10.2

## Steps

1. **State the problem precisely** — Write the failure in one sentence with observable facts: not "the site was slow" but "API p99 latency exceeded 2s for 14 minutes starting 14:32 UTC."
2. **Ask the first Why** — Why did this happen? Write the answer as a complete factual statement, not a guess. If you don't know, investigate before answering.
3. **Ask Why again** — Why did that cause occur? Each answer must be a direct, verifiable cause of the previous statement — not a parallel factor, not speculation.
4. **Repeat until systemic cause is reached** — Continue until the answer points to a process, standard, or system design that can be changed. Typical depth: 4–6 iterations. Stop when further "why" produces "human error" or "we don't have a process" — those are root causes.
5. **Verify the causal chain** — Read the chain top to bottom: "X caused Y, which caused Z, which caused the failure." If any step feels like a leap, insert another why.
6. **Identify the fix at the root** — The corrective action targets the deepest cause, not the symptom. A fix at symptom level is a workaround; a fix at root level prevents recurrence.
7. **Document the chain** — Record all five (or more) levels with the causal links written explicitly. This is the postmortem evidence, not just the conclusion.
8. **Validate the fix prevents recurrence** — After implementing, confirm the failure mode cannot repeat via the same chain. If it can, the root cause was misidentified.

## Rules

- Never accept "human error" as the final root cause — ask why the human error was possible (missing safeguard, absent process, inadequate training) and continue.
- Each why must have a verifiable answer before proceeding — stop and investigate rather than guess.
- A single incident may have multiple causal chains; run Five Whys separately on each branch.
- The fix must address the root cause level, not an intermediate level — otherwise the problem will recur via a different path.
- Five is a heuristic, not a rule — use four or seven if that's where the systemic cause genuinely lies.

## Common Mistakes

- **Stopping at the symptom level** — "The database crashed" answered with "restart the database." No root cause found, failure recurs within weeks.
- **Branching into multiple whys per level** — Pick one causal chain at a time; parallel branches become a separate Five Whys exercise.
- **Accepting blame as an answer** — "The engineer deployed untested code" is not a root cause; it's the second why. Ask: why was untested code deployable?
- **Working backwards from a predetermined fix** — Starting with the solution you already want and fitting the why chain around it. The chain must be built forward from evidence.

## When NOT to Use

- When the failure is a known, previously root-caused issue with a documented fix — apply the fix directly rather than re-running the analysis.
- When the incident is still active — stabilize and resolve first, then apply Five Whys in the postmortem phase.
- When the causal chain requires specialized domain expertise you don't have — bring in the domain expert before conducting the analysis, not after.
