---
name: design-cqrs-architecture
description: Use when designing a system where read and write workloads have significantly different scaling, consistency, or model requirements that a single unified model cannot serve efficiently.
source: Young "CQRS Documents" (2010); Fowler "CQRS" (martinfowler.com, 2011); Microsoft Azure Architecture Guide CQRS pattern; Vernon "Implementing Domain-Driven Design" (2013)
tags: [architecture, cqrs, patterns, scalability, design, distributed-systems, ddd, event-sourcing]
verified: true
---

# Design CQRS Architecture

Separate the system's command (write) model from its query (read) model so each can be optimized, scaled, and evolved independently — without forcing a single model to serve contradictory requirements.

## Why This Is Best Practice

**Adopted by:** Microsoft (Azure Architecture Center recommended pattern), Axon Framework (Java CQRS/ES framework, 10,000+ enterprise users), Netflix, LinkedIn (feed architecture), and widely in financial services and e-commerce systems with high read/write asymmetry. Greg Young coined CQRS; Martin Fowler documented and normalized it as an architectural pattern.
**Impact:** Systems with high read/write asymmetry (typical ratio: 80–95% reads, 5–20% writes) achieve 3–10× read throughput improvement by separating read stores optimized for query patterns from normalized write stores. LinkedIn's feed system migration to a CQRS-derived architecture handled 5× traffic growth without proportional infrastructure scaling (LinkedIn Engineering 2012). Microsoft reports CQRS enables independent scaling of read and write infrastructure, reducing over-provisioning costs by 40–60% in asymmetric workloads.
**Why best:** A unified data model optimized for write integrity (normalized, transactional) is structurally poor for read performance (denormalized, cached, pre-aggregated). Trying to optimize one model for both creates performance ceilings and forces compromises on both consistency and query flexibility. CQRS allows the write side to enforce business invariants cleanly, while the read side exposes optimized, denormalized projections tailored to each consumer's needs.

Sources: Young "CQRS Documents" (cqrs.files.wordpress.com 2010); Fowler "CQRS" (martinfowler.com/bliki/CQRS.html 2011); Microsoft Azure Architecture Patterns: CQRS; Vernon "Implementing Domain-Driven Design" (2013) Part 4; Richardson "Microservices Patterns" (2018) Ch. 7

## Steps

1. **Assess whether CQRS is warranted** — CQRS adds complexity; apply it only when one or more is true: (a) read/write workloads have different scaling needs, (b) read models need denormalized projections the write model cannot efficiently serve, (c) the system uses Event Sourcing as the write store, (d) different teams own read and write paths. Do not apply CQRS to simple CRUD systems.
2. **Define the command model (write side)** — Design the write model around business commands and invariants, not UI or query needs. Commands represent intent ("PlaceOrder", "CancelSubscription"). The command handler validates the command, applies domain logic, enforces invariants, and persists state changes. The write model should be normalized and optimized for consistency, not for querying.
3. **Choose the write store** — Options: relational database (normalized schema, ACID), Event Store (append-only event log — required for Event Sourcing), or document store. The write store is the system of record; choose based on consistency and transaction requirements, not query performance.
4. **Define query models (read side)** — For each consumer or use case (API response, dashboard, report), design a read model that directly matches what the consumer needs. Read models are denormalized, pre-aggregated, and optimized for query performance. One system can have multiple read models serving different consumers.
5. **Choose the read store(s)** — Read stores are optimized for the read model's access patterns: relational with read replicas, Redis (key-value lookup), Elasticsearch (full-text search), Cassandra (time-series), or materialized views. The read store is derived from the write store — it can be rebuilt.
6. **Design the synchronization mechanism** — The read store must be kept in sync with the write store. Options: (a) synchronous update in the same transaction (simple, limits scalability), (b) domain events published after each write and consumed by read model projectors (eventual consistency, recommended for scale), (c) database change data capture (CDC) feeding read model updates. Define acceptable staleness for each read model.
7. **Implement command validation at the boundary** — Commands are validated at entry: structural validation (request shape), authorization (permission to issue this command), and domain validation (business rule check using the write model). Invalid commands are rejected with clear error codes before any state change.
8. **Handle eventual consistency explicitly** — If read and write models synchronize asynchronously, the system exhibits eventual consistency — a read immediately after a write may not reflect the change. Decide per use case: (a) return the command result directly from the write side, (b) poll until consistent, (c) accept that the UI will show stale data with a refresh. Never hide eventual consistency from users without a compensating UX pattern.

## Rules

- Commands mutate state; queries must not — a query handler that triggers side effects (writes, emails, counters) breaks CQRS and makes the read side unpredictable.
- The read model is disposable — the read store must be rebuildable from the write store at any time; if it cannot be rebuilt, the write store is not the system of record.
- Do not share a single ORM model between command and query sides — sharing models couples them and defeats the purpose of separation; define explicit command DTOs and query DTOs.
- Apply CQRS at the service or bounded context level, not application-wide — not every aggregate or service needs CQRS; apply it only where read/write asymmetry justifies the complexity.

## Common Mistakes

- **Applying CQRS everywhere** — Fowler explicitly warns that CQRS should be applied selectively; most systems do not have sufficient read/write asymmetry to justify it. Applying it uniformly adds complexity with no benefit in simple CRUD contexts.
- **Synchronizing read models in-transaction** — Updating both write and read stores in the same transaction creates tight coupling and eliminates the independent scaling benefit. Use events or CDC for asynchronous projection.
- **Exposing the write model to the read side** — Returning the domain aggregate from a command handler as the query response uses the wrong model for the wrong purpose. Define explicit read DTOs from the read store.
- **Ignoring eventual consistency in the UI** — Showing a user a "success" confirmation and immediately reloading a list that hasn't yet reflected their action destroys trust. Design UX patterns (optimistic UI, polling, websockets) to handle the consistency window.

## Examples

**E-commerce order system:** Write side: normalized PostgreSQL schema with ACID transactions for order creation, inventory reservation, payment. Read side: Redis cache for order status (ms lookup), Elasticsearch for order history search, denormalized PostgreSQL view for order management dashboard. Events synchronize via Kafka.

**Financial portfolio system:** Write side: event store recording every trade, deposit, and adjustment as immutable events. Read side: materialized PostgreSQL view of current positions (rebuilt nightly), Redis for real-time P&L display, Elasticsearch for transaction history. Event sourcing + CQRS enables complete audit trail and point-in-time portfolio reconstruction.

## When NOT to Use

- When the system is primarily CRUD with no significant read/write asymmetry — the added complexity of separate models, synchronization, and eventual consistency management exceeds any benefit.
- When the team is new to distributed systems — CQRS requires operational maturity to manage synchronization failures, consistency windows, and read model rebuilds; ensure the team can operate the pattern before adopting it.
- When strong consistency on every read is non-negotiable — financial regulatory reporting, medical record systems, and some billing scenarios may not tolerate the eventual consistency window that asynchronous synchronization introduces.
