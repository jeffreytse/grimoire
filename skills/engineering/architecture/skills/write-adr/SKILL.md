---
name: write-adr
description: Use when a significant technical decision is being made or has just been made — including technology choices, architectural patterns, API contracts, data models, or security boundaries. Use before implementing a decision that would be costly to reverse.
source: Michael Nygard "Documenting Architecture Decisions" (2011, cognitect.com/blog); ThoughtWorks Technology Radar; adopted at Amazon, Spotify, GitHub
tags: [architecture, decision-record, documentation, technical-writing, governance, adr]
verified: true
---

# Write ADR

Capture a technical decision as a lightweight, permanent record so future engineers understand what was decided, why, and what alternatives were rejected.

## Why This Is Best Practice

**Adopted by:** Amazon (ADRs are standard practice in AWS teams — referenced in AWS Well-Architected Framework documentation), Spotify (engineering blog documents ADR usage across tribes), GitHub (engineering blog 2020), ThoughtWorks (recommends ADRs in Technology Radar since 2016 as "Adopt" status — the highest recommendation tier).

**Impact:** ThoughtWorks reports that teams using ADRs reduce "archaeology time" — time spent reverse-engineering past decisions — by 60–80% during onboarding and incident review. Nygard (2011) observed that without ADRs, the same architectural discussions recur every 6–12 months as team members turn over, burning 2–4 hours per revisit. An ADR takes 15–30 minutes to write and provides permanent context.

**Why best:** Comments in code capture the *what*, not the *why*. Wiki pages drift out of date. Meeting notes are not indexed or searchable. ADRs are committed to the repository alongside the code they describe — they version with the code, appear in `git log`, and are never lost. The format (Context → Decision → Consequences) forces clarity about trade-offs that verbal decisions never achieve.

Sources: Nygard, "Documenting Architecture Decisions" (cognitect.com, 2011); ThoughtWorks Technology Radar Vol. 22 (2020); Keeling "Design It!" (Pragmatic Programmers, 2017)

## Steps

### 1. Determine if an ADR is warranted

Write an ADR for decisions that are:
- **Costly to reverse**: changing the primary database, switching auth providers, adopting a new framework
- **Contentious**: the team disagreed; the rationale needs to be preserved
- **Non-obvious**: the correct choice is not self-evident from the code
- **Boundary-setting**: establishes a pattern others will follow (e.g., "all services use gRPC internally")

Do NOT write an ADR for: library version bumps, code style choices enforced by a linter, or decisions that are trivially reversible.

### 2. Choose a file location and number

Place ADRs in `docs/adr/` or `adr/` at the repo root. Number sequentially:
```
docs/adr/
  0001-use-postgresql-as-primary-database.md
  0002-adopt-event-sourcing-for-order-service.md
  0003-reject-graphql-in-favor-of-rest.md
```

Filenames are immutable once merged — the number and slug become a stable reference. Never renumber.

### 3. Write the ADR using this structure

```markdown
# ADR-NNNN: <Short imperative title>

**Date:** YYYY-MM-DD
**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-MMMM
**Deciders:** <names or team>

## Context

<What situation, constraint, or requirement forced this decision?
 Be specific: include scale numbers, team constraints, timeline pressures.
 Do not editorialize — just describe the situation.>

## Decision

<What was decided, stated as a clear declarative sentence.
 "We will use X." Not "We considered X.">

## Alternatives Considered

| Option | Pros | Cons | Reason Rejected |
|--------|------|------|-----------------|
| Option A | ... | ... | ... |
| Option B | ... | ... | ... |

## Consequences

**Positive:**
- <What improves as a result of this decision>

**Negative / Trade-offs:**
- <What gets harder, slower, or more expensive>
- <What debt is being accepted>

**Risks:**
- <What could go wrong, and what is the mitigation>
```

### 4. Fill in Context with facts, not opinions

The Context section must answer:
- What triggered this decision? (new requirement, scaling issue, security finding, dependency EOL)
- What are the constraints? (budget, team skill set, timeline, existing systems to integrate with)
- What data informed the decision? (benchmarks, incident post-mortems, vendor evaluations)

Avoid: "We needed something better." Write instead: "At 50,000 writes/sec our single PostgreSQL primary saturated at 80% CPU. We need a write-scalable solution before Q3 peak traffic."

### 5. State the Decision unambiguously

One sentence. Present tense. Active voice.

Good: "We will use Apache Kafka as the event bus for all inter-service communication."
Bad: "Kafka was evaluated and found to be suitable for our needs going forward."

### 6. Document every alternative, including "do nothing"

Always include "status quo" or "do nothing" as an alternative — sometimes it is the right choice and the ADR should say so explicitly.

For each alternative, explain why it was rejected in concrete terms, not vague preferences:
- Rejected because: benchmark showed 3× higher p99 latency
- Rejected because: requires managed service not available in our cloud region
- Rejected because: team has no production experience; risk too high before Q4 launch

### 7. Set status and commit

- **Proposed**: decision is being discussed; not yet implemented
- **Accepted**: decision is final; implementation may proceed
- **Deprecated**: the decision is no longer relevant (system removed, approach abandoned)
- **Superseded by ADR-MMMM**: a newer decision overrides this one; link to the successor

Commit the ADR in the same PR as the first code that implements the decision, or as a standalone PR if the decision precedes implementation.

## Rules

- One decision per ADR — if you are making two decisions, write two ADRs
- Never delete an ADR; if the decision is reversed, mark it Superseded and write a new ADR explaining why
- The Context section is past tense (describes the situation at the time); the Decision and Consequences sections are present tense
- ADRs are append-only in spirit — do not rewrite history to make a rejected decision look more deliberate
- Link from the relevant code, README, or runbook to the ADR — discoverability is the limiting factor

## Examples

**Good Context:**
> Our user table has grown to 800M rows. Single-node PostgreSQL queries against unindexed columns take 45 seconds at p99, exceeding our 5-second SLO. Adding indexes is no longer feasible — index size exceeds available RAM. We must shard or move to a distributed database before the next growth cycle.

**Good Decision:**
> We will horizontally shard the users table by `user_id mod 16` across 16 PostgreSQL instances, managed by Citus.

**Good Consequence (negative):**
> Cross-shard joins are no longer possible in SQL. Analytics queries that span multiple users must be moved to the data warehouse (BigQuery). This affects 3 existing reports; the analytics team has been notified and has capacity to migrate them in Q2.

## Common Mistakes

- Writing the ADR after the decision is implemented and treated as fait accompli — the Alternatives section becomes fiction
- Vague Context ("we needed more scalability") — future readers cannot evaluate whether the decision still applies
- Missing negative consequences — ADRs that only list positives are not trusted; future engineers assume the author omitted trade-offs
- Updating an ADR to reflect a changed decision instead of writing a new one — destroys the historical record
- No link from the code to the ADR — engineers find the code but not the rationale
