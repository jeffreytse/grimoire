---
name: design-system-architecture
description: Use when asked to design, evaluate, or document the high-level architecture of a system — including new systems, major feature additions, or scaling a system beyond its current design limits.
source: Amazon/Google SWE design practices; Kleppmann "Designing Data-Intensive Applications" (O'Reilly 2017); Google SRE Book (O'Reilly 2016)
tags: [system-design, architecture, scalability, trade-offs, distributed-systems, technical-leadership]
verified: true
---

# Design System Architecture

Derive requirements, identify components, define data flow, and articulate trade-offs — before writing any code.

## Why This Is Best Practice

**Adopted by:** Amazon (design documents required before any significant engineering work per Amazon's engineering tenets), Google (design docs are mandatory for projects above a size threshold — described in "Software Engineering at Google", O'Reilly 2020), Stripe (RFC process for architectural decisions).

**Impact:** IBM Systems Sciences Institute found that fixing a defect found in design costs 6× less than one found during implementation, and 15× less than one found post-release. Google's internal data (Winters et al., "Software Engineering at Google") shows that teams with upfront design docs ship features with 50% fewer post-launch incidents. The cost of a design document is 1–2 days; the cost of a wrong architecture is months of rework.

**Why best:** Writing code without a design is the equivalent of building a house without blueprints. Informal design (whiteboard-only) produces undocumented trade-offs that trap future engineers. Formal upfront design (big design up front, BDUF) goes too far and creates documents no one reads. The sweet spot is a lightweight, decision-focused design doc that documents *why*, not *how*.

Sources: Kleppmann "DDIA" (O'Reilly 2017), Winters et al. "Software Engineering at Google" (O'Reilly 2020), Google SRE Book (Beyer et al., O'Reilly 2016), IBM Systems Sciences Institute defect cost study

## Steps

### 1. Define functional requirements

List what the system must do. Write these as user-facing capabilities, not technical choices:
- "Users can upload files up to 5 GB"
- "The system sends email notifications within 60 seconds of an event"
- "Search returns results in under 500 ms at p99"

Distinguish must-have from nice-to-have. Scope creep in requirements is the single most common source of over-engineered systems.

### 2. Define non-functional requirements (SLOs)

Quantify quality attributes. Vague targets ("fast", "reliable") are undesignable. Require numbers:

| Attribute | Metric | Target |
|-----------|--------|--------|
| Latency | p99 response time | < 200 ms |
| Availability | Uptime per month | 99.9% (43 min downtime) |
| Throughput | Peak requests/sec | 10,000 RPS |
| Durability | Data loss tolerance | Zero (RPO = 0) |
| Storage | Data volume in 3 years | ~10 TB |

Work backwards from these numbers to size components.

### 3. Estimate scale (back-of-envelope)

Before drawing components, sanity-check the numbers:
- Daily active users × average requests per user = daily request volume
- Request volume / 86,400 seconds = average RPS; multiply by peak factor (typically 3–5×) = peak RPS
- Data written per request × daily volume × retention period = storage requirement
- These estimates expose whether you need a single database or a distributed system, a monolith or microservices

### 4. Identify core components

Map the functional requirements to components. Common building blocks:
- **Load balancer** — distribute traffic, terminate TLS
- **API gateway / edge** — auth, rate limiting, routing
- **Application servers** — stateless computation
- **Cache** (Redis, Memcached) — reduce read latency, absorb spikes
- **Message queue** (Kafka, SQS) — decouple producers from consumers, buffer bursts
- **Primary database** — source of truth (choose type: relational, document, wide-column, time-series)
- **Search index** (Elasticsearch, Typesense) — full-text and faceted search
- **Object store** (S3, GCS) — large binary data (files, images, video)
- **CDN** — static assets and cacheable responses at the edge

Only include components that the requirements justify. Every component adds operational complexity.

### 5. Design data flow

For each core user journey (typically 3–5), draw the request path end-to-end:
1. Client → edge → application → storage: the write path
2. Client → edge → cache → database → client: the read path
3. Event → queue → consumer → derived store: the async processing path

For each path, identify: latency budget, failure modes, and consistency requirements (is stale data acceptable?).

### 6. Choose the database(s)

The database choice is the highest-leverage and hardest-to-change architectural decision. Apply these heuristics:

| Need | Consider |
|------|----------|
| Relational data with transactions | PostgreSQL, MySQL |
| High-write time-series data | ClickHouse, TimescaleDB, InfluxDB |
| Flexible schema, document-oriented | MongoDB, DynamoDB |
| Global distribution, multi-region writes | CockroachDB, Spanner, DynamoDB Global Tables |
| Graph traversals | Neo4j, Amazon Neptune |
| Full-text search | Elasticsearch, OpenSearch |

Avoid polyglot persistence unless the requirements force it — each additional database is a synchronization problem.

### 7. Address the hard problems explicitly

Every non-trivial system has one or two hard problems. Name them and state your approach:

- **Consistency vs. availability**: Under a network partition, does the system serve stale data or return an error? (CAP theorem — pick a side per use case)
- **Hot spots**: Will certain keys (celebrity users, viral content) get disproportionate traffic? Pre-shard or use random jitter
- **Fan-out**: If one event must notify 1 million subscribers, push (write to each) or pull (subscribers poll) — Twitter uses a hybrid based on follower count
- **Data at scale**: At what volume does a single database node become a bottleneck? Plan horizontal sharding before you need it, not after

### 8. Document trade-offs

For each significant architectural choice, record: the option chosen, the alternatives considered, and the reason for the choice. This is the most valuable part of the design document — it prevents future engineers from re-litigating settled decisions.

Format:
> **Decision**: Use PostgreSQL as the primary store.
> **Alternatives**: MongoDB (rejected — relational joins needed), DynamoDB (rejected — complex queries not supported).
> **Reason**: Data is relational; referential integrity is required; team has PostgreSQL expertise.

## Rules

- No component enters the design without a justified requirement — "we might need it later" is not a requirement
- Every SLO must have a measurement plan — if you cannot measure it, you cannot enforce it
- State failure modes for every external dependency — what happens when the cache is down, when the queue is full?
- The design document is a decision log, not a tutorial — future readers want *why*, not *how*
- For systems under 1,000 RPS, prefer a monolith — microservices add operational complexity that small teams cannot absorb

## Examples

**Back-of-envelope for a photo sharing app:**
- 10M DAU, 2 uploads/day = 20M uploads/day = 230 uploads/sec average, ~1,000/sec peak
- Average photo 3 MB compressed = 60 TB/day new storage; at 90-day retention = ~5 PB total
- Conclusion: object store required (S3), CDN essential, database stores metadata only

**Trade-off example:**
> We chose eventual consistency for the "likes" counter. Exact counts are not user-critical; a 1–5 second lag is acceptable. Strong consistency would require distributed coordination that adds 50–100 ms latency per write — not worth it for a display-only metric.

## Common Mistakes

- Designing microservices before validating the monolith is too slow — premature decomposition
- Picking a database based on hype rather than requirements — Redis is not a primary store
- No numbers in non-functional requirements — "high availability" is not a target
- Designing for the peak of the peak (Black Friday × 10) — over-provisions by 10× for normal operation
- Skipping failure modes — the design looks perfect until a dependency goes down
