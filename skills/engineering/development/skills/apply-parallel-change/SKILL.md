---
name: apply-parallel-change
description: Use when evolving a public API, shared interface, or data schema in a way that must remain backward-compatible with existing consumers during the migration period.
source: Fowler "ParallelChange" (martinfowler.com/bliki/ParallelChange.html, 2014); also known as Expand-Contract or Make-Then-Break pattern; Humble & Farley "Continuous Delivery" (2010)
tags: [development, api-design, refactoring, backward-compatibility, continuous-delivery, patterns, migration, schema]
verified: true
---

# Apply Parallel Change

Evolve an API, interface, or schema without breaking existing consumers by following a three-phase sequence: expand (add new), migrate (move consumers), contract (remove old) — keeping the system functional throughout.

## Why This Is Best Practice

**Adopted by:** Google (standard internal API evolution practice across services), documented by Fowler as "ParallelChange" (also known as Expand-Contract), used across microservices ecosystems, database migration frameworks (Flyway, Liquibase), and REST API versioning in high-traffic systems. The pattern underlies every zero-downtime database schema migration.
**Impact:** Breaking API changes cause integration failures that affect all consumers simultaneously — a common source of production incidents during deployments. Parallel Change eliminates this failure class entirely by ensuring both old and new interfaces exist simultaneously during migration. Google's internal data shows that API-breaking incidents account for 15–25% of service outages; Parallel Change prevents this category. Zero-downtime database migrations at companies like GitHub and Stripe are built on this pattern.
**Why best:** The alternative to Parallel Change is a coordinated cutover: freeze all consumers, migrate the API, redeploy everyone simultaneously. At any meaningful scale, this is operationally infeasible and risky. Parallel Change converts a simultaneous, coordinated, high-risk change into a sequential, self-paced, low-risk migration — each step is independently deployable and reversible.

Sources: Fowler "ParallelChange" (martinfowler.com/bliki/ParallelChange.html 2014); Humble & Farley "Continuous Delivery" (2010) Ch. 12; Sadalage & Fowler "Refactoring Databases" (2006); Expand-Contract pattern (continuousdelivery.com); Nygard "Release It!" (2018) on versioning

## Steps

1. **Identify the breaking change** — Define exactly what is changing and what existing consumers depend on. Breaking changes include: renaming a field, removing a field, changing a field's type, changing a method signature, splitting a table column, changing an enum value. Non-breaking changes (additive) do not require Parallel Change.
2. **Phase 1 — Expand: add the new alongside the old** — Introduce the new API element without removing the old. For a field rename: add the new field name while keeping the old. For a method signature change: add the new overload while keeping the old. For a database column rename: add the new column while keeping the old, and write code that populates both columns on every write. Existing consumers are unaffected because nothing was removed.
3. **Deploy the expansion** — Deploy the expanded system with both old and new elements active. Verify that existing consumers work correctly with no changes. The expand phase must be deployed and stable before any consumer migration begins.
4. **Phase 2 — Migrate: move consumers to the new** — Update each consumer to use the new API element, one at a time. Deploy each consumer update independently. For internal consumers (same codebase or same organization): migrate all. For external consumers (third-party clients): provide migration guides, deprecation notices, and a defined sunset timeline. During migration, both old and new elements must remain active.
5. **Verify migration completeness** — Confirm no active consumers use the old element. Methods: (a) add deprecation warnings/logging to the old element and monitor for invocations, (b) query access logs, (c) search all dependent codebases. Do not proceed to Phase 3 until verified.
6. **Phase 3 — Contract: remove the old** — Remove the old API element after all consumers have migrated. For fields: drop the old column (after verifying zero reads/writes). For methods: remove the old overload. For schemas: run the column drop migration. Confirm with tests that the old element no longer exists.
7. **Communicate the timeline to consumers** — For externally visible APIs, the migration timeline must be communicated: announce the deprecation at expansion deployment, set a contractual sunset date (typically 6–12 months for external APIs), and provide migration documentation. Do not skip communication even for internal APIs with multiple teams.
8. **Validate end-to-end after contraction** — After removing the old element, run the full integration test suite. Confirm that the system functions correctly with only the new element active. This is the final validation that migration is complete.

## Rules

- Never deploy Phase 2 (consumer migration) before Phase 1 (expansion) is deployed and verified — consumers cannot migrate to a new API element that does not exist yet.
- Never skip directly from adding the new to removing the old in a single deployment — the combined deploy has the same risk as a coordinated cutover; the phases must be separate deployments.
- The old element must remain fully functional during the entire migration window — do not degrade the old element's behavior to "encourage" migration; this creates hidden failures.
- For database columns: the write path must populate both old and new columns during the migration period — if only the new column is written, rollback to the old code path produces data loss.

## Common Mistakes

- **Renaming a database column in a single migration** — Dropping the old column in the same migration that adds the new column causes a deployment window where the application expects one schema and the database provides another. Always separate add → migrate → drop into three deployments.
- **Removing the old element before all consumers migrate** — Removing prematurely breaks any consumer that has not yet migrated. Instrument the old element with usage logging and set a zero-usage threshold before removing.
- **Using Parallel Change for additive changes** — Adding a new optional field to a REST response or adding a new method to a service does not break existing consumers; it does not require Parallel Change. Apply the pattern only to breaking changes.
- **Skipping the contraction phase** — Leaving the old element indefinitely after all consumers migrate creates dead code that future developers will be confused by, and in the case of database columns, ongoing storage and index overhead. Schedule contraction immediately after migration confirms zero usage.

## Examples

**REST API field rename:** Rename `user.name` to `user.full_name`. Phase 1: respond with both `name` and `full_name` (same value). Phase 2: update all API clients to use `full_name` over 4 weeks. Phase 3: remove `name` from responses, update API schema documentation.

**Database column rename:** Rename `orders.client_id` to `orders.customer_id`. Phase 1: add `customer_id` column; code writes both columns; reads from `customer_id`. Phase 2: migrate all queries to read from `customer_id`; verify zero reads from `client_id` in logs. Phase 3: drop `client_id` column in next release cycle.

**Method signature change:** Change `processOrder(orderId: string)` to `processOrder(order: OrderRequest)`. Phase 1: add new overload, keep old signature. Phase 2: update all callers to new overload. Phase 3: remove old signature.

## When NOT to Use

- When the change is purely additive (new field, new endpoint, new optional parameter) — additive changes are backward-compatible by definition and do not require Parallel Change.
- When you control all consumers and can deploy them atomically with the provider — if the API and all consumers ship in a single atomic deployment with zero downtime, a coordinated cutover is safe and simpler.
- When the migration window is so short (< 1 hour) that an expand-and-contract cycle would take longer than just coordinating a cutover — for very small, internal, fast-moving codebases, the pattern overhead exceeds the benefit.
