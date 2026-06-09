---
name: design-data-quality-framework
description: Use when designing a systematic data quality framework covering completeness, accuracy, consistency, timeliness, validity, and uniqueness across data pipelines and data warehouses
source: 'DAMA International "DMBOK2: Data Management Body of Knowledge" 2nd ed. (2017); ISO 8000-8 "Data Quality" international standard; Redman T. "Data Quality: The Field Guide" (2001); Google "Data Quality Best Practices" (2022)'
tags: [data, data-quality, dmbok, data-engineering, pipelines, governance]
verified: true
---

# Design Data Quality Framework

Establish systematic measurement, monitoring, and remediation of data quality across six dimensions so that data consumers can trust the data they use for decisions and operations.

## Why This Is Best Practice

**Adopted by:** DAMA International (global data management standards body, 60,000+ members), ISO (international standard ISO 8000), financial regulators (BCBS 239 for banking data quality), healthcare (HL7 FHIR data quality requirements), and hyperscalers including Google, Amazon, and Microsoft for internal data platforms

**Impact:** Gartner estimates poor data quality costs organizations an average of $12.9 million per year (2021). IBM found that the US economy loses $3.1 trillion annually to poor data quality. Google's internal data quality initiative reduced data pipeline incidents caused by quality issues by 60% and cut analyst time spent on data validation from 40% to 15% of their workweek. BCBS 239 compliance at major banks required formalized data quality frameworks that reduced regulatory reporting errors by 85%.

**Why best:** Ad hoc data quality checks embedded in individual pipelines or queries are inconsistent, undiscoverable, and unmaintainable at scale. A systematic framework defines quality dimensions uniformly, automates measurement, surfaces results to data consumers, and assigns accountability for remediation — making data quality a first-class engineering and organizational concern.

Sources: DAMA International "DMBOK2" (2017); ISO 8000-8 Data Quality standard; Redman, T.C. "Data Quality: The Field Guide" Digital Press (2001); BCBS 239 "Principles for Effective Risk Data Aggregation and Risk Reporting" (2013); Gartner "How to Improve Your Data Quality" (2021)

## Steps

1. **Define the six quality dimensions — establish a shared vocabulary** — Standardize quality measurement across six DAMA-aligned dimensions: (1) Completeness — are all expected values present? (2) Accuracy — do values match the real-world entity they represent? (3) Consistency — are values consistent across systems and records? (4) Timeliness — is data available when needed, with acceptable freshness? (5) Validity — do values conform to defined formats, ranges, and business rules? (6) Uniqueness — are there duplicates that should be single records? Document definitions and acceptable thresholds for each dimension per data domain.

2. **Inventory critical data elements — focus measurement where it matters** — Identify Critical Data Elements (CDEs): the data fields whose quality directly impacts key business decisions, regulatory reporting, or operational processes. For each CDE, document: business owner, source system, downstream consumers, acceptable quality thresholds per dimension, and business impact of quality failure. Start quality measurement with CDEs before expanding to all data.

3. **Instrument data pipelines with quality checks — automate measurement** — Embed quality checks at key pipeline stages: at ingestion (source data profiling), after transformation (business rule validation), and before consumption (output validation). Use a data quality framework (Great Expectations, dbt tests, Soda Core, or Monte Carlo) to define checks as code, run them automatically, and capture results in a central quality store.

4. **Define quality SLOs — make quality measurable and contractual** — For each CDE and data product, define Service Level Objectives (SLOs) for quality: e.g., "completeness ≥ 99.5%, accuracy ≥ 99%, freshness ≤ 4 hours." Treat quality SLO breaches as incidents with defined severity levels and response SLAs. Publish quality SLOs in the data catalog so consumers know what to expect.

5. **Implement data observability — detect anomalies proactively** — Deploy data observability tooling (Monte Carlo, Bigeye, or dbt-based monitoring) to detect statistical anomalies in data distributions, row counts, null rates, and schema changes without pre-defined rules. Combine rule-based checks (step 3) with ML-based anomaly detection for issues that rules cannot anticipate.

