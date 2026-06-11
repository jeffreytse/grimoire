---
name: validate-llm-output
description: Use when LLM outputs are used to make decisions, generate content, answer questions, or drive application logic — to detect hallucinations, verify factual grounding, and prevent over-reliance on unverified AI-generated content.
source: 'OWASP Top 10 for LLM Applications 2025 LLM09 (owasp.org/www-project-top-10-for-large-language-model-applications/); NIST AI RMF 1.0 Measure 2.5; Anthropic Constitutional AI research; CWE-1188'
tags: [security, owasp, llm, hallucination, output-validation, grounding, ai-safety, emerging, developer]
emerging: true
---

# Validate LLM Output

Verify LLM outputs against ground truth sources, enforce structured schemas, and implement human review thresholds — preventing hallucinated facts, fabricated citations, and incorrect decisions from propagating into production systems.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM09 (Overreliance). NIST AI RMF 1.0 (2023) Measure 2.5 requires AI output monitoring and human review mechanisms. The EU AI Act (2024) Article 9 mandates human oversight for high-risk AI system outputs. Microsoft Azure AI Content Safety, AWS Bedrock Guardrails, and Google Vertex AI Evaluation all provide output validation tooling.
**Status:** Emerging — hallucination detection and output validation are active research areas; no single complete solution exists, and techniques are still rapidly improving.
**Impact:** A US federal attorney was sanctioned by a judge for submitting AI-generated legal briefs containing fabricated case citations (Mata v. Avianca, 2023). A medical AI system gave drug dosage recommendations that were factually incorrect due to hallucination. LLMs confidently generate plausible-sounding but false statistics, fake academic citations, and invented API parameters — all at a high confidence level that can mislead users and downstream systems.
**Why best:** User-level review for every output is the alternative — it's too slow for real-time applications and subject to automation bias (users trust AI outputs without scrutinizing them). Automated structural validation + confidence thresholds + selective human review provides scalable quality control.

Sources: OWASP LLM Top 10 2025 LLM09; Mata v. Avianca (2023); NIST AI RMF 1.0 Measure 2.5; EU AI Act Article 9

## Steps

1. **Use structured output with schema validation**:

   ```python
   from pydantic import BaseModel, validator
   from typing import Optional

   class ProductRecommendation(BaseModel):
       product_id: str  # must reference a real product
       reason: str
       confidence: float  # 0.0–1.0

       @validator('product_id')
       def product_must_exist(cls, v):
           if not db.product_exists(v):
               raise ValueError(f"LLM hallucinated product ID: {v}")
           return v

       @validator('confidence')
       def confidence_in_range(cls, v):
           if not 0.0 <= v <= 1.0:
               raise ValueError("Confidence must be 0–1")
           return v

   response = openai.beta.chat.completions.parse(
       model="gpt-4o",
       messages=messages,
       response_format=ProductRecommendation,
   )
   result = response.choices[0].message.parsed
   # If product_id doesn't exist in DB, pydantic raises ValueError — handle it
   ```

2. **Verify citations and references against source documents**:

   ```python
   def verify_citations(llm_response: str, source_documents: list[str]) -> dict:
       """Check that claims in LLM response are grounded in source documents."""
       # Extract claims from response
       claims = extract_claims(llm_response)
       verification_results = {}

       for claim in claims:
           # Check if claim is supported by any source document
           supported = any(
               is_semantically_similar(claim, chunk)
               for doc in source_documents
               for chunk in split_into_chunks(doc)
           )
           verification_results[claim] = 'supported' if supported else 'unverified'

       unverified_count = sum(1 for v in verification_results.values() if v == 'unverified')
       return {
           'verified_fraction': 1 - unverified_count / max(len(claims), 1),
           'claims': verification_results
       }
   ```

3. **Flag low-confidence outputs for human review**:

   ```python
   CONFIDENCE_THRESHOLDS = {
       'medical_advice': 0.95,
       'legal_summary': 0.90,
       'financial_analysis': 0.90,
       'general_query': 0.70,
   }

   def route_for_review(output: dict, task_type: str) -> dict:
       threshold = CONFIDENCE_THRESHOLDS.get(task_type, 0.80)
       confidence = output.get('confidence', 0.0)

       if confidence < threshold:
           queue_for_human_review(output, task_type, reason='low_confidence')
           return {
               'status': 'pending_review',
               'message': 'This response requires expert review before use.',
               'estimated_review_time': '2 hours'
           }
       return output
   ```

4. **Cross-validate with deterministic checks where possible**:

   ```python
   def validate_llm_analysis(llm_output: dict, raw_data: dict) -> bool:
       """Re-check LLM's numerical claims against source data."""
       # LLM said revenue grew 15% — verify against actual numbers
       if 'revenue_growth' in llm_output:
           actual_growth = (
               (raw_data['revenue_2024'] - raw_data['revenue_2023'])
               / raw_data['revenue_2023']
           )
           llm_claimed = llm_output['revenue_growth']
           if abs(actual_growth - llm_claimed) > 0.02:  # 2% tolerance
               logger.warning("LLM hallucinated revenue growth: claimed=%s actual=%s",
                              llm_claimed, actual_growth)
               return False
       return True
   ```

5. **Set appropriate user expectations** — surface uncertainty in the UI:

   ```python
   def format_response_with_confidence(response: dict) -> dict:
       confidence = response.get('confidence', 0.0)
       disclaimer = None

       if confidence < 0.70:
           disclaimer = "⚠️ This response has low confidence. Verify before acting on it."
       elif confidence < 0.85:
           disclaimer = "This response is AI-generated. Please verify important details."

       return {
           **response,
           'disclaimer': disclaimer,
           'sources_verified': response.get('sources_grounded', False),
       }
   ```

6. **Monitor output quality over time** — hallucination rates drift with model updates:

   ```python
   def sample_for_quality_audit(response_id: str, sample_rate: float = 0.05):
       """Send a sample of responses for human quality review."""
       if random.random() < sample_rate:
           queue_for_quality_review(response_id)
   ```

## Rules

- Never use LLM outputs as authoritative sources for safety-critical decisions (medical dosing, legal advice, financial transactions) without human review.
- "The model said it with confidence" is not validation — LLMs express high confidence on hallucinated content.
- Schema validation catches structural errors; it does not verify factual accuracy. Both are needed.
- Grounding (RAG + citation verification) is currently the most effective hallucination reduction technique for factual queries.

## Common Mistakes

- **Displaying LLM-generated citations as clickable links without verifying they exist** — hallucinated URLs and DOIs are common.
- **Using LLM output as input to another LLM without validation** — hallucinations compound across chains.
- **Treating structured output (JSON mode) as validated output** — the structure is correct, but the values may still be hallucinated.
- **No feedback loop for incorrect outputs** — without tracking which outputs were wrong, you cannot improve the system.
