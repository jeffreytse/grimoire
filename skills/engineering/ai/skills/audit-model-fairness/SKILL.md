---
name: audit-model-fairness
description: Use when deploying ML models that make decisions affecting people (hiring, lending, healthcare, content moderation) or when required by regulation to assess algorithmic bias
source: NIST AI Risk Management Framework (AI RMF 1.0, 2023); Barocas & Hardt "Fairness in Machine Learning" (2019); IEEE Ethically Aligned Design
tags: [ai, fairness, ethics, compliance]
verified: true
---

# Audit Model Fairness

Systematically measure and document a model's performance across demographic groups to identify discriminatory outcomes before or after deployment.

## Why This Is Best Practice

**Adopted by:** Required by EU AI Act (2024) for high-risk AI systems; CFPB fair lending requirements; EEOC guidelines for employment AI; NIST AI RMF adopted by US federal agencies
**Impact:** Biased models create legal liability (CFPB fines up to $1M/day for fair lending violations); Amazon famously retracted an AI hiring tool after discovering gender bias; proactive audits prevent reputational and regulatory harm
**Why best:** Models trained on historical data encode historical discrimination; without measurement, unfairness is invisible until harm occurs

Sources: NIST AI RMF 1.0 (2023); Barocas, Hardt & Narayanan "Fairness and Machine Learning" (2019); IEEE Ethically Aligned Design v2 (2019)

## Steps

1. **Define the protected attributes** — Identify legally and ethically relevant attributes for your context: race, gender, age, disability status, national origin, religion, sexual orientation. Determine which you can directly measure and which must be inferred from proxies. Document legal basis and jurisdiction.

2. **Select fairness metrics** — Choose metrics appropriate to the decision context. Demographic parity: equal positive prediction rates across groups (appropriate for representation goals). Equalized odds: equal TPR and FPR across groups (appropriate for classification). Calibration: equal prediction accuracy across groups (appropriate for risk scoring). No single metric satisfies all simultaneously (impossibility theorem); choose based on harm type.

3. **Assemble a stratified evaluation dataset** — Evaluation data must be representative of the deployment population. Oversample minority groups to ensure statistical significance (minimum 100 samples per subgroup for meaningful metrics). Use held-out data, not training data. Document dataset construction methodology.

4. **Measure overall model performance** — Establish baseline accuracy, precision, recall, and AUC for the full population. This is the reference point for group-level comparisons. Document evaluation date and model version.

5. **Measure per-group performance** — Compute the same performance metrics for every protected group. Calculate disparity ratios: group metric / majority group metric. Flag disparities above 0.8 (80% rule, EEOC 4/5ths rule) as potential adverse impact. Visualize as a fairness dashboard.

6. **Investigate sources of disparity** — Analyze: Is disparity in the training data (historical bias)? In feature selection (proxy discrimination)? In model architecture? In label quality (human labeling bias)? Use SHAP values to identify which features drive differential predictions across groups.

7. **Apply mitigation techniques** — Pre-processing: reweight training data, resample underrepresented groups. In-processing: add fairness constraints to the loss function (adversarial debiasing, regularization). Post-processing: adjust decision thresholds per group to equalize error rates. Document trade-offs with overall model performance.

8. **Conduct human review of edge cases** — Sample 50-100 misclassified cases per protected group. Have domain experts review for patterns. Automated metrics miss contextual harms that human review surfaces (e.g., stereotyped language in text models).

9. **Produce a model card or fairness audit report** — Document: model purpose, intended use, evaluation methodology, per-group performance metrics, known limitations, and mitigation steps taken. Publish internally and externally per your disclosure policy. EU AI Act requires this for high-risk systems.

10. **Establish ongoing monitoring** — Deploy fairness monitoring in production. Track per-group prediction distributions monthly. Set alerts if demographic disparity increases above threshold post-deployment. Retrain or retune when drift is detected. Fairness is not a one-time audit.

## Rules

- No model affecting consequential decisions (employment, credit, healthcare) should be deployed without a completed fairness audit.
- Fairness and accuracy are often in tension; document the trade-off explicitly and let stakeholders make the decision — do not hide it.
- Never use protected attributes as model features in jurisdictions where this is prohibited by law; proxy variables (zip code as income proxy) can have equivalent discriminatory effect.
- Consult legal counsel before publishing fairness audit results; findings may create legal exposure and require qualified communication.

## Common Mistakes

- **Measuring only aggregate accuracy** — a model with 90% accuracy may have 95% accuracy for one group and 75% for another; aggregate metrics obscure group-level harm.
- **Selecting fairness metrics after seeing results** — choose metrics before analysis; post-hoc selection biases toward favorable metrics.
- **Treating a one-time audit as sufficient** — distribution shift in production data causes fairness properties to degrade; only ongoing monitoring detects this.
- **Conflating correlation with causation** — a model performs worse for a group because of data or feature issues, not because the group is inherently harder to predict.

## When NOT to Use

- Models with no human impact (weather prediction, scientific simulation, internal analytics without decision-making use)
- Contexts where protected attribute data is legally prohibited from collection and proxy analysis is not feasible