6. **Build a data quality scorecard — surface quality to consumers** — Create a dashboard showing quality scores per data domain, per dimension, and trend over time. Integrate quality scores into the data catalog (Datahub, Amundsen, Atlan) so users see quality metadata alongside dataset descriptions. A dataset with a completeness score of 94% tells a consumer whether it is fit for their purpose.

7. **Assign data stewards — create accountability for quality** — For each data domain, assign a Data Steward responsible for quality: defining business rules, investigating quality failures, coordinating remediation across systems, and approving quality exceptions. Without ownership, quality issues remain unresolved. Stewards are business-side roles empowered with data engineering support.

8. **Implement root cause classification — learn from failures** — When a quality check fails, classify the root cause: source system change, ETL transformation bug, schema drift, upstream data deletion, or business rule change. Maintain a quality incident log. Analyze root causes monthly to identify systemic patterns and prioritize prevention over repeated remediation.

9. **Define remediation workflows — close the feedback loop** — For each quality failure class, define the remediation workflow: automated self-healing (backfill, deduplication), source system correction, quarantine and manual review, or consumer notification and hold. Integrate remediation workflow with the incident management system. Track mean time to resolve (MTTR) quality incidents.

10. **Govern quality through a data quality committee** — Establish a recurring forum (monthly) where data stewards, engineering, and business stakeholders review quality trends, approve threshold changes, prioritize remediation investments, and publish a quality roadmap. Embed data quality objectives into engineering OKRs.

## Rules

- Quality checks must run automatically in CI/CD — no manual quality validation before production deployment
- Every CDE must have a documented owner (data steward) — unowned data has no accountability for quality
- Quality failures must never be silently ignored — every failed check must generate an alert routed to a responsible party
- Quality thresholds must be agreed with data consumers, not set unilaterally by the engineering team
- Raw source data must be preserved in a landing zone before any transformation — enables reprocessing when quality issues are fixed at the source

## Common Mistakes

- **Measuring quality but not acting on it** — Building dashboards that show quality scores without alert routing, ownership, or remediation processes creates reporting theater with no quality improvement.
- **Defining quality without business context** — Setting a 99% completeness threshold without knowing whether 1% missing values in a specific field causes failed transactions or regulatory violations makes the metric meaningless.
- **Over-investing in validation rules upfront** — Writing hundreds of rules before understanding the actual quality failure modes leads to false positives, alert fatigue, and abandonment of the framework. Start with 5–10 rules for each CDE and expand based on observed failure patterns.
- **Ignoring data quality at the source** — Most quality frameworks catch problems late (in the warehouse). Root cause is usually the source system. Invest in upstream data contracts and source system quality agreements to prevent problems at origin.
- **Treating quality as a one-time project** — Data quality degrades continuously as source systems change, volumes grow, and business rules evolve. Quality must be an ongoing operational discipline with recurring measurement and governance.

## Examples

**Great Expectations pipeline gate:** A dbt model outputs a customer table. A Great Expectations suite validates: `expect_column_values_to_not_be_null` for `customer_id`, `expect_column_values_to_be_between` for `age` (18–120), `expect_column_unique_value_count_to_be_between` for duplicate customer keys. The suite runs in CI; if any expectation fails, the deployment is blocked and an alert fires to the data steward.

**Quality SLO breach incident:** A Monte Carlo alert fires: the daily row count of the `orders` table dropped 40% from the 30-day average. The on-call data engineer investigates, finds an upstream API change that dropped a required field causing rejections. The data steward is notified. The source system team applies a fix. The quality SLO breach is logged with root cause "source system schema change" and resolved within the 4-hour SLA.

## When NOT to Use

- Exploratory data science notebooks with single-use data where quality is assessed manually by the analyst
- Real-time event streams where the latency of quality checks would violate throughput requirements — use lightweight statistical sampling instead of full validation
