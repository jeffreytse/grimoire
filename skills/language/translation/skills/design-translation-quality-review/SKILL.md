---
name: design-translation-quality-review
description: Use when establishing a quality review process for translated content to ensure accuracy, fluency, and cultural appropriateness before publication.
source: ISO 17100:2015 (Translation Services Requirements); ASTM F2575-14 (Quality Aspects in Translation); Localization Industry Standards Association (LISA) QA model; TAUS (Translation Automation User Society) DQF framework
tags: [translation, localization, quality, review, language, communication, internationalization]
verified: true
---

# Design Translation Quality Review

Design a systematic translation quality review process that catches accuracy errors, fluency issues, and cultural missteps before translated content reaches its audience.

## Why This Is Best Practice

**Adopted by:** ISO 17100:2015 (international translation services standard, adopted in 40+ countries), EU Directorate-General for Translation (produces 2M+ pages annually), ASTM International translation quality standard, major localization buyers including Microsoft, Adobe, SAP, and the United Nations translation services.
**Impact:** Errors in translated content cost organizations an average of $50,000 per major legal or regulatory mistranslation (GALA Globalization survey 2020). ISO 17100-compliant processes reduce translation error rates by 60–80% versus unreviewed single-translator output (CSA Research 2019). A two-step translation + editing workflow is the minimum standard for any content with legal, medical, or safety implications. TAUS data shows structured QA reduces post-publication correction costs by 70%.
**Why best:** Machine translation and non-native translator output both fail in predictable ways: literal translation of idioms, false cognates, culturally inappropriate register, and terminology inconsistency. A designed review process assigns each error type to the right reviewer — a native-speaking editor for fluency, a domain expert for terminology, a local cultural consultant for appropriateness — rather than hoping one reviewer catches everything.

Sources: ISO 17100:2015; ASTM F2575-14 Standard Guide for Quality Aspects in Translation; TAUS Dynamic Quality Framework (DQF); CSA Research "Translation Technology Insights" (2019); Colina "Translation Quality Evaluation" (2015)

## Steps

1. **Define quality requirements by content type** — Legal/regulatory and medical content requires full TEP (Translation, Editing, Proofreading) with subject-matter expert review. Marketing/brand content needs native-speaker fluency review plus cultural adaptation check. Internal/informational content can use post-edit machine translation with light review. Match review depth to stakes.
2. **Select qualified reviewers** — Editor must be a native speaker of the target language, different from the original translator, with domain expertise. A Spanish-speaking engineer should not edit a medical translation — domain knowledge and linguistic fluency are both required. Match reviewer expertise to content domain.
3. **Create a style guide and glossary** — Before reviewing, define: approved terminology for key terms (product names, technical terms, legal phrases), preferred register (formal vs. informal address), brand voice guidelines, and terms that must not be translated (trademarks, proper nouns). Reviewers can only apply consistent standards if standards exist.
4. **Run accuracy review** — Compare target text to source text segment by segment. Flag: omissions (content present in source but missing in target), additions (content in target not in source), mistranslations (meaning changed), and terminology errors (wrong term per glossary). Use a translation memory tool (Memsource, memoQ, SDL Trados) for segment-by-segment comparison.
5. **Run fluency and naturalness review** — Read target text without reference to source. Does it read like native writing or translated text? Flag: unnatural word order, calques (literal idiom translation), inconsistent register, awkward phrasing. This review is done by a native speaker who may not read the source language.
6. **Run cultural and localization review** — Check: date/time/number formats match locale (DD/MM/YYYY vs. MM/DD/YYYY), currency symbols, units of measurement, color and image cultural associations, legal disclaimers applicable to target market, culturally sensitive references (humor, idioms, taboos). This requires local market knowledge.
7. **Score and categorize errors** — Use a standardized error typology (MQM — Multidimensional Quality Metrics, or LISA QA model) to rate each error by type and severity: critical (changes meaning or creates safety/legal risk), major (significantly impairs understanding), minor (style or preference). Calculate a quality score for pass/fail determination.
8. **Provide structured feedback to translator** — Return errors with: source segment, target segment, error type, severity, and corrected version. Aggregate errors across a batch to identify systematic patterns (e.g., consistently wrong term) that indicate training needs, not one-time mistakes.

## Rules

- Never have the original translator review their own work — they have the same blind spots as their translation; independent review is the quality control mechanism.
- The glossary must exist before translation begins, not during review — retroactively applying terminology standards to finished translations doubles rework.
- Error severity must be defined before review starts — "critical" and "major" mean different things to different reviewers unless operationally defined.
- Machine translation output always requires human post-editing before publication — raw MT output contains systematic errors that automated QA tools cannot reliably catch.

## Common Mistakes

- **Reviewing for preference, not quality** — Reviewers who change translations to their personal style rather than correcting genuine errors undermine the translator and inflate revision cycles. Reviews must be grounded in the style guide and error typology.
- **Skipping cultural review for "similar" languages** — Brazilian and European Portuguese, Mexican and Spanish Spanish, and Simplified and Traditional Chinese share grammar but differ significantly in vocabulary, idiom, and cultural context. "Close enough" localization generates customer complaints.
- **Using monolingual review for technical content** — A fluency reviewer who cannot read the source cannot catch omissions or meaning changes. Accuracy review requires bilingual reviewers.
- **Not tracking error patterns over time** — One-time correction is not quality improvement. Error pattern analysis across batches identifies whether issues are translator skill gaps, ambiguous source content, or missing glossary terms — each requiring a different fix.

## Examples

**Software UI localization:** Translation team produces 5,000 strings in Spanish. Accuracy reviewer (bilingual) checks 10% sample for mistranslations — finds button labels using formal "usted" when UI style guide specifies informal "tú". Full batch flagged for terminology correction before fluency review proceeds.

**Legal contract translation:** English contract → French translation. Domain expert (bilingual lawyer) reviews liability clauses. Fluency reviewer (native French speaker) reviews readability. Cultural review confirms French legal terminology matches French jurisdiction requirements, not Quebec or Belgian equivalents.

## When NOT to Use

- When producing draft-quality translations for internal reference only — light post-edited machine translation with a disclaimer is sufficient and review cost is unjustified.
- When the translator and target audience share a very close dialect and the content is low-stakes — a brief spot-check replaces full TEP workflow.
- When speed is the only constraint and accuracy is acknowledged as non-critical — document the risk acceptance explicitly rather than silently skipping the process.
