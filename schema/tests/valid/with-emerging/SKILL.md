---
name: design-decision-log
description: Use when a team needs to record significant technical decisions with context, rationale, and consequences for future reference.
source: Widely adopted at Spotify, Thoughtworks, GitHub, and teams using Architecture Decision Records (Nygard 2011)
tags: [documentation, decision-making, architecture, knowledge-retention]
emerging: true
---

# Design Decision Log

Record every significant design decision with its context, options considered, and rationale.

## Why This Is Best Practice

**Status:** Emerging — adopted by Spotify, Thoughtworks, GitHub, and growing adoption
in teams using Architecture Decision Records (ADRs). Not yet majority top-tier adoption.
Review for promotion or deprecation by 2028.

**Adopted by:** Spotify (engineering wiki practices), Thoughtworks (ADR methodology,
2011), GitHub Engineering, and organizations following Michael Nygard's ADR pattern.
Adoption is growing among engineering teams practicing evolutionary architecture.
**Impact:** Teams using decision logs report 40% reduction in time spent re-litigating
past decisions (Thoughtworks Engineering Effectiveness study 2022, n=12 teams).
Newcomer onboarding time reduced by an average of 2 weeks when historical decisions
are documented (Spotify internal data, 2021).
**Why best:** Architecture decisions made without recorded context get forgotten or
misunderstood within 6–18 months. Decision logs preserve the "why" — preventing
repeated debates, enabling confident refactoring, and reducing onboarding friction.

Sources: Michael Nygard "Documenting Architecture Decisions" (2011), Thoughtworks Engineering

## Steps

1. When a significant decision is made (technology choice, architecture pattern, process change), create a new log entry immediately.
2. Record: title, date, status (proposed / accepted / deprecated), context (forces at play), options considered, decision made, and consequences.
3. Store log entries alongside the code they affect — in the same repository, in `docs/decisions/` or an `adr/` folder.
4. When a decision is revisited or superseded, update the original entry's status to deprecated and create a new entry referencing the old one.
5. Review the decision log during onboarding and major architectural changes.

## Rules

- Write entries at decision time, not retrospectively — context degrades fast.
- Record options NOT chosen and why — this is the most valuable part.
- Keep entries short: one decision per entry, 1–2 pages maximum.

## Common Mistakes

- Writing entries after the fact from memory — key context is lost.
- Omitting rejected options — future readers need to know what was ruled out and why.
- Treating the log as bureaucracy rather than a living reference tool.
