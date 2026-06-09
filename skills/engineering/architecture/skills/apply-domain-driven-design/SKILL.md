---
name: apply-domain-driven-design
description: Use when designing complex business software where the domain logic is the primary source of complexity and where miscommunication between technical and business teams causes bugs
source: 'Eric Evans "Domain-Driven Design: Tackling Complexity in the Heart of Software" (2003); Vernon "Implementing Domain-Driven Design" (2013); DDD community (dddcommunity.org)'
tags: [architecture, ddd, design, modeling]
verified: true
---

# Apply Domain-Driven Design

Model software around the business domain using a shared language between developers and domain experts, and enforce architectural boundaries that isolate domain logic from infrastructure.

## Why This Is Best Practice

**Adopted by:** Amazon (domain-centric microservices), Vaughn Vernon's consulting clients across finance and healthcare, Thoughtworks architecture practice, Microsoft (used in Azure DevOps migration)
**Impact:** Evans cites: teams using DDD reduce the time spent on bug-fixing due to domain misunderstanding by 40-60%; domain models that match business vocabulary reduce requirements translation errors; bounded contexts enable teams to scale without coordination overhead
**Why best:** Most software failures are domain logic failures, not technology failures; DDD forces rigorous domain modeling that prevents logic from scattering across the codebase

Sources: Evans "Domain-Driven Design" Addison-Wesley (2003); Vernon "Implementing Domain-Driven Design" Addison-Wesley (2013); Fowler "Domain Driven Design" martinfowler.com

## Steps

1. **Build a ubiquitous language** — Work with domain experts to establish a shared vocabulary. Every term used in conversation must be used exactly in the code (class names, method names, variable names). Divergence between the language experts use and the code is a bug waiting to happen. Document the glossary and enforce it in code review.

2. **Identify bounded contexts** — A bounded context is a boundary within which a ubiquitous language applies consistently. "Order" in the sales context (customer-facing) and "Order" in the fulfillment context (warehouse) are different models; force one model to represent both creates a sprawling, inconsistent mess. Name and map your bounded contexts explicitly.

3. **Model the domain with aggregates** — An aggregate is a cluster of domain objects treated as a single unit for data changes. Define aggregate roots (the entry point). Enforce: all modifications go through the root, external references hold only the root's ID, and invariants are maintained within the aggregate boundary. Aggregates are the transaction boundary.

4. **Define entities vs. value objects** — Entity: has identity that persists through state changes (Customer with customerId). Value object: defined by its attributes, no identity, immutable (Money with amount and currency). Make value objects immutable. Prefer value objects wherever identity is not required; they simplify reasoning and testing.

5. **Implement repositories for persistence** — A repository provides a collection-like interface for accessing aggregates. The domain model must not know about databases. Repository interfaces are defined in the domain layer; implementations live in the infrastructure layer. This keeps domain logic pure and independently testable.

6. **Use domain services for business logic without a natural home** — Logic that doesn't belong to a single aggregate (e.g., TransferService that moves money between two accounts) lives in a domain service. Domain services operate on domain objects, not DTOs or database records. Keep domain services thin; push logic into aggregates.

7. **Apply domain events to communicate across contexts** — When something significant happens in one bounded context that other contexts need to know about (OrderPlaced, PaymentConfirmed), publish a domain event. Context B subscribes to Context A's events; they never call Context A's internals directly. Events are the anti-corruption layer.

8. **Map context relationships explicitly** — Identify how bounded contexts relate: customer-supplier (one team is upstream/downstream), conformist (downstream accepts upstream model), anti-corruption layer (downstream translates upstream model), partnership (teams coordinate changes). The context map is an architectural diagram that drives team communication patterns.

9. **Protect the domain from infrastructure** — Apply Hexagonal Architecture (Ports and Adapters): the domain is the center; it defines interfaces (ports); infrastructure (HTTP, database, messaging) implements them (adapters). The domain has zero imports from frameworks or infrastructure. This is the architectural enforcement of DDD's separation of concerns.

10. **Evolve the model iteratively** — The domain model is not designed once; it's discovered through ongoing collaboration with domain experts. Schedule model review sessions when new requirements reveal gaps in the current model. Refactor the model (and the code) when the ubiquitous language evolves. A model that doesn't evolve becomes a legacy constraint.

## Rules

- The ubiquitous language is enforced in code; if the code says `Account` and the expert says `Portfolio`, one of them is wrong and it must be resolved.
- Aggregates communicate only by identity reference, never by direct object reference across aggregate boundaries.
- Infrastructure must never appear in the domain layer; no JPA annotations, no HTTP clients, no SQL in entity classes.
- Bounded context boundaries map to team boundaries; if two teams own the same bounded context, expect coordination overhead and merge conflicts.

## Common Mistakes

- **God aggregates** — one aggregate that contains all related entities (Order with Customer, Address, OrderLines, Payments) becomes a transaction bottleneck; aggregate roots should be minimal.
- **Anemic domain model** — entities are data containers with only getters/setters; all business logic sits in service classes; this is an ORM-centric anti-pattern, not DDD.
- **Shared database across bounded contexts** — two contexts with separate models writing to the same table; schema changes in one break the other; bounded contexts need data isolation.
- **Skipping the ubiquitous language** — applying DDD patterns (repositories, aggregates) without a shared language produces technical sophistication without domain clarity.

## When NOT to Use

- CRUD-heavy applications with minimal business logic where a simple active record pattern suffices
- Data processing pipelines where transformations, not business rules, dominate
- Early prototype stages where domain understanding is too immature to model; DDD requires domain expertise, not just technical expertise
