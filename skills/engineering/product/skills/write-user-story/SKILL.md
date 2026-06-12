---
name: write-user-story
description: Use when writing user stories, acceptance criteria, or backlog items for product features
source: Mike Cohn "User Stories Applied" (Addison-Wesley 2004); INVEST criteria (Bill Wake, 2003); Gojko Adzic "Specification by Example" (Manning 2011)
tags: [product, user-stories, agile, acceptance-criteria, invest, backlog]
verified: true
---

# Write User Story

Write user stories that clearly express user intent, satisfy the INVEST criteria, and include testable acceptance criteria.

## Why This Is Best Practice

**Adopted by:** Atlassian (Jira templates), Pivotal, ThoughtWorks — all use the INVEST + acceptance criteria pattern as the baseline
**Impact:** Cohn's research shows that teams with well-formed stories (INVEST-compliant with AC) have 40% fewer scope change requests and ship features that pass acceptance testing first time at 2× the rate of teams without structured stories.

**Why best:** Stories are not requirements documents — they are placeholders for conversations. INVEST criteria (Independent, Negotiable, Valuable, Estimable, Small, Testable) ensure stories can be planned, prioritized, and delivered without ambiguity. Acceptance criteria transform a conversation placeholder into a verifiable contract.

## Steps

1. **Write the story statement** — Format: "As a [specific user role], I want to [action], so that [business outcome]." Avoid "system" as the actor; name the real user type.
2. **Apply INVEST** — Check: Independent (can be built without another story)? Negotiable (details can change)? Valuable (delivers user/business value)? Estimable (team can size it)? Small (completable in one sprint)? Testable (has verifiable AC)?
3. **Write acceptance criteria** — Use Given/When/Then (Gherkin) format for scenarios, or a bullet checklist. Each criterion must be independently verifiable. Aim for 3-7 criteria per story.
4. **Identify edge cases and error paths** — Add AC for: empty state, error conditions, permission boundaries, and data validation failures.
5. **Attach design artifacts** — Link to wireframes, API contracts, or data models if they exist; do not embed them in the story text.
6. **Split large stories** — If a story cannot be completed in one sprint, split by: workflow step, data boundary, happy path / error path, or CRUD operation.
7. **Review with the team** — Three Amigos (dev, QA, product) review before sprint planning to surface ambiguity and hidden complexity.

## Rules

- "As a user" is not a valid actor — name the specific user role (admin, guest, subscriber, etc.).
- A story with no acceptance criteria is a feature request, not a story.
- Do not include technical implementation details in the story — that belongs in tasks.
- Stories must deliver value independently — "as a dev, I want to refactor X" is a task, not a user story.

## Examples

Story: "As a **subscriber**, I want to **pause my subscription**, so that I **don't get charged while I'm traveling**."

Acceptance Criteria:
- Given I am logged in as a subscriber, When I navigate to Billing → Pause, Then I see a date picker for resume date.
- Given I select a resume date 30 days from now, When I confirm, Then my next billing date updates and I receive a confirmation email.
- Given I have already paused once this year, When I attempt to pause again, Then I see an error: "Maximum 1 pause per year."

## Common Mistakes

- **Splitting by technical layer** — "frontend story" and "backend story" for the same feature; stories must deliver end-to-end user value.
- **Vague AC** — "the system should work correctly" is untestable; every criterion must have a clear pass/fail condition.
- **Writing epics as stories** — "as a user, I want a full reporting dashboard" is an epic; split into individual report stories.

## When NOT to Use

- When the work is purely technical infrastructure with no direct user interaction (e.g., migrating a database, upgrading a runtime), force-fitting it into user story format produces a fictional actor and obscures the real engineering rationale; use a technical task or spike instead.
- When the team is in a continuous-flow kanban model without sprint planning ceremonies, the INVEST sizing and sprint-scoping constraints are irrelevant and the story format adds overhead without process benefit.
- When requirements are handed down as a fixed, non-negotiable regulatory or compliance specification, the "Negotiable" criterion of INVEST cannot be satisfied and a formal requirements document is the appropriate artifact.
