---
name: audit-training-data-integrity
description: Use when curating, collecting, or managing datasets for training or fine-tuning machine learning models — to detect and mitigate poisoned, biased, mislabeled, or adversarially crafted training data.
source: 'OWASP Top 10 for LLM Applications 2025 LLM03 (owasp.org/www-project-top-10-for-large-language-model-applications/); NIST AI RMF 1.0; Barreno et al. "The Security of Machine Learning" (2010); CWE-20'
tags: [security, owasp, llm, training-data, data-poisoning, ml-security, emerging, developer]
emerging: true
---

# Audit Training Data Integrity

Validate, deduplicate, and provenance-track training data to detect poisoning attacks, privacy violations, and systematic bias before they are baked into model weights.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM03 (Training Data Poisoning). NIST AI RMF 1.0 (2023) Govern 1.1 and Map 2.2 address training data governance. Google, Meta, and Hugging Face publish data cards and datasheets for datasets as standard practice. The ML safety community (papers by Carlini, Wallace, Schuster) has demonstrated practical data poisoning attacks since 2020.
**Status:** Emerging — data poisoning is a well-researched attack (since ~2017) but detection tooling and industry standards are still maturing.
**Impact:** Carlini et al. (2021) demonstrated extracting verbatim training data from GPT-2, including PII. Schuster et al. (2021) showed code suggestion models trained on poisoned GitHub data could insert security vulnerabilities into suggestions (later confirmed in real systems). BadNL attacks (Chen et al., 2021) showed NLP classifiers could be given hidden backdoors via 50 poisoned examples in 100,000 training samples — 0.05% poisoning rate sufficient for reliable triggering.
**Why best:** Post-hoc model auditing (testing the trained model for backdoors) is the alternative — it's expensive, slow, and incomplete. Data-time auditing catches issues before they're encoded into weights, where removal is computationally expensive and may require full retraining.

Sources: OWASP LLM Top 10 2025 LLM03; NIST AI RMF 1.0; Carlini et al. "Extracting Training Data from Large Language Models" (2021); Schuster et al. "You Autocomplete Me: Poisoning Vulnerabilities in Neural Code Completion" (2021)

## Steps

1. **Track data provenance with a data card for every dataset**:

   ```yaml
   # data_card.yaml
   dataset_name: customer-support-v2
   version: 2.3.1
   collected_by: data-team
   collection_date: 2024-01-15
   sources:
     - type: internal
       description: Customer support tickets Q3-Q4 2023
       pii_review: completed
       pii_review_date: 2024-01-10
   licenses: [internal-proprietary]
   known_biases:
     - English-only, over-represents US timezone users
   exclusions:
     - removed: all tickets containing SSN, credit card numbers
     - removed: tickets from users who opted out of data use
   hash_sha256: <hash of final dataset file>
   ```

2. **Scan for PII before training**:

   ```python
   from presidio_analyzer import AnalyzerEngine

   analyzer = AnalyzerEngine()

   def scan_for_pii(text):
       results = analyzer.analyze(text=text, language='en')
       high_confidence = [r for r in results if r.score > 0.85]
       return high_confidence

   def audit_dataset(dataset_path, sample_size=1000):
       pii_count = 0
       for sample in load_dataset_sample(dataset_path, n=sample_size):
           findings = scan_for_pii(sample['text'])
           if findings:
               pii_count += 1
               log_pii_finding(sample['id'], findings)

       pii_rate = pii_count / sample_size
       if pii_rate > 0.01:  # >1% PII rate — investigate
           raise DataQualityError(f"PII rate {pii_rate:.2%} exceeds threshold")
   ```

3. **Detect and remove near-duplicate entries** — duplicates cause memorization and model overfitting:

   ```python
   from datasketch import MinHash, MinHashLSH

   def deduplicate_dataset(texts, threshold=0.9):
       lsh = MinHashLSH(threshold=threshold, num_perm=128)
       unique_texts = []

       for i, text in enumerate(texts):
           m = MinHash(num_perm=128)
           for word in text.split():
               m.update(word.encode('utf8'))

           if not lsh.query(m):  # no near-duplicate found
               lsh.insert(str(i), m)
               unique_texts.append(text)

       print(f"Removed {len(texts) - len(unique_texts)} near-duplicates")
       return unique_texts
   ```

4. **Audit label consistency for supervised learning**:

   ```python
   from sklearn.metrics import cohen_kappa_score

   def audit_label_quality(labels_annotator1, labels_annotator2):
       kappa = cohen_kappa_score(labels_annotator1, labels_annotator2)
       if kappa < 0.6:
           raise DataQualityError(
               f"Inter-annotator agreement κ={kappa:.2f} below threshold 0.6"
           )
       return kappa

   # Flag for review: examples where annotators disagree
   def find_disagreements(data, labels1, labels2):
       return [data[i] for i in range(len(data)) if labels1[i] != labels2[i]]
   ```

5. **Test for backdoor triggers in the trained model** (post-training validation):

   ```python
   def test_for_backdoor(model, trigger_phrases, expected_safe_output):
       """Check if specific phrases cause anomalous outputs (backdoor indicator)."""
       anomalous = []
       for trigger in trigger_phrases:
           output = model.generate(f"Normal sentence with {trigger}")
           if output != expected_safe_output:
               anomalous.append({'trigger': trigger, 'output': output})
       return anomalous
   ```

6. **Implement consent and opt-out tracking for user-generated training data**:

   ```python
   def include_in_training(user_id, content):
       consent = db.get_training_consent(user_id)
       if consent is None or consent.opted_out:
           return False
       if content.created_at < consent.consent_date:
           return False  # pre-dates consent
       return True
   ```

## Rules

- Cryptographically hash and version-stamp each training dataset snapshot — enables auditing exactly which data was used if a model issue is discovered post-deployment.
- Data contributed by users requires explicit consent for training use — check privacy policy and applicable regulations (GDPR Article 5, CCPA).
- Regular expressions alone are insufficient for PII detection — use NER-based tools (Presidio, spaCy) for higher coverage.
- Poisoning attacks are designed to be undetectable in individual samples — statistical methods on the full dataset are required.

## Common Mistakes

- **Treating training data as inherently trusted because it's "internal"** — internal datasets sourced from user content or web scrapes can contain adversarial examples.
- **Not tracking data provenance** — when a model exhibits a problem behavior, inability to trace which training data caused it makes remediation impossible.
- **Using the full training set for deduplication checks** — exact deduplication misses near-duplicates; use MinHash LSH or similar approximate methods.
- **Skipping PII scanning because data is "anonymized"** — anonymization is often incomplete; structured PII like SSNs and credit cards persist through naive scrubbing.
