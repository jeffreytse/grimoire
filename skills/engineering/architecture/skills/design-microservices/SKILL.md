---
name: design-microservices
description: Use when decomposing a monolith or designing a new system using microservices architecture
source: Sam Newman "Building Microservices" (2nd ed., O'Reilly 2021); Eric Evans "Domain-Driven Design" (Addison-Wesley 2003)
tags: [microservices, architecture, ddd, bounded-context, decomposition, distributed-systems]
verified: true
---

# Design Microservices

Decompose a system into independently deployable services aligned to business capabilities using bounded contexts.

## Why This Is Best Practice

**Adopted by:** Netflix, Amazon, Uber, Spotify — all publicly documented their migrations from monoliths
**Impact:** Amazon's service decomposition enabled teams to deploy independently, reducing release cycle from months to minutes; Netflix processes 2+ billion API requests/day across ~700 microservices.

**Why best:** DDD bounded contexts provide the natural seam for service decomposition — they align services to business domains, minimize cross-service coupling, and give each team a clear ownership boundary. Sam Newman's patterns (strangler fig, anti-corruption layer) provide proven migration paths.

## Steps

1. **Map business capabilities** — Identify the core domains (e.g., Orders, Inventory, Payments) using event storming or domain mapping workshops.
2. **Define bounded contexts** — Draw context boundaries where the ubiquitous language changes; each context is a candidate service.
3. **Apply the right decomposition pattern** — Decompose by business capability (preferred) or subdomain; avoid decomposing by technical layer (e.g., "database service").
4. **Define service contracts** — Specify each service's public API (REST, gRPC, or events) before implementation; treat contracts as stable interfaces.
5. **Design for failure** — Every inter-service call can fail; apply timeouts, retries with backoff, and circuit breakers (Hystrix/Resilience4j pattern).
6. **Choose a data strategy** — Each service owns its data store; share data via events or APIs, never via shared databases.
7. **Plan observability** — Agree on distributed tracing (OpenTelemetry), structured logging, and health endpoints before first service ships.

## Rules

- One database per service — shared databases couple services and eliminate independent deployability.
- Services must be independently deployable without coordinating with other teams.
- Avoid nano-services — if a service only exists to wrap a CRUD table, it belongs in a larger service.
- Prefer choreography (events) over orchestration for business workflows to avoid a central point of failure.

## Examples

E-commerce decomposition:
- `order-service` — creates and tracks orders; publishes `OrderPlaced` events.
- `inventory-service` — listens for `OrderPlaced`, reserves stock, publishes `StockReserved`.
- `payment-service` — listens for `StockReserved`, charges customer.

Each service has its own DB; communication is via Kafka events; no service calls another's DB directly.

## Common Mistakes

- **Distributed monolith** — services that must deploy together or share a database gain none of the benefits of microservices.
- **Decomposing too early** — start with a well-structured monolith ("modular monolith") until domain boundaries are understood.
- **No contract versioning** — breaking API changes cascade across services; version APIs from day one.

## When NOT to Use

- When the team has fewer than 8 engineers, the operational overhead of independent deployments, distributed tracing, and per-service databases will consume more capacity than the autonomy gains.
- When the domain boundaries are not yet understood — building microservices before the domain model stabilizes results in services that must be merged or split within months, producing a distributed monolith with extra steps.
- When the system has no requirement for independent scaling or independent team ownership, a well-structured modular monolith delivers the same maintainability at a fraction of the operational complexity.
