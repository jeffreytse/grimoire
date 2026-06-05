---
name: write-tech-spec
description: Use when a non-trivial engineering change needs to be documented before implementation — including new services, significant refactors, API design, database schema changes, cross-team dependencies, or any work estimated above 1 week. Trigger before the sprint that contains implementation work.
source: Stripe engineering spec process (described in Stripe Press engineering culture), Notion design doc format (Notion Engineering Blog), Google technical writing standards (developers.google.com/tech-writing)
tags: [technical-spec, design-doc, architecture, engineering, alignment, api-design, pre-implementation]
verified: true
---

# Write Tech Spec

Write a technical specification that aligns engineers on what to build and how before implementation begins — capturing the problem, proposed solution, alternatives considered, and open questions.

## Why This Is Best Practice

**Adopted by:** Stripe (engineering blog documents mandatory design docs for services), Google (design doc culture is central to the engineering process, documented in "Software Engineering at Google"), Notion, Airbnb, Dropbox, and most engineering organizations above ~20 engineers.

**Impact:** Google's "Software Engineering at Google" (O'Reilly, 2020) reports that design docs catch ~50% of design flaws before any code is written. A spec review that takes 2 engineer-hours prevents weeks of rework. Studies cited in the book show that the cost to fix a bug found during design is 10× cheaper than finding it during testing and 100× cheaper than finding it in production.

**Why best:** Technical specs externalize design reasoning. Without them, critical decisions exist only in one engineer's head, making onboarding harder, debugging slower, and post-mortems incomplete. The discipline of writing forces the author to think through edge cases, failure modes, and alternatives they might skip when coding. Specs also create a historical record of why decisions were made — invaluable when those decisions are questioned six months later.

Sources: Titus Winters, Tom Manshreck, Hyrum Wright, "Software Engineering at Google" (O'Reilly, 2020, Chapter 10); Stripe Engineering Blog; Google Technical Writing Course (developers.google.com/tech-writing); Gergely Orosz, "The Software Engineer's Guidebook" (2023).

## Steps

1. **Write the problem statement** (2–4 sentences): what technical problem are you solving? Why does it exist now? What breaks if you don't solve it? Include relevant metrics or incident data if available.

2. **State the goals and non-goals**: goals are the specific outcomes the solution must achieve. Non-goals are explicitly out of scope. Both must be measurable or clearly bounded. Example goal: "Reduce p99 latency of the search API from 2.1s to under 500ms." Example non-goal: "This spec does not address search ranking — that is handled separately."

3. **Describe the proposed solution** in enough detail that a senior engineer who wasn't in any of the planning meetings could implement it. Include:
   - System/component diagram if the change spans multiple services
   - Data model changes (schema diffs, new fields, removed fields)
   - API contract (endpoint paths, request/response shapes, error codes)
   - Key algorithms or data structures
   - Dependencies on other teams or services

4. **List alternatives considered** — at least two. For each alternative: what it is, why you considered it, and why you rejected it. This section is not optional. Without it, reviewers will suggest alternatives you already evaluated, and future engineers won't know why the current approach was chosen.

5. **Address operational concerns**:
   - How will you deploy this safely? (feature flags, gradual rollout, dark launch)
   - How will you monitor it? (metrics, alerts, dashboards)
   - What does rollback look like?
   - What are the failure modes, and what is the blast radius of each?

6. **Write the open questions section** — anything unresolved that needs a decision before or during implementation. Assign each question to an owner with a target resolution date.

7. **Estimate scope**: rough implementation timeline by phase, and which teams or engineers are needed. This is for planning purposes, not a commitment.

8. **Get async written review** from: the tech lead or staff engineer for the affected system, engineers who will implement the work, and any teams that own dependencies. Collect feedback in the document itself, not in Slack threads.

9. **Resolve open questions and update the spec** before engineering kicks off. The spec should reflect the final agreed design, not just the initial proposal.

## Rules

- The problem statement must appear before the solution. If you lead with the solution, you prevent reviewers from questioning whether the solution matches the problem.
- Alternatives considered is a required section. A spec with only one option is a proposal, not a design. Reviewers will always ask "did you consider X?" — answer it in advance.
- Diagrams must be editable source (Mermaid, draw.io XML, Excalidraw) not image screenshots. Diagrams that can't be updated become stale and misleading.
- API contracts must include error responses, not just happy paths. Every error code the API returns must be specified.
- The spec is frozen when implementation begins. After that, changes are tracked as amendments with a date and author, not silent edits.
- Length scales with scope: a 2-day refactor warrants 1–2 pages. A new microservice warrants 4–8 pages. A platform migration may warrant 10–15 pages with appendices.

## Examples

**Problem statement (bad):** "We need to refactor the payments service to be more scalable."

**Problem statement (good):** "The payments service processes all transactions synchronously in a single process. At current growth (20% MoM), we will exceed our single-process throughput ceiling within 60 days. The P0 incident on 2024-01-15 (COE-447) was caused by exactly this bottleneck. This spec proposes async processing via a job queue to decouple ingestion from processing."

**Alternatives considered (bad):** "We considered other approaches but this one was best."

**Alternatives considered (good):**
```
Alternative 1: Vertical scaling (larger instance type)
- Considered: Yes. Estimated cost: $4,200/month additional.
- Rejected: Buys ~4 months at current growth before hitting the next ceiling. Does not address the architectural bottleneck.

Alternative 2: Database read replicas
- Considered: Yes. Addresses read load but not write throughput.
- Rejected: Our bottleneck is write processing, not reads (see profiling data in Appendix A).
```

**API contract (bad):** "The endpoint will return user data in JSON."

**API contract (good):**
```
POST /v1/payments/queue
Request: { "amount": integer (cents), "currency": "USD"|"EUR", "idempotency_key": string (UUID) }
Response 202: { "job_id": string, "estimated_completion_ms": integer }
Response 400: { "error": "invalid_currency", "message": string }
Response 409: { "error": "duplicate_idempotency_key", "job_id": string }
Response 429: { "error": "rate_limit_exceeded", "retry_after": integer (seconds) }
```

## Common Mistakes

- **Writing the spec after the implementation:** Post-hoc specs rationalize decisions already made instead of evaluating them. Reviewers can't change the design when the code is already in review.
- **Skipping operational concerns:** A spec that describes what to build but not how to deploy, monitor, or roll back it is incomplete. The implementation is done when it is running safely in production, not when the PR is merged.
- **Too much detail on obvious parts, too little on hard parts:** Specs often over-explain standard CRUD operations and under-explain the one tricky distributed systems problem the whole design hinges on. Focus depth where the risk is.
- **No alternatives section:** This is the most commonly skipped section and the one reviewers most frequently request. Write it even when the answer is obvious — it builds trust that you evaluated the space.
- **Letting it go stale:** A spec that diverges from the implementation is a liability. Either update it when decisions change or add a header: "Implementation diverged from spec. See [PR link] for final design." Never silently let a spec become fiction.
