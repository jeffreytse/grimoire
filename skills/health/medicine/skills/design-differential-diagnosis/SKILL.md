---
name: design-differential-diagnosis
description: Use when generating and systematically ranking the possible diagnoses for a patient's presenting symptoms and clinical findings
source: Harrison's "Principles of Internal Medicine" (21st ed.); Kassirer & Kopelman "Learning Clinical Reasoning" (1991); Sackett "Evidence-Based Medicine" (2000) clinical decision-making
tags: [medicine, diagnosis, clinical-reasoning, differential-diagnosis]
verified: true
---

# Design Differential Diagnosis

Generate a comprehensive, prioritized differential diagnosis using systematic anatomical, physiological, and probabilistic frameworks, then narrow it with targeted evaluation.

## Why This Is Best Practice

**Adopted by:** USMLE clinical reasoning assessment, Royal College of Physicians and Surgeons certification (Canada, UK), Johns Hopkins/Mayo Clinic clinical reasoning curricula, WHO IMCI (Integrated Management of Childhood Illness) diagnostic algorithms.

**Impact:** Systematic differential diagnosis reduces diagnostic error — which affects 12 million Americans/year (IOM 2015) — by 40% compared to pattern-recognition-only approaches (Graber et al. Arch Int Med 2005); premature closure (settling on a diagnosis too early) accounts for 36% of diagnostic errors.

**Why best:** Structured differential generation forces clinicians to consider diagnoses outside their initial pattern match, preventing anchoring bias and premature closure — the two most common cognitive errors in diagnosis. Probabilistic ranking ensures the most dangerous conditions are evaluated first.

Sources: Harrison's 21st ed. Part 1; Kassirer & Kopelman (1991) ch. 3–5; Sackett et al. (2000) ch. 3; Graber et al. Arch Intern Med 165:1493–1499 (2005).

## Steps

1. **Extract the problem representation** — create a one-sentence clinical summary: "A [age][sex] with [key risk factors] presents with [duration] [chief complaint] plus [2–3 key associated findings], most notable for [pivotal finding]." This activates the correct disease schema.

2. **Generate the initial differential broadly** — list all conditions that could explain the chief complaint; use a systematic framework to avoid omission:
   - Anatomical: work through each organ system that could produce this symptom
   - Pathophysiological (VITAMIN-D): Vascular, Infectious, Traumatic, Autoimmune/Metabolic, Iatrogenic/Idiopathic, Neoplastic, Degenerative/Drug, Congenital/Endocrine
   - Pattern-based: serious conditions not to miss + most common conditions for this demographic.

3. **Apply pre-test probability** — for each candidate diagnosis, estimate base rate given: age, sex, risk factors, geographic prevalence, and referral context. High prior probability diagnoses deserve more prominent position even with fewer specific findings.

4. **Identify pivot features** — find 2–3 findings that most powerfully discriminate between diagnoses: features that are highly specific (LR+>10) for one diagnosis or highly sensitive (LR-<0.1) for ruling out another. Pivot features drive the diagnostic workup.

5. **Rank the differential in three tiers:**
   - Must-not-miss: life-threatening or limb-threatening conditions with even low probability — evaluate first (e.g., PE, MI, meningitis, ectopic pregnancy)
   - Most likely: highest posterior probability given all clinical data
   - Possible: plausible but lower probability or less urgent.

6. **Apply LR to update probabilities** — for key clinical findings and test results: post-test odds = pre-test odds × LR+/−. Use published LR values (EvidenceAlerts, DynaMed); a LR+10 with 10% pre-test probability gives 53% post-test probability.

7. **Identify the discriminating workup** — order tests that have the highest LR+ for the top "must-not-miss" and most-likely diagnoses; avoid ordering tests that won't change management regardless of result.

8. **Apply diagnostic thresholds** — test when: uncertainty is high enough to warrant testing (>test threshold) but not high enough to treat without confirmation (below treatment threshold); treat without testing when probability is above treatment threshold in time-sensitive conditions.

9. **Iteratively update after each result** — each finding (positive or negative) updates probabilities; re-rank the differential with each new piece of information. Do not anchor to the initial most-likely diagnosis if contradictory evidence accumulates.

10. **Document clinical reasoning** — in the Assessment section of the note, state the leading diagnosis with supporting evidence, the differential with key distinguishing features, and the rationale for chosen workup. Undocumented reasoning is indistinguishable from no reasoning.

## Rules

- Always include the most dangerous diagnosis on the differential, even if probability is low — missing aortic dissection at 2% probability is more consequential than missing tension headache at 40%.
- Premature closure is the most dangerous diagnostic error — explicitly ask "What else could this be?" after reaching a leading diagnosis.
- Treat the patient, not the diagnosis — if clinical status deteriorates, widen the differential rather than explaining deterioration as a complication of the current diagnosis.
- Negative test results rule out diagnoses only if test sensitivity is high enough — a negative D-dimer rules out PE only when pre-test probability is low; a negative D-dimer with high pre-test probability does not.

## Common Mistakes

- **Anchoring** — fixing on the first plausible diagnosis and seeking confirmatory evidence while ignoring contradictory findings; corrected by explicitly generating alternatives before ordering tests.
- **Availability bias** — over-weighting diagnoses seen recently or memorably; corrected by using base rate and systematic differential generation rather than pattern recall alone.
- **Overordering without a diagnostic strategy** — ordering a "shotgun" panel without knowing which results would change management generates noise, cost, and false positives without improving diagnostic accuracy.
- **Forgetting zebras in high-risk contexts** — rare diagnoses are rare in primary care but common in referral centers and in patients who have already been evaluated; context shifts the prior.

## When NOT to Use

- In a time-critical emergency where pattern recognition and immediate treatment take precedence over comprehensive differential generation (treat shock first, refine diagnosis after stabilization)
- For straightforward, high-specificity presentations where the clinical picture is pathognomonic (e.g., classic shingles rash — treat, do not generate a broad differential)
- As a standalone AI exercise for individual patient medical decisions — differential diagnosis for real patients requires licensed clinical assessment, physical examination, and full access to the complete history
