---
name: design-data-pipeline
description: Use when designing, building, or reviewing a data pipeline for ingestion, transformation, or analytics
source: Maxime Beauchemin "Functional Data Engineering" (medium.com/@maximebeauchemin, 2018); dbt documentation (docs.getdbt.com); Apache Airflow documentation
tags: [data-engineering, data-pipeline, dbt, airflow, etl, elt, functional-data-engineering]
verified: true
---

# Design Data Pipeline

Design data pipelines using functional principles — idempotency, immutability, and declarative transformations — for reliability and maintainability.

## Why This Is Best Practice

**Adopted by:** Airbnb (Airflow originator), Fishtown Analytics (dbt originator, now dbt Labs), Lyft, GitLab
**Impact:** Beauchemin's functional data engineering principles, adopted by thousands of data teams, eliminate an entire class of pipeline bugs (non-idempotent transforms, mutable state) that cause silent data corruption. dbt's adoption grew from 0 to 30,000+ companies in 5 years due to its application of software engineering practices to data transformation.

**Why best:** Traditional ETL pipelines are stateful, brittle, and difficult to test. Functional data engineering applies software engineering principles: pure transformations (same input → same output), immutable historical data, idempotent operations (safe to re-run), and declarative SQL-based transforms that are version-controlled and testable.

## Steps

1. **Choose ELT over ETL** — Load raw data into the warehouse first (ELT), then transform in-warehouse using SQL/dbt. ELT preserves the raw audit trail and uses the warehouse's compute efficiently.
2. **Design idempotent stages** — Every pipeline stage must be safely re-runnable: use `INSERT OVERWRITE` or `CREATE OR REPLACE` semantics; never `INSERT APPEND` without deduplication guards.
3. **Partition by time** — Partition tables by event date or processing date. Load and reprocess data in date partitions; never update historical rows — append corrections with a `processed_at` timestamp.
4. **Model in layers (dbt convention)** — Staging (raw → typed, renamed) → Intermediate (business logic joins) → Marts (aggregated, analytics-ready). Each layer is independently testable.
5. **Add data quality tests** — Test at the staging layer: not-null, unique, accepted-values, referential integrity. Use dbt tests or Great Expectations. Fail the pipeline on data quality violations.
6. **Orchestrate with DAGs** — Define dependencies as directed acyclic graphs (Airflow, Prefect, Dagster). No pipelines with implicit ordering — explicit dependency graphs are observable and retryable.
7. **Monitor data freshness and volume** — Alert on: SLA breach (table not updated within expected window), row count anomaly (>20% deviation from baseline), and schema change.

## Rules

- Raw data is immutable — never overwrite source data; append new loads with `ingested_at` timestamps.
- Every transformation must be deterministic — the same source data must always produce the same output.
- Pipeline code is version-controlled — no notebook-only pipelines in production; promote to .sql or .py files in a repo.
- Separate ingestion, transformation, and serving layers — each can be monitored, tested, and scaled independently.

## Examples

dbt model structure:
```
models/
  staging/
    stg_orders.sql       -- raw → typed, renamed
    stg_payments.sql
  intermediate/
    int_order_items.sql  -- join orders + line items
  marts/
    fct_orders.sql       -- fact table for analytics
    dim_customers.sql    -- dimension table
```

Idempotent load (BigQuery):
```sql
CREATE OR REPLACE TABLE `project.dataset.orders_2026_03_01`
AS SELECT * FROM `project.raw.orders`
WHERE DATE(created_at) = '2026-03-01';
```

## Common Mistakes

- **Append-only pipelines without deduplication** — re-running on failure creates duplicate rows; always use `INSERT OVERWRITE` or MERGE semantics.
- **No data quality tests** — silent schema changes (a column renamed upstream) corrupt downstream models without failing the pipeline.
- **God DAG** — one Airflow DAG with 200 tasks; impossible to debug; split by domain and layer.
